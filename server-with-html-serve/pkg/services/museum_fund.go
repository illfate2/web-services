package service

import (
	"github.com/illfate2/web-services/server-with-html-serve/pkg/entities"
)

func (s *Service) CreateMuseumFund(fund entities.MuseumFund) (entities.MuseumFund, error) {
	return s.repo.InsertMuseumFund(fund)
}
