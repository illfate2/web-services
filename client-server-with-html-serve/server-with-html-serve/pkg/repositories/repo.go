package repo

import (
	"context"
	"strconv"

	"github.com/shurcooL/graphql"
)

func New(client *graphql.Client) *Repo {
	return &Repo{
		client: client,
	}
}

type Repo struct {
	client *graphql.Client
}

func (r *Repo) findByID(id int, q interface{}) error {
	err := r.client.Query(context.TODO(), q, map[string]interface{}{
		"id": graphql.ID(strconv.Itoa(id)),
	})
	if err != nil {
		return err
	}
	return nil
}
