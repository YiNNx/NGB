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

	Likes       []User `pg:"many2many:likes"`
	Collections []User `pg:"many2many:collections"`
}

func InsertPost(p *Post) error {
	_, err := tx.Model(p).Insert()
	if err != nil {
		return err
	}
	return nil
}

func GetPostByPid(pid int) (*Post, error) {
	p := &Post{Pid: pid}
	err := tx.Model(p).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GetPostsByTag(tag string, limit int, offset int) ([]Post, error) {
	var posts []Post
	err := tx.Model(&posts).
		Where("tags::jsonb ?& array['" + tag + "']").
		Limit(limit).Offset(offset).Select()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func GetPostsByPids(pids []int) ([]Post, error) {
	if pids != nil {
		var posts []Post
		err := tx.Model(&posts).
			Where("pid in (?)", pg.In(pids)).
			Select()
		if err != nil {
			return nil, err
		}
		return posts, nil
	} else {
		return nil, nil
	}
}

func GetPostsByUid(uid int) ([]Post, error) {
	var posts []Post
	err := tx.Model(&posts).
		Where("author = ?", uid).
		Select()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func GetPostsByBoard(bid int, limit int, offset int) ([]Post, error) {
	var posts []Post
	err := tx.Model(&posts).
		Where("board = ?", bid).
		Limit(limit).Offset(offset).Select()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func SelectAllPosts(limit int, offset int) ([]Post, error) {
	var posts []Post
	err := tx.Model(&posts).Limit(limit).Offset(offset).Select()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func CheckPostId(pid int) error {
	p := &Post{Pid: pid}
	err := tx.Model(p).WherePK().Select()
	if err != nil {
		return err
	}
	return nil
}

func GetBoardByPost(pid int) (int, error) {
	p := &Post{Pid: pid}
	err := tx.Model(p).WherePK().Select()
	if err != nil {
		return 0, err
	}
	return p.Board, nil
}

func DeletePost(pid int) error {
	p := &Post{Pid: pid}
	_, err := tx.Model(p).WherePK().Delete()
	if err != nil {
		return err
	}
	return nil
}
