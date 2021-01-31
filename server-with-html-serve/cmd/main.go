package main

import (
	"log"
	"net/http"

	"github.com/illfate2/web-services/server-with-html-serve/pkg/api"
)

func main() {
	s := api.NewServer()
	log.Fatal(http.ListenAndServe(":8080", s))
}
