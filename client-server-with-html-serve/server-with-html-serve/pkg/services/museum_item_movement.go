package service

import (
	"log"

	"github.com/illfate2/web-services/client-server-with-html-serve/server-with-html-serve/pkg/entities"
)

func (s *Service) CreateMuseumItemMovement(movement entities.MuseumItemMovement) (entities.MuseumItemMovement, error) {
	movement.MuseumItemID = movement.Item.ID
	movement, err := s.insertMuseumItemMovement(movement)
	if err != nil {
		log.Print(err)
		return entities.MuseumItemMovement{}, err
	}
	return movement, nil
}

func (s *Service) GetMuseumItemMovement(id int) (entities.MuseumItemMovement, error) {
	return s.repo.FindMuseumItemMovement(id)
}

func (s *Service) GetMuseumItemMovements() ([]entities.MuseumItemMovement, error) {
	return s.repo.FindMuseumItemMovements()
}

func (s *Service) insertMuseumItemMovement(movement entities.MuseumItemMovement) (entities.MuseumItemMovement, error) {
	return s.repo.InsertMuseumItemMovement(movement)
}

func (s *Service) UpdateMuseumItemMovement(movement entities.MuseumItemMovement) error {
	return s.repo.UpdateMuseumItemMovement(movement)
}

func (s *Service) DeleteMuseumItemMovement(id int) error {
	return s.repo.DeleteMuseumItemMovement(id)
}
