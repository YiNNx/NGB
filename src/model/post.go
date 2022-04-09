package model

import "time"

type Post struct {
	tableName struct{}

	Pid    int       `pg:",pk"`
	Board  int       `pg:",notnull"`
	Time   time.Time `pg:"default:now()"`
	Author int       `pg:",notnull"`
	Tags   []string

	Title   string `pg:",notnull"`
	Content string `pg:",notnull"`

	Comments    []int
	Likes       []int
	Collections []int
}

func AddPost(p *Post) error {
	_, err := db.Model(p).Insert()
	if err != nil {
		return err
	}

	u := &User{Uid: p.Author}
	if err := db.Model(u).WherePK().Select(); err != nil {
		return err
	}
	u.Posts = append(u.Posts, p.Pid)
	if _, err := db.Model(u).Column("posts").WherePK().Update(); err != nil {
		return err
	}

	b := &Board{Bid: p.Board}
	if err := db.Model(b).WherePK().Select(); err != nil {
		return err
	}
	b.Posts = append(b.Posts, p.Pid)
	if _, err := db.Model(b).Column("posts").WherePK().Update(); err != nil {
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

func AddLike(pid int, uid int) error {
	p := &Post{Pid: pid}
	err := db.Model(p).WherePK().Select()
	if err != nil {
		return err
	}
	p.Likes = append(p.Likes, uid)
	_, err = db.Model(p).
		Column("likes").
		WherePK().
		Update()
	if err != nil {
		return err
	}

	u := &User{Uid: uid}
	if err := db.Model(u).WherePK().Select(); err != nil {
		return err
	}
	u.Likes = append(u.Likes, pid)
	if _, err := db.Model(u).Column("likes").WherePK().Update(); err != nil {
		return err
	}
	return nil
}

func AddCollection(pid int, uid int) error {
	p := &Post{Pid: pid}
	err := db.Model(p).WherePK().Select()
	if err != nil {
		return err
	}
	p.Collections = append(p.Collections, uid)
	_, err = db.Model(p).
		Column("collections").
		WherePK().
		Update()
	if err != nil {
		return err
	}

	u := &User{Uid: uid}
	if err := db.Model(u).WherePK().Select(); err != nil {
		return err
	}
	u.Collections = append(u.Collections, pid)
	if _, err := db.Model(u).Column("collections").WherePK().Update(); err != nil {
		return err
	}
	return nil
}
