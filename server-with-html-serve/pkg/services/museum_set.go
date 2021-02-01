package service

import (
	"github.com/illfate2/web-services/server-with-html-serve/pkg/entities"
)

func (s *Service) CreateMuseumSet(set entities.MuseumSet) (entities.MuseumSet, error) {
	set, err := s.repo.InsertMuseumSet(set)
	if err != nil {
		return entities.MuseumSet{}, err
	}
	return set, nil
}

func (s *Service) GetMuseumSets() ([]entities.MuseumSet, error) {
	sets, err := s.repo.FindMuseumSets()
	if err != nil {
		return nil, err
	}
	return sets, nil
}

func (s *Service) FindMuseumSet(id int) (entities.MuseumSetWithDetails, error) {
	return s.repo.FindMuseumSet(id)
}
