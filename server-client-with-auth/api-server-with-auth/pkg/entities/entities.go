package entities

import "time"

type User struct {
	ID       int
	Email    string
	Password string
}

type MuseumFund struct {
	ID        int
	Name      string
	CreatedBy int
}

type MuseumItemMovement struct {
	ID                  int
	MuseumItemID        int
	Item                MuseumItem
	AcceptDate          *time.Time
	ExhibitTransferDate *time.Time
	ExhibitReturnDate   *time.Time
	MuseumExhibitionID  int
	ResponsiblePersonID int
	ResponsiblePerson   Person
	CreatedBy           int
}

type MuseumItem struct {
	ID              int
	InventoryNumber string
	Name            string
	CreationDate    Date
	KeeperID        int
	MuseumSetID     int
	MuseumFundID    int
	Annotation      string
	CreatedBy       int
}

type MuseumItemWithKeeper struct {
	MuseumItem
	Keeper Person
}

type SearchMuseumItemsArgs struct {
	ItemName *string
	SetName  *string
	Date     *time.Time
}

type MuseumItemWithDetails struct {
	MuseumItem
	Keeper Person
	Set    MuseumSet
	Fund   MuseumFund
}

type Person struct {
	ID         int
	FirstName  string
	LastName   string
	MiddleName string
}

type MuseumSet struct {
	ID        int
	Name      string
	CreatedBy int
}

type MuseumSetWithDetails struct {
	MuseumSet
	Items []MuseumItemWithKeeper
}

type Date struct {
	Time time.Time
}

func (d Date) String() string {
	return d.Time.Format("2006-01-02")
}

func NewDate(time time.Time) Date {
	return Date{
		Time: roundToDay(time),
	}
}

func roundToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
