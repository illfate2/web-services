package api

import (
	"log"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/illfate2/web-services/server-with-html-serve/pkg/entities"
)

func (s *Server) initMuseumSet(e *echo.Echo) {
	e.GET("/museumSets", s.getMuseumSets)
	e.GET("/museumSet/:id", s.getMuseumSet)
	e.GET("/museumSet", func(c echo.Context) error {
		_ = s.tmpl.ExecuteTemplate(c.Response().Writer, "New museum item set", nil)
		return nil
	})
	e.POST("/museumSet", s.createMuseumSet)
	e.GET("/deleteMuseumSet/:id", s.deleteMuseumSet)
	e.GET("/editMuseumSet/:id", s.getEditMuseumItempage)
	e.POST("/editMuseumSet/:id", s.updateMuseumItem)
}

func (s *Server) getMuseumSets(c echo.Context) error {
	items, err := s.service.GetMuseumSets()
	if err != nil {
		log.Printf("Failed to find museum sets: %+v", err)
		return err
	}
	_ = s.tmpl.ExecuteTemplate(c.Response().Writer, "Museum sets", items)
	return nil
}

func (s *Server) getMuseumSet(c echo.Context) error {
	id := getIDFromURL(c)
	set, err := s.service.GetMuseumSet(id)
	if err != nil {
		log.Printf("Failed to find item with details: %s", err)
		return err
	}

	err = s.tmpl.ExecuteTemplate(c.Response(), "ShowMuseumSet", set)
	return nil
}

func (s *Server) createMuseumSet(c echo.Context) error {
	setName := c.FormValue("name")
	_, err := s.service.CreateMuseumSet(entities.MuseumSet{
		Name: setName,
	})
	if err != nil {
		return err
	}
	return c.Redirect(301, "/museumSets")
}

func (s *Server) getEditMuseumSetPage(c echo.Context) error {
	id := c.Param("id")
	parsedID, _ := strconv.ParseInt(id, 10, 64)
	item, err := s.service.GetMuseumFundByID(int(parsedID))
	if err != nil {
		return err
	}
	_ = s.tmpl.ExecuteTemplate(c.Response(), "Edit Museum Set", item)
	return nil
}

func (s *Server) updateMuseumSet(c echo.Context) error {
	fund := entities.MuseumFund{
		ID:   getIDFromURL(c),
		Name: getNameFromForm(c),
	}
	err := s.service.UpdateMuseumFund(fund)
	if err != nil {
		return err
	}
	return c.Redirect(301, "/museumSets")
}

func (s *Server) deleteMuseumSet(c echo.Context) error {
	err := s.service.DeleteMuseumSet(getIDFromURL(c))
	if err != nil {
		return err
	}
	return c.Redirect(301, "/museumSets")
}
