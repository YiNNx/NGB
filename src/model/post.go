package model

import (
	"github.com/go-pg/pg/v10"
	"time"
)

type Post struct {
	tableName struct{}

	Pid    int       `pg:",pk"`
	Board  int       `pg:",notnull"`
	Time   time.Time `pg:"default:now()"`
	Author int       `pg:",notnull"`
	Tags   []string

	Title   string `pg:",notnull"`
	Content string `pg:",notnull"`

	//Comments    []Comment `pg:"rel:has-many"`
	Likes       []User `pg:"many2many:likes"`
	Collections []User `pg:"many2many:collections"`
}

func InsertPost(p *Post) error {
	_, err := db.Model(p).Insert()
	if err != nil {
		return err
	}
	return nil
}

func GetPostByPid(pid int) (*Post, error) {
	p := &Post{Pid: pid}
	err := db.Model(p).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GetPostsByTag(tag string) ([]Post, error) {
	var posts []Post
	err := db.Model(&posts).
		Where("tags::jsonb ?& array['" + tag + "']").
		Select()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func GetPostsByPids(pids []int) ([]Post, error) {
	var posts []Post
	err := db.Model(&posts).
		Where("pid in (?)", pg.In(pids)).
		Select()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func GetPostsByUid(uid int) ([]Post, error) {
	var posts []Post
	err := db.Model(&posts).
		Where("author = ?", uid).
		Select()
	if err != nil {
		return nil, err
	}
	return posts, nil
}
