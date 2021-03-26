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
	ServerPort  string `envconfig:"SERVER_PORT"`
	DBAddr      string `envconfig:"DB_ADDR"`
	LooFilePath string `envconfig:"LOG_FILE_PATH" default:"log.txt"`
	JWTSecret   string `envconfig:"SERVER_PORT" default:"secret"`
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var c serverConfig
	err := envconfig.Process("", &c)
	if err != nil {
		panic(err)
	}
	conn, err := pgxpool.Connect(context.TODO(), c.DBAddr) // todo
	engine, err := xorm.NewEngine("postgres", c.DBAddr)
	panicOnErr(err)
	defer engine.Close()
	repo := repo.New(conn, engine)
	panicOnErr(err)
	defer conn.Close()
	service := service.NewService(repo)
	jwtService := auth.NewJWTService([]byte(c.JWTSecret))
	server := rest.NewServer(jwtService, service)
	log.Fatal(http.ListenAndServe(":"+c.ServerPort, server))
}
