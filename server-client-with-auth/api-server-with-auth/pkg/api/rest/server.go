package rest

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/github"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/api/generated"
	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/api/graph"
	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/auth"
	service "github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/services"
)

type Server struct {
	http.Handler
	oauth      *OAuth
	router     *chi.Mux
	service    *service.Service
	jwtService *auth.JWTService
}

func NewServer(jwtSvc *auth.JWTService, service *service.Service, configs ...OAuthConfig) *Server {
	router := chi.NewRouter()
	server := &Server{
		oauth:      NewOAuth(service, jwtSvc),
		Handler:    router,
		jwtService: jwtSvc,
		service:    service,
	}
	for _, c := range configs {
		server.oauth.WithConfig(c)
	}
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() { h.ServeHTTP(w, r) }()
			bearerToken := r.Header.Get("Authorization")
			if bearerToken == "" {
				return
			}
			token := strings.Split(bearerToken, " ")[1]
			userID, err := jwtSvc.UserIDFromToken(token)
			if err != nil {
				return
			}
			ctx := context.WithValue(r.Context(), "userID", userID)
			r = r.WithContext(ctx)
		})
	})
	router.Use(corsMiddleware())
	router.Handle("/login/github", github.StateHandler(
		gologin.DebugOnlyCookieConfig,
		github.LoginHandler(server.oauth.GetConfig(Github), nil)))
	router.Handle("/callback/github", github.StateHandler(gologin.DebugOnlyCookieConfig,
		github.CallbackHandler(server.oauth.GetConfig(Github), server.oauth.issueSession(), nil)))

	router.HandleFunc("/login/google", server.oauth.HandleGoogleLogin)
	router.HandleFunc("/callback/google", func(w http.ResponseWriter, req *http.Request) {
		server.oauth.HandleCallBackFromGoogle(w, req)
	})
	router.HandleFunc("/login/discord", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, server.oauth.GetConfig(Discord).AuthCodeURL(server.oauth.state), http.StatusTemporaryRedirect)
	})
	router.HandleFunc("/callback/discord", func(w http.ResponseWriter, req *http.Request) {
		server.oauth.discordCallback(w, req)
	})
	router.Handle("/", playground.Handler("Dataloader", "/query"))
	router.HandleFunc("/query", server.handleGQLEndpoint)
	return server
}

func (s *Server) handleGQLEndpoint(w http.ResponseWriter, req *http.Request) {
	handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: graph.NewResolver(s.service, s.jwtService),
		Directives: generated.DirectiveRoot{
			Auth: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
				return authMiddleware(ctx, req, s.jwtService, next)
			},
			Self: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
				return selfAccessMiddleware(ctx, s.service, next)
			},
		},
	})).ServeHTTP(w, req)
}

func selfAccessMiddleware(ctx context.Context, repo *service.Service, next graphql.Resolver) (res interface{}, err error) {
	fieldContext := graphql.GetFieldContext(ctx)
	alias := fieldContext.Field.Alias
	switch alias {
	case "updateMuseumItem", "deleteMuseumItem":
		err := validateCreatedByWithCtxUser(ctx, func(i int) (int, error) {
			item, err := repo.GetMuseumSet(i)
			return item.CreatedBy, err
		}, fieldContext)
		if err != nil {
			return nil, err
		}
	case "updateMuseumSet", "deleteMuseumSet":
		err := validateCreatedByWithCtxUser(ctx, func(i int) (int, error) {
			set, err := repo.GetMuseumSet(i)
			return set.CreatedBy, err
		}, fieldContext)
		if err != nil {
			return nil, err
		}
	case "updateMuseumItemMovement", "deleteMuseumItemMovement":
		err := validateCreatedByWithCtxUser(ctx, func(i int) (int, error) {
			m, err := repo.GetMuseumItemMovement(i)
			return m.CreatedBy, err
		}, fieldContext)
		if err != nil {
			return nil, err
		}
	case "updateMuseumFund", "deleteMuseumFund":
		err := validateCreatedByWithCtxUser(ctx, func(i int) (int, error) {
			fund, err := repo.GetMuseumFundByID(i)
			return fund.CreatedBy, err
		}, fieldContext)
		if err != nil {
			return nil, err
		}
	}
	return next(ctx)
}

func authMiddleware(ctx context.Context, req *http.Request, jwtService *auth.JWTService, next graphql.Resolver) (res interface{}, err error) {
	bearerToken := req.Header.Get("Authorization")
	if bearerToken == "" {
		return nil, errors.New("empty token")
	}
	token := strings.Split(bearerToken, " ")[1]
	_, err = jwtService.UserIDFromToken(token)
	if err != nil {
		return nil, err
	}
	return next(ctx)
}
func validateCreatedByWithCtxUser(ctx context.Context, f func(int) (int, error), fieldContext *graphql.FieldContext) error {
	id, err := f(fieldContext.Args["id"].(int))
	if err != nil {
		return err
	}
	if id != ctx.Value("userID").(int) {
		return errors.New("token user id and created by are not the same")
	}
	return nil
}

func corsMiddleware() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}
