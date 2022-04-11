package model

import "github.com/go-pg/pg/v10/orm"

type Like struct {
	Pid int `pg:",notnull"`
	Uid int `pg:",notnull"`
}

type Collection struct {
	Pid int `pg:",notnull"`
	Uid int `pg:",notnull"`
}

type FollowShip struct {
	tableName struct{} `pg:"follow_ships"`
	Followee  int      `pg:",notnull"`
	Follower  int      `pg:",notnull"`
}

type JoinShip struct {
	Bid int `pg:",notnull"`
	Uid int `pg:",notnull"`
}

type ManageShip struct {
	Bid int `pg:",notnull"`
	Uid int `pg:",notnull"`
}

func InsertLike(pid int, uid int) error {
	l := &Like{
		Pid: pid,
		Uid: uid,
	}
	_, err := db.Model(l).Insert()
	if err != nil {
		return err
	}
	return nil
}

//func InsertCollection(pid int, uid int) error {
//	f := &Collection{
//		Pid: pid,
//		Uid: uid,
//	}
//	_, err := db.Model(f).Insert()
//	if err != nil {
//		return err
//	}
//	return nil
//}

func InsertFollowShip(followee int, follower int) error {
	f := FollowShip{
		Followee: followee,
		Follower: follower,
	}
	_, err := db.Model(f).Insert()
	if err != nil {
		return err
	}
	return nil
}

func InsertJoinShip(bid int, uid int) error {
	f := &JoinShip{
		Bid: bid,
		Uid: uid,
	}
	_, err := db.Model(f).Insert()
	if err != nil {
		return err
	}
	return nil
}

func InsertManageShip(bid int, uid int) error {
	f := ManageShip{
		Bid: bid,
		Uid: uid,
	}
	_, err := db.Model(f).Insert()
	if err != nil {
		return err
	}
	return nil
}

//func GetLikesByUid(uid int) ([]Post, error) {
//	var likes []Like
//	if err := db.Model(&likes).Where("uid = ?", uid).Select(); err != nil {
//		return nil, err
//	}
//	var pids []int
//	for i, _ := range likes {
//		pids = append(pids, likes[i].Pid)
//	}
//	posts, err := GetPostsByPids(pids)
//	if err != nil {
//		return nil, err
//	}
//	return posts, nil
//}

func GetLikesByPid(pid int) ([]Post, error) {
	var likes []Like
	if err := db.Model(&likes).Where("pid = ?", pid).Select(); err != nil {
		return nil, err
	}
	var uids []int
	for i, _ := range likes {
		uids = append(uids, likes[i].Uid)
	}
	posts, err := GetPostsByUids(uids)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

//SELECT posts.*
//FROM likes
//inner join posts on likes.pid = posts.pid where uid=2 order by time ;

func GetLikesByUid(uid int) ([]*Post, error) {
	orm.RegisterTable((*Like)(nil))
	var posts []*Post
	err := db.Model(&posts).Column("posts.*").Join("inner join likes on likes.pid = posts.pid").Where("likes.uid = ?", uid).Select()
	if err != nil {
		return nil, err
	}
	return posts, nil
}
