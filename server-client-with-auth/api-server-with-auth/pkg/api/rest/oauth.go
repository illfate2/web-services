package rest

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/gologin/v2/github"
	"github.com/pollen5/discord-oauth2"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/auth"
	service "github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/services"
)

type OAuth struct {
	state      string
	service    *service.Service
	jwtService *auth.JWTService
	configs    map[OAuthProvider]*oauth2.Config
}

type OAuthProvider string

const (
	Google  OAuthProvider = "google"
	Discord OAuthProvider = "discord"
	Github  OAuthProvider = "github"
)

type OAuthConfig struct {
	ProviderType OAuthProvider
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func NewOAuth(service *service.Service, jwtService *auth.JWTService) *OAuth {
	return &OAuth{
		state:      "thisshouldberandom",
		service:    service,
		jwtService: jwtService,
		configs: map[OAuthProvider]*oauth2.Config{
			Google: {
				Scopes:   []string{"https://www.googleapis.com/auth/userinfo.email"},
				Endpoint: google.Endpoint,
			},
			Github: {
				Endpoint: githuboauth.Endpoint,
			},
			Discord: {
				Scopes:   []string{discord.ScopeIdentify},
				Endpoint: discord.Endpoint,
			},
		},
	}
}

func (o *OAuth) GetConfig(provider OAuthProvider) *oauth2.Config {
	return o.configs[provider]
}

func (o *OAuth) WithConfig(cfg OAuthConfig) {
	config, ok := o.configs[cfg.ProviderType]
	if !ok {
		return
	}
	config.ClientID = cfg.ClientID
	config.ClientSecret = cfg.ClientSecret
	config.RedirectURL = cfg.RedirectURL
}

func (o *OAuth) HandleCallBackFromGoogle(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != o.state {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := o.GetConfig(Google).Exchange(context.Background(), code)
	if err != nil {
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	id := gjson.GetBytes(response, "id").Str
	user, err := o.service.SignupWithOAuth(id, "google")
	if err != nil {
		log.Printf("Got error from signup google: %v", err)
		return
	}
	o.setAuthInfo(w, r, user.ID)

}

func (o *OAuth) setAuthInfo(w http.ResponseWriter, req *http.Request, userID int) {
	accessToken, _ := o.jwtService.GenerateAccessToken(userID)
	tokenResp := struct {
		AccessToken string
	}{
		AccessToken: accessToken,
	}
	cookie := http.Cookie{
		Name:     "auth",
		Value:    tokenResp.AccessToken,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: false,
		MaxAge:   50000,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, req, "http://localhost:3000/", http.StatusFound)
}

func (o *OAuth) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	googleCfg := o.GetConfig(Google)
	authURL, _ := url.Parse(googleCfg.Endpoint.AuthURL)
	parameters := url.Values{}
	parameters.Add("client_id", googleCfg.ClientID)
	parameters.Add("scope", strings.Join(googleCfg.Scopes, " "))
	parameters.Add("redirect_uri", googleCfg.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", o.state)
	authURL.RawQuery = parameters.Encode()
	http.Redirect(w, r, authURL.String(), http.StatusTemporaryRedirect)
}

func (o *OAuth) issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user, err := o.service.SignupWithOAuth(strconv.FormatInt(*githubUser.ID, 10), "github")
		if err != nil {
			log.Printf("Got error signup github: %v", err)
			return
		}
		o.setAuthInfo(w, req, user.ID)
		http.Redirect(w, req, "/query", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

func (o *OAuth) discordCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != o.state {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("State does not match."))
		return
	}
	discordCfg := o.GetConfig(Discord)
	token, err := discordCfg.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	res, err := discordCfg.Client(context.Background(), token).Get("https://discordapp.com/api/v6/users/@me")
	if err != nil || res.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte(res.Status))
		}
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	id := gjson.GetBytes(body, "id").Str
	user, err := o.service.SignupWithOAuth(id, "discord")
	if err != nil {
		log.Printf("Got error: %v", err)
		return
	}
	o.setAuthInfo(w, r, user.ID)
}
