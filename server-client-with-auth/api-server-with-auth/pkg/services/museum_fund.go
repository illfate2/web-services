package service

import (
	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/entities"
)

func (s *Service) CreateMuseumFund(fund entities.MuseumFund) (entities.MuseumFund, error) {
	return s.repo.InsertMuseumFund(fund)
}

func (s *Service) GetMuseumFundByID(id int) (entities.MuseumFund, error) {
	return s.repo.FindMuseumFundByID(id)
}

func (s *Service) GetMuseumFunds() ([]entities.MuseumFund, error) {
	return s.repo.FindMuseumFunds()
}

func (s *Service) UpdateMuseumFund(fund entities.MuseumFund) error {
	return s.repo.UpdateMuseumFund(fund)
}

func (s *Service) DeleteMuseumFund(id int) error {
	return s.repo.DeleteMuseumFund(id)
}
