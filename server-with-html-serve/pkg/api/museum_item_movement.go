package api

import (
	"log"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/illfate2/web-services/server-with-html-serve/pkg/entities"
)

func (s *Server) initMuseumItemMovement(e *echo.Echo) {
	e.POST("/museumItemMovement", s.createMuseumItemMovement)
	e.GET("/museumItemMovements", s.getMuseumItemMovements)

	e.GET("/museumItemMovement", s.getMuseumItemMovementPage)
	e.GET("/museumItemMovement/:id", s.getMuseumItemMovement)

	e.GET("/deleteMuseumItemMovement/:id", s.deleteMuseumMovement)
	e.GET("/editMuseumItemMovement/:id", s.getEditMuseumItemMovementPage)
	e.POST("/editMuseumItemMovement/:id", s.updateMuseumItemMovement)
}

func (s *Server) getMuseumItemMovements(c echo.Context) error {
	movements, err := s.service.GetMuseumItemMovements()
	if err != nil {
		log.Printf("Failed to find museum item movement: %+v", err)
		return err
	}
	_ = s.tmpl.ExecuteTemplate(c.Response().Writer, "Museum item Movements", movements)
	return nil
}

func (s *Server) getMuseumItemMovement(c echo.Context) error {
	id := getIDFromURL(c)
	item, err := s.service.GetMuseumItemMovement(id)
	if err != nil {
		log.Printf("Failed to find item with details: %s", err)
		return err
	}

	_ = s.tmpl.ExecuteTemplate(c.Response(), "ShowMuseumItemMovement", item)
	return nil
}

func (s *Server) createMuseumItemMovement(c echo.Context) error {
	_, err := s.service.CreateMuseumItemMovement(getMovementFromForm(c))
	if err != nil {
		log.Print(err)
		return err
	}
	return c.Redirect(301, "/museumItemMovements")
}

func (s *Server) getMuseumItemMovementPage(c echo.Context) error {
	items, err := s.service.SearchMuseumItems(entities.SearchMuseumItemsArgs{})
	if err != nil {
		return err
	}
	_ = s.tmpl.ExecuteTemplate(c.Response().Writer, "New museum item movement", items)
	return nil
}

func (s *Server) getEditMuseumItemMovementPage(c echo.Context) error {
	id := c.Param("id")
	parsedID, _ := strconv.ParseInt(id, 10, 64)
	movement, err := s.service.GetMuseumItemMovement(int(parsedID))
	if err != nil {
		return err
	}
	_ = s.tmpl.ExecuteTemplate(c.Response(), "Update Museum Item Movement", movement)
	return nil
}

func (s *Server) updateMuseumItemMovement(c echo.Context) error {
	id := getIDFromURL(c)
	movement := getMovementFromForm(c)
	movement.ID = id
	err := s.service.UpdateMuseumItemMovement(movement)
	if err != nil {
		return err
	}
	return c.Redirect(301, "/museumItemMovements")
}

func (s *Server) deleteMuseumMovement(c echo.Context) error {
	err := s.service.DeleteMuseumItemMovement(getIDFromURL(c))
	if err != nil {
		return err
	}
	return c.Redirect(301, "/museumItemMovements")
}

func getMovementFromForm(c echo.Context) entities.MuseumItemMovement {
	var movement entities.MuseumItemMovement
	movement.AcceptDate = getParsedTime(c.FormValue("accept_date"))
	movement.ExhibitTransferDate = getParsedTime(c.FormValue("exhibit_transfer_date"))
	movement.ExhibitReturnDate = getParsedTime(c.FormValue("exhibit_return_date"))
	movement.ResponsiblePerson = getPersonFromForm(c)
	movement.Item.ID, _ = strconv.Atoi(c.FormValue("item"))
	return movement
}

func getParsedTime(t string) *time.Time {
	if t == "" {
		return nil
	}
	t = t + ":00"
	parsed, err := time.Parse("2006-01-02T15:04:05", t)
	log.Print(err)
	return &parsed
}
