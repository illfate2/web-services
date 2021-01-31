package server

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	http.Handler
	tmpl *template.Template
}

func NewServer() *Server {
	e := echo.New()
	s := &Server{
		Handler: e,
		tmpl:    template.Must(template.ParseGlob("static/*")),
	}
	return s
}
