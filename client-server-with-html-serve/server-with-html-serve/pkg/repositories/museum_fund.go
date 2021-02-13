package repo

import (
	"context"
	"strconv"

	"github.com/shurcooL/graphql"

	"github.com/illfate2/web-services/client-server-with-html-serve/server-with-html-serve/pkg/entities"
)

type museumFund struct {
	ID   graphql.ID     `json:"id"`
	Name graphql.String `json:"name"`
}

func (f museumFund) convertToEntity() entities.MuseumFund {
	id, _ := strconv.Atoi(f.ID.(string))
	return entities.MuseumFund{
		ID:   id,
		Name: string(f.Name),
	}
}

func (r *Repo) InsertMuseumFund(fund entities.MuseumFund) (entities.MuseumFund, error) {
	m := struct {
		CreateMuseumFund museumFund `graphql:"createMuseumFund(input: $input)"`
	}{}
	type MuseumFundInput struct {
		Name graphql.String `json:"name"`
	}
	input := MuseumFundInput{
		Name: graphql.String(fund.Name),
	}
	err := r.client.Mutate(context.TODO(), &m, map[string]interface{}{
		"input": input,
	})
	if err != nil {
		return entities.MuseumFund{}, err
	}
	return m.CreateMuseumFund.convertToEntity(), nil
}

func (r *Repo) DeleteMuseumFund(id int) error {
	m := struct {
		ID graphql.ID `graphql:"deleteMuseumFund(id: $id)"`
	}{}

	err := r.client.Mutate(context.TODO(), &m, map[string]interface{}{
		"id": graphql.ID(strconv.Itoa(id)),
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) FindMuseumFunds() ([]entities.MuseumFund, error) {
	q := struct {
		MuseumFunds []museumFund `graphql:"museumFunds"`
	}{}
	err := r.client.Query(context.TODO(), &q, nil)
	if err != nil {
		return nil, err
	}
	res := make([]entities.MuseumFund, 0, len(q.MuseumFunds))
	for _, s := range q.MuseumFunds {
		res = append(res, s.convertToEntity())
	}
	return res, nil
}

func (r *Repo) FindMuseumFundByID(id int) (entities.MuseumFund, error) {
	q := struct {
		MuseumFund museumFund `graphql:"museumFund(id: $id)"`
	}{}
	err := r.client.Query(context.TODO(), &q, map[string]interface{}{
		"id": graphql.ID(strconv.Itoa(id)),
	})
	if err != nil {
		return entities.MuseumFund{}, err
	}
	return q.MuseumFund.convertToEntity(), nil
}

func (r *Repo) UpdateMuseumFund(fund entities.MuseumFund) error {
	m := struct {
		UpdateMuseumFund museumFund `graphql:"updateMuseumFund(id: $id, input: $input)"`
	}{}
	type UpdateMuseumFundInput struct {
		Name graphql.String `json:"name"`
	}
	input := UpdateMuseumFundInput{
		Name: graphql.String(fund.Name),
	}
	err := r.client.Mutate(context.TODO(), &m, map[string]interface{}{
		"id":    graphql.ID(strconv.Itoa(fund.ID)),
		"input": input,
	})
	return err
}
