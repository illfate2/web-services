package api

import (
	"log"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/illfate2/web-services/client-server-with-html-serve/server-with-html-serve/pkg/entities"
)

func (s *Server) initMuseumFund(e *echo.Echo) {
	e.GET("/museumFunds", s.getMuseumFunds)
	e.GET("/museumFund/:id", s.getMuseumFund)
	e.GET("/museumFund", func(c echo.Context) error {
		_ = s.tmpl.ExecuteTemplate(c.Response().Writer, "New museum item fund", nil)
		return nil
	})
	e.POST("/museumFund", s.createMuseumFund)
	e.GET("/deleteMuseumFund/:id", s.deleteMuseumFund)
	e.GET("/editMuseumFund/:id", s.getEditMuseumFundPage)
	e.POST("/editMuseumFund/:id", s.updateMuseumFund)
}

func (s *Server) getMuseumFunds(c echo.Context) error {
	items, err := s.service.GetMuseumFunds()
	if err != nil {
		log.Printf("Failed to find museum sets: %+v", err)
		return err
	}
	_ = s.tmpl.ExecuteTemplate(c.Response().Writer, "Museum funds", items)
	return nil
}

func (s *Server) getMuseumFund(c echo.Context) error {
	id := getIDFromURL(c)
	set, err := s.service.GetMuseumFundByID(id)
	if err != nil {
		log.Printf("Failed to find item with details: %s", err)
		return err
	}

	err = s.tmpl.ExecuteTemplate(c.Response(), "Museum fund", set)
	return nil
}

func (s *Server) createMuseumFund(c echo.Context) error {
	fundName := c.FormValue("name")
	_, err := s.service.CreateMuseumFund(entities.MuseumFund{
		Name: fundName,
	})
	if err != nil {
		return err
	}
	return c.Redirect(301, "/museumFunds")
}

func (s *Server) getEditMuseumFundPage(c echo.Context) error {
	id := c.Param("id")
	parsedID, _ := strconv.ParseInt(id, 10, 64)
	item, err := s.service.GetMuseumFundByID(int(parsedID))
	if err != nil {
		return err
	}
	_ = s.tmpl.ExecuteTemplate(c.Response(), "Edit Museum Fund", item)
	return nil
}

func (s *Server) updateMuseumFund(c echo.Context) error {
	fund := entities.MuseumFund{
		ID:   getIDFromURL(c),
		Name: getNameFromForm(c),
	}
	err := s.service.UpdateMuseumFund(fund)
	if err != nil {
		log.Print(err)
		return err
	}
	return c.Redirect(301, "/museumFunds")
}

func (s *Server) deleteMuseumFund(c echo.Context) error {
	err := s.service.DeleteMuseumFund(getIDFromURL(c))
	if err != nil {
		return err
	}
	return c.Redirect(301, "/museumFunds")
}
