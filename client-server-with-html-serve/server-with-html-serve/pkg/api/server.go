package api

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	service "github.com/illfate2/web-services/client-server-with-html-serve/server-with-html-serve/pkg/services"
)

type Server struct {
	http.Handler
	service *service.Service
	tmpl    *template.Template
}

func NewServer(serv *service.Service, writer io.Writer) *Server {
	e := echo.New()
	s := &Server{
		Handler: e,
		service: serv,
		tmpl:    template.Must(template.ParseGlob("static/*")),
	}

	fileLogger := log.New(writer, "", 0)
	e.Use(func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fileLogger.Print(c.Request())
			return handlerFunc(c)
		}
	})
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
