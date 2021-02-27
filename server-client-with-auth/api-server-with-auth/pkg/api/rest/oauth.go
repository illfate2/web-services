package rest

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dghubble/gologin/v2/github"
	"github.com/jackc/pgx/v4"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/auth"
	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/entities"
	repo "github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/repositories"
	service "github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/services"
)

type Auth struct {
}

var (
	oauthGoogleConf = &oauth2.Config{
		ClientID:     "920271564807-d49j9ob3b02li21rbcl0u75k0brdnf45.apps.googleusercontent.com",
		ClientSecret: "1ZYVDfGj4-C7gefMZs88b8tR",
		RedirectURL:  "http://localhost:8082/callback/google",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	OauthGithubConf = &oauth2.Config{
		ClientID:     "d8efd64a876723cf1c30",
		ClientSecret: "964cb7bd972d91de80108a1c665f4623926ac15d",
		RedirectURL:  "http://localhost:8082/callback/github",
		Endpoint:     githuboauth.Endpoint,
	}
	oauthStateString = "thisshouldberandom"
)

func HandleCallBackFromGoogle(w http.ResponseWriter, r *http.Request,
	service *service.Service, repo *repo.Repo, jwtService *auth.JWTService) {
	state := r.FormValue("state")
	if state != oauthStateString {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code") // TODO possible empty
	token, err := oauthGoogleConf.Exchange(context.Background(), code)
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

	email := gjson.GetBytes(response, "email").Str
	id := gjson.GetBytes(response, "id").Str
	user, err := repo.FindUserByEmail(email)
	if err != nil && err != pgx.ErrNoRows {
		return
	}
	if err == pgx.ErrNoRows {
		user, err = service.CreateUser(entities.User{
			Email:    email,
			Password: id,
		})
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	}
	accessToken, _ := jwtService.GenerateAccessToken(user.ID)
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
	http.Redirect(w, r, "http://localhost:3000/", http.StatusFound)
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	authURL, _ := url.Parse(oauthGoogleConf.Endpoint.AuthURL)
	parameters := url.Values{}
	parameters.Add("client_id", oauthGoogleConf.ClientID)
	parameters.Add("scope", strings.Join(oauthGoogleConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthGoogleConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	authURL.RawQuery = parameters.Encode()
	http.Redirect(w, r, authURL.String(), http.StatusTemporaryRedirect)
}

func IssueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, "/query", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
