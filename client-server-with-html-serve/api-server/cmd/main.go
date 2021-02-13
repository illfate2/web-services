package main

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/example/dataloader"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"

	"github.com/illfate2/web-services/client-server-with-html-serve/api-server/pkg/api/generated"
	"github.com/illfate2/web-services/client-server-with-html-serve/api-server/pkg/api/graph"
	repo "github.com/illfate2/web-services/client-server-with-html-serve/api-server/pkg/repositories"
	service "github.com/illfate2/web-services/client-server-with-html-serve/api-server/pkg/services"
)

type serverConfig struct {
	ServerPort  string `envconfig:"SERVER_PORT"`
	DBAddr      string `envconfig:"DB_ADDR"`
	LooFilePath string `envconfig:"LOG_FILE_PATH" default:"log.txt"`
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
	defer conn.Close()

	router := chi.NewRouter()
	router.Use(dataloader.LoaderMiddleware)

	router.Handle("/", playground.Handler("Dataloader", "/query"))
	router.Handle("/query", handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: graph.NewResolver(service.NewService(repo)),
	})))

	log.Fatal(http.ListenAndServe(":"+c.ServerPort, router))
}
