package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/shurcooL/graphql"

	"github.com/illfate2/web-services/client-server-with-html-serve/server-with-html-serve/pkg/api"
	repo "github.com/illfate2/web-services/client-server-with-html-serve/server-with-html-serve/pkg/repositories"
	service "github.com/illfate2/web-services/client-server-with-html-serve/server-with-html-serve/pkg/services"
)

type serverConfig struct {
	ServerPort  string `envconfig:"SERVER_PORT" default:"5555"`
	LooFilePath string `envconfig:"LOG_FILE_PATH" default:"log.txt"`
}

func main() {
	var c serverConfig
	err := envconfig.Process("", &c)
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile(c.LooFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	client := graphql.NewClient("http://localhost:8082/query", nil) // todo
	r := repo.New(client)
	s := api.NewServer(service.NewService(r), file)
	log.Print("Starting server")
	log.Fatal(http.ListenAndServe(":"+c.ServerPort, s))
}
