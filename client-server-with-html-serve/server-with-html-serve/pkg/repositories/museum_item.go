package repo

import (
	"context"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/shurcooL/graphql"

	"github.com/illfate2/web-services/client-server-with-html-serve/server-with-html-serve/pkg/entities"
)

type museumItem struct {
	ID              graphql.ID      `json:"id"`
	Name            graphql.String  `json:"name"`
	Annotation      *graphql.String `json:"annotation"`
	InventoryNumber graphql.String  `json:"inventoryNumber"`
	CreationDate    time.Time       `json:"creationDate"`
}

func (s museumItem) convertToEntity() entities.MuseumItem {
	id, _ := strconv.Atoi(s.ID.(string))
	return entities.MuseumItem{
		ID:              id,
		Name:            string(s.Name),
		Annotation:      getStrFromPtr(s.Annotation),
		InventoryNumber: string(s.InventoryNumber),
		CreationDate:    entities.NewDate(s.CreationDate),
	}
}

type person struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	MiddleName string `json:"middleName"`
}

type museumItemWithDetails struct {
	museumItem
	Set    museumSet  `json:"set"`
	Fund   museumFund `json:"fund"`
	Person person     `json:"person"`
}

func (s museumItemWithDetails) convertToEntity() entities.MuseumItemWithDetails {
	return entities.MuseumItemWithDetails{
		MuseumItem: s.museumItem.convertToEntity(),
		Keeper: entities.Person{
			FirstName:  s.Person.FirstName,
			LastName:   s.Person.LastName,
			MiddleName: s.Person.MiddleName,
		},
		Set:  s.Set.convertToEntity(),
		Fund: s.Fund.convertToEntity(),
	}
}

func getStrFromPtr(s *graphql.String) string {
	if s == nil {
		return ""
	}
	return string(*s)
}

func (r *Repo) FindMuseumItem(id int) (entities.MuseumItem, error) {
	q := struct {
		MuseumItem museumItem `graphql:"museumItem(id: $id)"`
	}{}
	err := r.findByID(id, &q)
	if err != nil {
		return entities.MuseumItem{}, err
	}
	return q.MuseumItem.convertToEntity(), nil
}

func (r *Repo) FindMuseumItemWithDetails(id int) (entities.MuseumItemWithDetails, error) {
	q := struct {
		MuseumItem museumItemWithDetails `graphql:"museumItem(id: $id)"`
	}{}
	err := r.findByID(id, &q)
	if err != nil {
		return entities.MuseumItemWithDetails{}, err
	}
	return q.MuseumItem.convertToEntity(), nil
}

func (r *Repo) FindMuseumItems(args entities.SearchMuseumItemsArgs) ([]entities.MuseumItem, error) {
	q := struct {
		MuseumItems []museumItem `graphql:"museumItems"`
	}{}
	err := r.client.Query(context.TODO(), &q, nil)
	if err != nil {
		return nil, err
	}
	res := make([]entities.MuseumItem, 0, len(q.MuseumItems))
	for _, s := range q.MuseumItems {
		res = append(res, s.convertToEntity())
	}
	return res, nil
}

type personInput struct {
	FirstName  graphql.String `json:"firstName"`
	MiddleName graphql.String `json:"middleName"`
	LastName   graphql.String `json:"lastName"`
}

func graphqlIDFromInt(i int) graphql.ID {
	return graphql.ID(strconv.Itoa(i))
}
func intFromGraphqlID(i int) graphql.ID {
	return graphql.ID(strconv.Itoa(i))
}
func (r *Repo) InsertMuseumItem(item entities.MuseumItemWithDetails) (entities.MuseumItem, error) {
	m := struct {
		CreateMuseumItem museumItem `graphql:"createMuseumItem(input: $input)"`
	}{}
	type MuseumItemInput struct {
		Name            graphql.String  `json:"name"`
		InventoryNumber graphql.String  `json:"inventoryNumber"`
		Annotation      *graphql.String `json:"annotation"`
		CreationDate    time.Time       `json:"creationDate"`
		SetID           graphql.ID      `json:"setID"`
		FundID          graphql.ID      `json:"fundID"`
		PersonInput     personInput     `json:"personInput"`
	}
	input := MuseumItemInput{
		Name:            graphql.String(item.Name),
		InventoryNumber: graphql.String(item.InventoryNumber),
		CreationDate:    item.CreationDate.Time,
		SetID:           graphqlIDFromInt(item.MuseumSetID),
		FundID:          graphqlIDFromInt(item.MuseumFundID),
		PersonInput: personInput{
			FirstName:  graphql.String(item.Keeper.FirstName),
			MiddleName: graphql.String(item.Keeper.MiddleName),
			LastName:   graphql.String(item.Keeper.LastName),
		},
	}
	if item.Annotation != "" {
		input.Annotation = graphql.NewString(graphql.String(item.Annotation))
	}
	err := r.client.Mutate(context.TODO(), &m, map[string]interface{}{
		"input": input,
	})
	if err != nil {
		return entities.MuseumItem{}, err
	}
	return m.CreateMuseumItem.convertToEntity(), nil
}

func (r *Repo) UpdateMuseumItem(item entities.MuseumItem) error {
	m := struct {
		UpdateMuseumItem museumItem `graphql:"updateMuseumItem(id: $id, input: $input)"`
	}{}
	type UpdateMuseumItemInput struct {
		Name            graphql.String  `json:"name"`
		InventoryNumber graphql.String  `json:"inventoryNumber"`
		Annotation      *graphql.String `json:"annotation"`
		CreationDate    time.Time       `json:"creationDate"`
	}
	input := UpdateMuseumItemInput{
		Name:            graphql.String(item.Name),
		InventoryNumber: graphql.String(item.InventoryNumber),
		CreationDate:    item.CreationDate.Time,
	}
	if item.Annotation != "" {
		input.Annotation = graphql.NewString(graphql.String(item.Annotation))
	}
	err := r.client.Mutate(context.TODO(), &m, map[string]interface{}{
		"id":    graphql.ID(strconv.Itoa(item.ID)),
		"input": input,
	})
	return errors.Wrap(err, "failed to mutate items")
}

func (r *Repo) DeleteMuseumItem(id int) error {
	m := struct {
		ID graphql.ID `graphql:"deleteMuseumItem(id: $id)"`
	}{}

	err := r.client.Mutate(context.TODO(), &m, map[string]interface{}{
		"id": graphql.ID(strconv.Itoa(id)),
	})
	return err
}
