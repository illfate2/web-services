package api

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"

	service "github.com/illfate2/web-services/server-with-html-serve/pkg/services"
)

type Server struct {
	http.Handler
	service service.Service
	tmpl    *template.Template
}

func NewServer() *Server {
	e := echo.New()
	s := &Server{
		Handler: e,
		tmpl:    template.Must(template.ParseGlob("static/*")),
	}
	s.initMuseumItemAPI(e)
	e.Static("/", "static")
	return s
}
