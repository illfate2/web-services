// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type MuseumExhibition struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	PhoneNo string `json:"phoneNo"`
}

type MuseumFund struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MuseumFundInput struct {
	Name string `json:"name"`
}

type MuseumItem struct {
	ID              int         `json:"id"`
	InventoryNumber string      `json:"inventoryNumber"`
	Name            string      `json:"name"`
	CreationDate    time.Time   `json:"creationDate"`
	Annotation      *string     `json:"annotation"`
	Person          *Person     `json:"person"`
	Set             *MuseumSet  `json:"set"`
	Fund            *MuseumFund `json:"fund"`
}

type MuseumItemInput struct {
	InventoryNumber string       `json:"inventoryNumber"`
	Name            string       `json:"name"`
	Annotation      *string      `json:"annotation"`
	CreationDate    time.Time    `json:"creationDate"`
	SetID           int          `json:"setID"`
	FundID          int          `json:"fundID"`
	PersonInput     *PersonInput `json:"personInput"`
}

type MuseumItemMovement struct {
	ID                  int         `json:"id"`
	AcceptDate          *time.Time  `json:"acceptDate"`
	ExhibitTransferDate *time.Time  `json:"exhibitTransferDate"`
	ExhibitReturnDate   *time.Time  `json:"exhibitReturnDate"`
	Item                *MuseumItem `json:"item"`
	Person              *Person     `json:"person"`
}

type MuseumMovementInput struct {
	AcceptDate          *time.Time   `json:"acceptDate"`
	ItemID              int          `json:"itemID"`
	ExhibitTransferDate *time.Time   `json:"exhibitTransferDate"`
	ExhibitReturnDate   *time.Time   `json:"exhibitReturnDate"`
	PersonInput         *PersonInput `json:"personInput"`
}

type MuseumSet struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MuseumSetInput struct {
	Name string `json:"name"`
}

type Person struct {
	ID         *int   `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	MiddleName string `json:"middleName"`
}

type PersonInput struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	MiddleName string `json:"middleName"`
}

type Signup struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateMuseumFundInput struct {
	Name *string `json:"name"`
}

type UpdateMuseumItemInput struct {
	InventoryNumber *string    `json:"inventoryNumber"`
	Name            *string    `json:"name"`
	Annotation      *string    `json:"annotation"`
	SetID           *int       `json:"setID"`
	FundID          *int       `json:"fundID"`
	CreationDate    *time.Time `json:"creationDate"`
}

type UpdateMuseumMovementInput struct {
	Name                *string    `json:"name"`
	AcceptDate          *time.Time `json:"acceptDate"`
	ExhibitTransferDate *time.Time `json:"exhibitTransferDate"`
	ExhibitReturnDate   *time.Time `json:"exhibitReturnDate"`
}

type UpdateMuseumSetInput struct {
	Name string `json:"name"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
