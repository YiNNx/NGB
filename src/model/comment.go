package model

import (
	"time"
)

type Comment struct {
	tableName struct{}

	Cid      int `pg:",pk"`
	Post     int `pg:",notnull"`
	IsAuthor bool
	SubCid   []int

	Time    time.Time `pg:"default:now()"`
	From    int       `pg:",notnull"`
	To      int       `pg:",notnull"`
	Content string    `pg:",notnull"`
}

func InsertComment(c *Comment) error {
	if _, err := db.Model(c).Insert(); err != nil {
		return err
	}
	return nil
}

func GetCommentByCid(cid int) (*Comment, error) {
	c := &Comment{Cid: cid}
	err := db.Model(c).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetCommentsByPid(pid int) ([]Post, error) {
	var posts []Post
	err := db.Model(&posts).
		Where("pid = ?", pid).
		Select()
	if err != nil {
		return nil, err
	}
	return posts, nil
}
