package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/github"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/api/generated"
	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/api/graph"
	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/api/rest"
	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/auth"
	repo "github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/repositories"
	service "github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/services"
)

type serverConfig struct {
	ServerPort  string `envconfig:"SERVER_PORT"`
	DBAddr      string `envconfig:"DB_ADDR"`
	LooFilePath string `envconfig:"LOG_FILE_PATH" default:"log.txt"`
	JWTSecret   string `envconfig:"SERVER_PORT" default:"secret"`
}

func main() {
	var c serverConfig
	err := envconfig.Process("", &c)
	if err != nil {
		panic(err)
	}
	conn, err := pgxpool.Connect(context.TODO(), c.DBAddr) // todo
	repo := repo.New(conn)
	if err != nil {
		panic(err)
	}
	service := service.NewService(repo)
	defer conn.Close()
	jwtService := auth.NewJWTService([]byte(c.JWTSecret))
	router := chi.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() { h.ServeHTTP(w, r) }()
			bearerToken := r.Header.Get("Authorization")
			if bearerToken == "" {
				return
			}
			token := strings.Split(bearerToken, " ")[1]
			userID, err := jwtService.UserIDFromToken(token)
			if err != nil {
				return
			}
			ctx := context.WithValue(r.Context(), "userID", userID)
			r = r.WithContext(ctx)
		})
	})
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Handle("/login/github", github.StateHandler(
		gologin.DebugOnlyCookieConfig,
		github.LoginHandler(rest.OauthGithubConf, nil)))
	router.Handle("/callback/github", github.StateHandler(gologin.DebugOnlyCookieConfig,
		github.CallbackHandler(rest.OauthGithubConf, rest.IssueSession(), nil)))

	router.HandleFunc("/login/google", rest.HandleGoogleLogin)
	router.HandleFunc("/callback/google", func(w http.ResponseWriter, req *http.Request) {
		rest.HandleCallBackFromGoogle(w, req, service, repo, jwtService)
	})
	router.Handle("/", playground.Handler("Dataloader", "/query"))
	router.HandleFunc("/query", func(w http.ResponseWriter, req *http.Request) {
		handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
			Resolvers: graph.NewResolver(service, jwtService),
			Directives: generated.DirectiveRoot{
				Auth: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
					return authMiddleware(ctx, req, jwtService, next)
				},
				Self: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
					return selfAccessMiddleware(ctx, repo, next)
				},
			},
		})).ServeHTTP(w, req)
	})
	log.Fatal(http.ListenAndServe(":"+c.ServerPort, router))
}

func selfAccessMiddleware(ctx context.Context, repo *repo.Repo, next graphql.Resolver) (res interface{}, err error) {
	fieldContext := graphql.GetFieldContext(ctx)
	alias := fieldContext.Field.Alias
	switch alias {
	case "updateMuseumItem", "deleteMuseumItem":
		err := validateCreatedByWithCtxUser(ctx, func(i int) (int, error) {
			item, err := repo.FindMuseumItem(i)
			return item.CreatedBy, err
		}, fieldContext)
		if err != nil {
			return nil, err
		}
	case "updateMuseumSet", "deleteMuseumSet":
		err := validateCreatedByWithCtxUser(ctx, func(i int) (int, error) {
			set, err := repo.FindMuseumSet(i)
			return set.CreatedBy, err
		}, fieldContext)
		if err != nil {
			return nil, err
		}
	case "updateMuseumItemMovement", "deleteMuseumItemMovement":
		err := validateCreatedByWithCtxUser(ctx, func(i int) (int, error) {
			m, err := repo.FindMuseumItemMovement(i)
			return m.CreatedBy, err
		}, fieldContext)
		if err != nil {
			return nil, err
		}
	case "updateMuseumFund", "deleteMuseumFund":
		err := validateCreatedByWithCtxUser(ctx, func(i int) (int, error) {
			fund, err := repo.FindMuseumFundByID(i)
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
