package repo

import (
	"context"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/shurcooL/graphql"

	"github.com/illfate2/web-services/client-server-with-html-serve/server-with-html-serve/pkg/entities"
)

type museumMovement struct {
	ID                  graphql.ID `json:"id"`
	AcceptDate          *time.Time `json:"acceptDate"`
	ExhibitTransferDate *time.Time `json:"exhibitTransferDate"`
	ExhibitReturnDate   *time.Time `json:"exhibitReturnDate"`
}

func (m museumMovement) convertToEntity() entities.MuseumItemMovement {
	id, _ := strconv.Atoi(m.ID.(string))
	return entities.MuseumItemMovement{
		ID:                  id,
		AcceptDate:          m.AcceptDate,
		ExhibitTransferDate: m.ExhibitTransferDate,
		ExhibitReturnDate:   m.ExhibitReturnDate,
	}
}

type museumMovementWithDetails struct {
	ID                  graphql.ID `json:"id"`
	AcceptDate          *time.Time `json:"acceptDate"`
	ExhibitTransferDate *time.Time `json:"exhibitTransferDate"`
	ExhibitReturnDate   *time.Time `json:"exhibitReturnDate"`
	Item                museumItem `json:"item"`
	Person              person     `json:"person"`
}

func (m museumMovementWithDetails) convertToEntity() entities.MuseumItemMovement {
	id, _ := strconv.Atoi(m.ID.(string))
	return entities.MuseumItemMovement{
		ID:                  id,
		Item:                m.Item.convertToEntity(),
		AcceptDate:          m.AcceptDate,
		ExhibitTransferDate: m.ExhibitTransferDate,
		ExhibitReturnDate:   m.ExhibitReturnDate,
		ResponsiblePerson: entities.Person{
			FirstName:  m.Person.FirstName,
			LastName:   m.Person.LastName,
			MiddleName: m.Person.MiddleName,
		},
	}
}

func (r *Repo) FindMuseumItemMovement(id int) (entities.MuseumItemMovement, error) {
	q := struct {
		MuseumMovement museumMovementWithDetails `graphql:"museumMovement(id: $id)"`
	}{}
	err := r.findByID(id, &q)
	if err != nil {
		return entities.MuseumItemMovement{}, err
	}
	return q.MuseumMovement.convertToEntity(), nil
}

func (r *Repo) FindMuseumItemMovements() ([]entities.MuseumItemMovement, error) {
	q := struct {
		MuseumMovements []museumMovement `graphql:"museumMovements"`
	}{}
	err := r.client.Query(context.TODO(), &q, nil)
	if err != nil {
		return nil, err
	}
	res := make([]entities.MuseumItemMovement, 0, len(q.MuseumMovements))
	for _, s := range q.MuseumMovements {
		res = append(res, s.convertToEntity())
	}
	return res, nil
}

func (r *Repo) InsertMuseumItemMovement(movement entities.MuseumItemMovement) (entities.MuseumItemMovement, error) {
	m := struct {
		CreateMuseumMovement museumMovement `graphql:"createMuseumItemMovement(input: $input)"`
	}{}
	type MuseumMovementInput struct {
		ItemID              graphql.ID  `json:"itemID,omitempty"`
		AcceptDate          *time.Time  `json:"acceptDate,omitempty"`
		ExhibitTransferDate *time.Time  `json:"exhibitTransferDate,omitempty"`
		ExhibitReturnDate   *time.Time  `json:"exhibitReturnDate,omitempty"`
		PersonInput         personInput `json:"personInput"`
	}
	input := MuseumMovementInput{
		ItemID:              graphqlIDFromInt(movement.MuseumItemID),
		AcceptDate:          movement.AcceptDate,
		ExhibitTransferDate: movement.ExhibitTransferDate,
		ExhibitReturnDate:   movement.ExhibitReturnDate,
		PersonInput: personInput{
			FirstName:  graphql.String(movement.ResponsiblePerson.FirstName),
			MiddleName: graphql.String(movement.ResponsiblePerson.MiddleName),
			LastName:   graphql.String(movement.ResponsiblePerson.LastName),
		},
	}
	err := r.client.Mutate(context.TODO(), &m, map[string]interface{}{
		"input": input,
	})
	if err != nil {
		return entities.MuseumItemMovement{}, err
	}
	return m.CreateMuseumMovement.convertToEntity(), nil
}

func (r *Repo) UpdateMuseumItemMovement(movement entities.MuseumItemMovement) error {
	m := struct {
		CreateMuseumMovement museumMovement `graphql:"createMuseumItemMovement(input: $input)"`
	}{}
	type MuseumMovementInput struct {
		AcceptDate          *time.Time `json:"acceptDate,omitempty"`
		ExhibitTransferDate *time.Time `json:"exhibitTransferDate,omitempty"`
		ExhibitReturnDate   *time.Time `json:"exhibitReturnDate,omitempty"`
	}
	input := MuseumMovementInput{
		AcceptDate:          movement.AcceptDate,
		ExhibitTransferDate: movement.ExhibitTransferDate,
		ExhibitReturnDate:   movement.ExhibitReturnDate,
	}
	err := r.client.Mutate(context.TODO(), &m, map[string]interface{}{
		"id":    graphql.ID(strconv.Itoa(movement.ID)),
		"input": input,
	})
	return errors.Wrap(err, "failed to mutate item movements")
}

func (r *Repo) DeleteMuseumItemMovement(id int) error {
	m := struct {
		ID graphql.ID `graphql:"deleteMuseumItemMovement(id: $id)"`
	}{}

	err := r.client.Mutate(context.TODO(), &m, map[string]interface{}{
		"id": graphql.ID(strconv.Itoa(id)),
	})
	return err
}
