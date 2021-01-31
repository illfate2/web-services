package service

import (
	"github.com/pkg/errors"

	"github.com/illfate2/web-services/server-with-html-serve/pkg/entities"
)

// todo: add tx
func (s *Service) CreateMuseumItem(item entities.MuseumItemWithDetails) (entities.MuseumItemWithDetails, error) {
	var err error
	item.MuseumItem, err = s.insertMuseumItem(item.MuseumItem)
	if err != nil {
		return entities.MuseumItemWithDetails{}, errors.Wrap(err, "failed to insert item")
	}

	return item, nil
}

func (s *Service) GetMuseumItem(id int) (entities.MuseumItem, error) {
	return s.repo.FindMuseumItem(id)
}

func (s *Service) GetMuseumItemByName(name string) (entities.MuseumItem, error) {
	return s.repo.FindMuseumItemByName(name)
}

func (s *Service) GetMuseumItemWithDetails(id int) (entities.MuseumItemWithDetails, error) {
	return s.repo.FindMuseumItemWithDetails(id)
}

func (s *Service) FindMuseumItems(args entities.SearchMuseumItemsArgs) ([]entities.MuseumItem, error) {
	return s.repo.FindMuseumItems(args)
}

func (s *Service) UpdateMuseumItem(item entities.MuseumItem) error {
	return s.repo.UpdateMuseumItem(item)
}

func (s *Service) DeleteMuseumItem(id int) error {
	return s.repo.DeleteMuseumItem(id)
}

func (s *Service) insertMuseumItem(item entities.MuseumItem) (entities.MuseumItem, error) {
	return s.repo.InsertMuseumItem(item)
}
