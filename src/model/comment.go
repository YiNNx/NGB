package model

import "time"

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

func AddComment(c *Comment) error {
	if _, err := db.Model(c).Insert(); err != nil {
		return err
	}

	p := &Post{Pid: c.Post}
	if err := db.Model(p).WherePK().Select(); err != nil {
		return err
	}
	p.Comments = append(p.Comments, c.Cid)
	_, err := db.Model(p).Column("comments").WherePK().Update()
	if err != nil {
		return err
	}

	u := &User{Uid: c.From}
	if err := db.Model(u).WherePK().Select(); err != nil {
		return err
	}
	u.Comments = append(u.Comments, c.Post)
	if _, err := db.Model(u).Column("comments").WherePK().Update(); err != nil {
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
