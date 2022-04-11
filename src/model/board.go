package model

import (
	"github.com/go-pg/pg/v10"
	"time"
)

type Board struct {
	tableName struct{}

	Bid    int    `pg:",pk"`
	Name   string `pg:",notnull"`
	Avatar string
	Time   time.Time `pg:"default:now()"`
	Intro  string

	Managers []User `pg:",many2many:manages"`
	Members  []User `pg:",many2many:joins"`
	//Posts    []Post `pg:"rel:has-many"`
}

func InsertBoard(b *Board) error {
	_, err := db.Model(b).Insert()
	if err != nil {
		return err
	}
	return nil
}

func GetBoardByBid(bid int) (*Board, error) {
	b := &Board{Bid: bid}
	err := db.Model(b).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GetBoardByBids(bids []int) ([]Board, error) {
	var boards []Board
	err := db.Model(&boards).
		Where("rid in (?)", pg.In(bids)).
		Select()
	if err != nil {
		return nil, err
	}
	return boards, nil
}
