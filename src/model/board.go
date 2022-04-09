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

	Managers []int
	Members  []int
	Posts    []int
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

func JoinBoard(uid int, bid int) error {
	b := &Board{Bid: bid}
	if err := db.Model(b).WherePK().Select(); err != nil {
		return err
	}
	b.Members = append(b.Members, uid)
	if _, err := db.Model(b).Column("members").WherePK().Update(); err != nil {
		return err
	}

	u := &User{Uid: uid}
	if err := db.Model(u).WherePK().Select(); err != nil {
		return err
	}
	u.BoardsJoin = append(u.BoardsJoin, bid)
	if _, err := db.Model(u).Column("boards_join").WherePK().Update(); err != nil {
		return err
	}
	return nil
}

func GetPostsFromBoard(b Board) error {
	posts := b.Posts
	err := db.Model((*Post)(nil)).
		Where("pid in (?)", pg.In(posts)).
		Select()
	if err != nil {
		return err
	}
	return nil
}
