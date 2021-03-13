package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/kelseyhightower/envconfig"

	"github.com/illfate2/web-services/server-with-html-serve/pkg/api"
	"github.com/illfate2/web-services/server-with-html-serve/pkg/repositories"
	service "github.com/illfate2/web-services/server-with-html-serve/pkg/services"
)

type serverConfig struct {
	ServerPort  string `envconfig:"SERVER_PORT"`
	DBAddr      string `envconfig:"DB_ADDR"`
	LogFilePath string `envconfig:"LOG_FILE_PATH" default:"log.log"`
}

func main() {
	var c serverConfig
	err := envconfig.Process("", &c)
	if err != nil {
		panic(err)
	}
	conn, err := pgx.Connect(context.TODO(), c.DBAddr) // todo
	repo := repo.New(conn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.TODO())
	file, err := os.OpenFile(c.LogFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	s := api.NewServer(service.NewService(repo), file)
	log.Print("Starting server")
	log.Fatal(http.ListenAndServe(":"+c.ServerPort, s))
}
