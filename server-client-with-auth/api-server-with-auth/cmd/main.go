package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"xorm.io/xorm"

	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/api/rest"
	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/auth"
	repo "github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/repositories"
	service "github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/services"
)

type serverConfig struct {
	ServerPort              string `envconfig:"SERVER_PORT"`
	DBAddr                  string `envconfig:"DB_ADDR"`
	LooFilePath             string `envconfig:"LOG_FILE_PATH" default:"log.txt"`
	JWTSecret               string `envconfig:"SERVER_PORT" default:"secret"`
	OAuthGithubID           string `envconfig:"OAUTH_GITHUB_ID" required:"true"`
	OAuthGithubSecret       string `envconfig:"OAUTH_GITHUB_SECRET" required:"true"`
	OAuthGithubRedirectURL  string `envconfig:"OAUTH_GITHUB_REDIRECT_URL" required:"true"`
	OAuthGoogleID           string `envconfig:"OAUTH_GOOGLE_ID" required:"true"`
	OAuthGoogleSecret       string `envconfig:"OAUTH_GOOGLE_SECRET" required:"true"`
	OAuthGoogleRedirectURL  string `envconfig:"OAUTH_GOOGLE_REDIRECT_URL" required:"true"`
	OAuthDiscordID          string `envconfig:"OAUTH_DISCORD_ID" required:"true"`
	OAuthDiscordSecret      string `envconfig:"OAUTH_DISCORD_SECRET" required:"true"`
	OAuthDiscordRedirectURL string `envconfig:"OAUTH_DISCORD_REDIRECT_URL" required:"true"`
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var config serverConfig
	err := envconfig.Process("", &config)
	if err != nil {
		panic(err)
	}
	conn, err := pgxpool.Connect(context.TODO(), config.DBAddr) // todo
	engine, err := xorm.NewEngine("postgres", config.DBAddr)
	panicOnErr(err)
	defer engine.Close()
	repo := repo.New(conn, engine)
	panicOnErr(err)
	defer conn.Close()
	service := service.NewService(repo)
	jwtService := auth.NewJWTService([]byte(config.JWTSecret))
	server := rest.NewServer(jwtService, service,
		rest.OAuthConfig{
			ProviderType: rest.Google,
			ClientID:     config.OAuthGoogleID,
			ClientSecret: config.OAuthGoogleSecret,
			RedirectURL:  config.OAuthGoogleRedirectURL,
		},
		rest.OAuthConfig{
			ProviderType: rest.Github,
			ClientID:     config.OAuthGithubID,
			ClientSecret: config.OAuthGithubSecret,
			RedirectURL:  config.OAuthGithubRedirectURL,
		},
		rest.OAuthConfig{
			ProviderType: rest.Discord,
			ClientID:     config.OAuthDiscordID,
			ClientSecret: config.OAuthDiscordSecret,
			RedirectURL:  config.OAuthDiscordRedirectURL,
		},
	)
	log.Fatal(http.ListenAndServe(":"+config.ServerPort, server))
}
