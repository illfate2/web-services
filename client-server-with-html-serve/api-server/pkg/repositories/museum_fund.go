package repo

import (
	"context"

	"github.com/illfate2/web-services/client-server-with-html-serve/api-server/pkg/entities"
)

func (r *Repo) InsertMuseumFund(fund entities.MuseumFund) (entities.MuseumFund, error) {
	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO museum_funds(name)
			VALUES($1)
			RETURNING id`,
		fund.Name).
		Scan(&fund.ID)
	if err != nil {
		return entities.MuseumFund{}, err
	}
	return fund, nil
}

func (r *Repo) DeleteMuseumFund(id int) error {
	_, err := r.conn.Exec(context.TODO(), `DELETE FROM museum_funds WHERE id = $1`, id)
	return err
}

func (r *Repo) FindMuseumFunds() ([]entities.MuseumFund, error) {
	rows, err := r.conn.Query(context.Background(),
		`SELECT id, name FROM museum_funds`)
	if err != nil {
		return nil, err
	}
	var funds []entities.MuseumFund
	for rows.Next() {
		var fund entities.MuseumFund
		err := rows.Scan(&fund.ID, &fund.Name)
		if err != nil {
			return nil, err
		}
		funds = append(funds, fund)
	}
	return funds, err
}

func (r *Repo) FindMuseumFundByID(id int) (entities.MuseumFund, error) {
	var fund entities.MuseumFund
	err := r.conn.QueryRow(context.Background(),
		`SELECT id, name FROM museum_funds WHERE id = $1`,
		id).
		Scan(&fund.ID, &fund.Name)
	return fund, err
}

func (r *Repo) UpdateMuseumFund(fund entities.MuseumFund) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE museum_funds 
			SET name = $1
			WHERE id = $2`,
		fund.Name, fund.ID)
	return err
}
