package model

import (
	"errors"
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

type Like struct {
	UserUid int
	PostPid int
}

type Collection struct {
	UserUid int
	PostPid int
}

// ----------- Post -------------

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

func GetPostsByPidList(pids []int) ([]Post, error) {
	if pids == nil {
		return nil, nil
	}
	var posts []Post
	err := tx.Model(&posts).
		Where("pid in (?)", pg.In(pids)).
		Select()
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

// -------- Like & Collection -----------

func InsertLikeOrCollection(m interface{}, pid int, uid int) error {
	if err := CheckPK(&Post{Pid: pid}); err != nil {
		return err
	}

	if res, err := tx.Model(m).Where("user_uid = ?", uid).Where("post_pid = ?", pid).SelectOrInsert(); err != nil {
		return err
	} else if res == false {
		return errors.New("already exist")
	}
	return nil
}

func DeleteLikeOrCollection(m interface{}, postPid int, userUid int) error {
	_, err := tx.Model(m).Where("post_pid = ?", postPid).Where("user_uid = ?", userUid).Delete()
	if err != nil {
		return err
	}
	return nil
}

func GetLikesOfUser(uid int) (posts []Post, count int, err error) {
	var user User
	count, err = tx.Model(&user).Relation("Likes").Where("uid = ?", uid).SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return user.Likes, count, nil
}

func GetLikesOfPost(pid int) (users []User, count int, err error) {
	var post Post
	count, err = tx.Model(&post).Relation("Likes").Where("pid = ?", pid).SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return post.Likes, count, nil
}

func GetCollectionsOfUser(uid int) ([]Post, error) {
	var user User
	err := tx.Model(&user).Relation("Collections").Where("uid = ?", uid).Select()
	if err != nil {
		return nil, err
	}
	return user.Collections, nil
}

func GetCollectionsOfPost(pid int) ([]User, error) {
	var post Post
	err := tx.Model(&post).Relation("Collections").Where("pid = ?", pid).Select()
	if err != nil {
		return nil, err
	}
	return post.Collections, nil
}
