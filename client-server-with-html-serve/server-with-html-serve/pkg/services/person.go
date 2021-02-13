package service

import (
	"github.com/illfate2/web-services/client-server-with-html-serve/server-with-html-serve/pkg/entities"
)

func (s *Service) CreatePerson(person entities.Person) (entities.Person, error) {
	person, err := s.repo.InsertPerson(person)
	if err != nil {
		return entities.Person{}, err
	}
	return person, nil
}
