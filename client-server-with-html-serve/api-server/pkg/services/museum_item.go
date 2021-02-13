package service

import (
	"github.com/pkg/errors"

	"github.com/illfate2/web-services/client-server-with-html-serve/api-server/pkg/entities"
)

func (s *Service) CreateMuseumItem(item entities.MuseumItemWithDetails) (entities.MuseumItemWithDetails, error) {
	var err error
	item.Keeper, err = s.CreatePerson(item.Keeper)
	if err != nil {
		return entities.MuseumItemWithDetails{}, err
	}
	item.KeeperID = item.Keeper.ID

	item.MuseumItem, err = s.insertMuseumItem(item.MuseumItem)
	if err != nil {
		return entities.MuseumItemWithDetails{}, errors.Wrap(err, "failed to insert item")
	}

	return item, nil
}

func (s *Service) GetMuseumItem(id int) (entities.MuseumItem, error) {
	return s.repo.FindMuseumItem(id)
}

func (s *Service) GetMuseumItemWithDetails(id int) (entities.MuseumItemWithDetails, error) {
	return s.repo.FindMuseumItemWithDetails(id)
}

func (s *Service) SearchMuseumItems(args entities.SearchMuseumItemsArgs) ([]entities.MuseumItem, error) {
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
