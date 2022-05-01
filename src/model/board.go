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

	//Posts    []Post `pg:"rel:has-many"`
}

func InsertBoard(b *Board) error {
	tx, err := db.Begin()
	defer tx.Close()
	_, err = db.Model(b).Insert()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
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

func GetBoardsByBids(bids []int) ([]Board, error) {
	if bids != nil {
		var boards []Board
		err := db.Model(&boards).
			Where("bid in (?)", pg.In(bids)).
			Select()
		if err != nil {
			return nil, err
		}
		return boards, nil
	} else {
		return nil, nil
	}
}

func SelectAllBoards() ([]Board, error) {
	var boards []Board
	err := db.Model(&boards).Select()
	if err != nil {
		return nil, err
	}
	return boards, nil
}

func CheckBoardId(bid int) error {
	b := &Board{Bid: bid}
	err := db.Model(b).WherePK().Select()
	if err != nil {
		return err
	}
	return nil
}
