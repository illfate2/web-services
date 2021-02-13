package repo

import (
	"context"
	"strconv"

	"github.com/shurcooL/graphql"

	"github.com/illfate2/web-services/client-server-with-html-serve/server-with-html-serve/pkg/entities"
)

type museumSet struct {
	ID   graphql.ID     `json:"id"`
	Name graphql.String `json:"name"`
}

func (s museumSet) convertToEntity() entities.MuseumSet {
	id, _ := strconv.Atoi(s.ID.(string))
	return entities.MuseumSet{
		ID:   id,
		Name: string(s.Name),
	}
}

func (r *Repo) InsertMuseumSet(set entities.MuseumSet) (entities.MuseumSet, error) {
	m := struct {
		CreateMuseumSet museumSet `graphql:"createMuseumSet(input: $input)"`
	}{}
	type MuseumSetInput struct {
		Name graphql.String `json:"name"`
	}
	input := MuseumSetInput{
		Name: graphql.String(set.Name),
	}
	err := r.client.Mutate(context.TODO(), &m, map[string]interface{}{
		"input": input,
	})
	if err != nil {
		return entities.MuseumSet{}, err
	}
	return m.CreateMuseumSet.convertToEntity(), nil
}

func (r *Repo) FindMuseumSets() ([]entities.MuseumSet, error) {
	q := struct {
		MuseumSets []museumSet `graphql:"museumSets"`
	}{}
	err := r.client.Query(context.TODO(), &q, nil)
	if err != nil {
		return nil, err
	}
	res := make([]entities.MuseumSet, 0, len(q.MuseumSets))
	for _, s := range q.MuseumSets {
		res = append(res, s.convertToEntity())
	}
	return res, nil
}

func (r *Repo) FindMuseumSetByName(name string) (entities.MuseumSet, error) {
	return entities.MuseumSet{}, nil
}

func (r *Repo) FindMuseumSet(id int) (entities.MuseumSetWithDetails, error) {
	q := struct {
		MuseumSet museumSet `graphql:"museumSet(id: $id)"`
	}{}
	err := r.findByID(id, &q)
	if err != nil {
		return entities.MuseumSetWithDetails{}, err
	}
	return entities.MuseumSetWithDetails{
		MuseumSet: q.MuseumSet.convertToEntity(),
	}, nil
}

func (r *Repo) UpdateMuseumSet(set entities.MuseumSet) error {
	m := struct {
		UpdateMuseumSet museumSet `graphql:"updateMuseumSet(id: $id, input: $input)"`
	}{}
	type MuseumSetInput struct {
		Name graphql.String `json:"name"`
	}
	input := MuseumSetInput{
		Name: graphql.String(set.Name),
	}
	err := r.client.Mutate(context.TODO(), &m, map[string]interface{}{
		"id":    graphql.ID(strconv.Itoa(set.ID)),
		"input": input,
	})
	return err
}

func (r *Repo) DeleteMuseumSet(id int) error {
	m := struct {
		ID graphql.ID `graphql:"deleteMuseumSet(id: $id)"`
	}{}

	err := r.client.Mutate(context.TODO(), &m, map[string]interface{}{
		"id": graphql.ID(strconv.Itoa(id)),
	})
	return err
}
