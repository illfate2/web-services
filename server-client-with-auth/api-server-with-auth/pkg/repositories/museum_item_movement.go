package repo

import (
	"context"

	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/entities"
)

func (r *Repo) FindMuseumItemMovement(id int) (entities.MuseumItemMovement, error) {
	var m entities.MuseumItemMovement
	err := r.conn.QueryRow(context.Background(),
		`SELECT mim.id, mim.item_id, mim.responsible_person_id, mim.accept_date,
		mim.exhibit_transfer_date, mim.exhibit_return_date, mi.name, p.first_name, p.middle_name, p.second_name
		FROM museum_item_movements mim 
		LEFT JOIN museum_items mi ON mi.id = mim.item_id
		LEFT JOIN persons p ON p.id = mim.responsible_person_id
		WHERE mim.id = $1`, id).
		Scan(&m.ID, &m.MuseumItemID, &m.ResponsiblePersonID, &m.AcceptDate, &m.ExhibitTransferDate, &m.ExhibitReturnDate,
			&m.Item.Name, &m.ResponsiblePerson.FirstName, &m.ResponsiblePerson.MiddleName, &m.ResponsiblePerson.LastName)
	if err != nil {
		return entities.MuseumItemMovement{}, err
	}
	return m, nil
}

func (r *Repo) FindMuseumItemMovements() ([]entities.MuseumItemMovement, error) {
	var movements []entities.MuseumItemMovement
	rows, err := r.conn.Query(context.Background(),
		`SELECT id, item_id, responsible_person_id, accept_date, exhibit_transfer_date, exhibit_return_date
			FROM museum_item_movements`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var m entities.MuseumItemMovement
		err := rows.Scan(&m.ID, &m.MuseumItemID, &m.ResponsiblePersonID, &m.AcceptDate, &m.ExhibitTransferDate, &m.ExhibitReturnDate)
		if err != nil {
			return nil, err
		}
		movements = append(movements, m)
	}
	return movements, nil
}

func (r *Repo) InsertMuseumItemMovement(movement entities.MuseumItemMovement) (entities.MuseumItemMovement, error) {
	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO museum_item_movements(item_id,responsible_person_id,accept_date,exhibit_transfer_date,exhibit_return_date)
		VALUES($1,$2,$3,$4,$5) RETURNING id`,
		movement.MuseumItemID, movement.ResponsiblePersonID, movement.AcceptDate, movement.ExhibitTransferDate,
		movement.ExhibitReturnDate).
		Scan(&movement.ID)
	if err != nil {
		return entities.MuseumItemMovement{}, err
	}
	return movement, nil
}

func (r *Repo) UpdateMuseumItemMovement(movement entities.MuseumItemMovement) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE museum_item_movements 
			SET accept_date = $1, exhibit_transfer_date = $2, exhibit_return_date= $3
			WHERE id = $4`,
		movement.AcceptDate, movement.ExhibitTransferDate, movement.ExhibitReturnDate, movement.ID)
	return err
}

func (r *Repo) DeleteMuseumItemMovement(id int) error {
	_, err := r.conn.Exec(context.Background(),
		`DELETE FROM museum_item_movements WHERE id = $1`, id)
	return err
}
