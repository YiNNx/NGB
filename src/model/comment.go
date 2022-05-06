package model

import (
	"time"
)

type Comment struct {
	tableName struct{}

	Cid       int `pg:",pk"`
	Post      int `pg:",notnull"`
	IsAuthor  bool
	ParentCid int

	Time    time.Time `pg:"default:now()"`
	From    int       `pg:",notnull"`
	To      int       `pg:",notnull"`
	Content string    `pg:",notnull"`
}

func GetCommentsOfPost(pid int) (comments []Comment, count int, err error) {
	count, err = tx.Model(&comments).
		Where("post = ?", pid).
		SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return comments, count, nil
}
