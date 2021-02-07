package api

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"

	service "github.com/illfate2/web-services/server-with-html-serve/pkg/services"
)

type Server struct {
	http.Handler
	service *service.Service
	tmpl    *template.Template
}

func NewServer(serv *service.Service) *Server {
	e := echo.New()
	s := &Server{
		Handler: e,
		service: serv,
		tmpl:    template.Must(template.ParseGlob("static/*")),
	}
	s.initMuseumItemAPI(e)
	s.initMuseumItemMovement(e)
	s.initMuseumSet(e)
	s.initMuseumFund(e)
	e.GET("/", func(c echo.Context) error {
		_ = s.tmpl.ExecuteTemplate(c.Response().Writer, "Home page", nil)
		return nil
	})
	e.Static("/", "static")
	return s
}
