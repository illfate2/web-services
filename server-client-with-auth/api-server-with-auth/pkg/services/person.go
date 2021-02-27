package service

import (
	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/entities"
)

func (s *Service) CreatePerson(person entities.Person) (entities.Person, error) {
	return s.repo.InsertPerson(person)
}

func (s *Service) FindPerson(id int) (entities.Person, error) {
	return s.repo.FindPerson(id)
}
