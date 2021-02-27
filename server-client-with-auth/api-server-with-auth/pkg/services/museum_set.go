package service

import (
	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/entities"
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

func (s *Service) GetMuseumSet(id int) (entities.MuseumSetWithDetails, error) {
	return s.repo.FindMuseumSet(id)
}

func (s *Service) GetMuseumSetByName(name string) (entities.MuseumSet, error) {
	return s.repo.FindMuseumSetByName(name)
}

func (s *Service) DeleteMuseumSet(id int) error {
	return s.repo.DeleteMuseumSet(id)
}

func (s *Service) UpdateMuseumSet(set entities.MuseumSet) error {
	return s.repo.UpdateMuseumSet(set)
}
