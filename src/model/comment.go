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

func InsertComment(c *Comment) error {
	if _, err := db.Model(c).Insert(); err != nil {
		return err
	}
	return nil
}

func GetCommentsCountOfPost(pid int) (int, error) {
	count, err := db.Model((*Comment)(nil)).Where("post = ?", pid).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetCommentByCid(cid int) (*Comment, error) {
	c := &Comment{Cid: cid}
	err := db.Model(c).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetCommentsByPid(pid int) ([]Comment, error) {
	var comments []Comment
	err := db.Model(&comments).
		Where("post = ?", pid).
		Select()
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func CheckCommentId(cid int) error {
	c := &Comment{Cid: cid}
	err := db.Model(c).WherePK().Select()
	if err != nil {
		return err
	}
	return nil
}
