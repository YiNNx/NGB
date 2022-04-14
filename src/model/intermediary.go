package model

import "errors"

//点赞

type Like struct {
	UserUid int
	PostPid int
}

func InsertLike(pid int, uid int) error {
	var li []Like
	if db.Model(&li).Where("user_uid = ?", uid).Where("post_pid = ?", pid).Select(); li != nil {
		return errors.New("already liked")
	}

	if err := CheckUserId(uid); err != nil {
		return err
	}
	if err := CheckPostId(pid); err != nil {
		return err
	}

	l := &Like{
		PostPid: pid,
		UserUid: uid,
	}

	_, err := db.Model(l).Insert()
	if err != nil {
		return err
	}
	return nil
}

func DeleteLike(postPid int, userUid int) error {
	like := &Like{}
	_, err := db.Model(like).Where("post_pid = ?", postPid).Where("user_uid = ?", userUid).Delete()
	if err != nil {
		return err
	}

	return nil
}

func GetLikesOfUser(uid int) ([]Post, error) {
	var user User
	err := db.Model(&user).Relation("Likes").Where("uid = ?", uid).Select()
	if err != nil {
		return nil, err
	}
	return user.Likes, nil
}

func GetLikesOfPost(pid int) ([]User, error) {
	var post Post
	err := db.Model(&post).Relation("Likes").Where("post_pid = ?", pid).Select()
	if err != nil {
		return nil, err
	}
	return post.Likes, nil
}

func GetLikesCountOfPost(pid int) (int, error) {
	count, err := db.Model((*Like)(nil)).Where("post_pid = ?", pid).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

//收藏

type Collection struct {
	UserUid int
	PostPid int
}

func InsertCollection(pid int, uid int) error {
	var co []Collection
	if db.Model(&co).Where("user_uid = ?", uid).Where("post_pid = ?", pid).Select(); co != nil {
		return errors.New("already collected")
	}

	if err := CheckUserId(uid); err != nil {
		return err
	}
	if err := CheckPostId(pid); err != nil {
		return err
	}
	f := &Collection{
		PostPid: pid,
		UserUid: uid,
	}
	_, err := db.Model(f).Insert()
	if err != nil {
		return err
	}
	return nil
}

func DeleteCollection(postPid int, userUid int) error {
	collection := &Collection{}
	_, err := db.Model(collection).Where("post_pid = ?", postPid).Where("user_uid = ?", userUid).Delete()
	if err != nil {
		return err
	}
	return nil
}

func GetCollectionsOfUser(uid int) ([]Post, error) {
	var user User
	err := db.Model(&user).Relation("Collections").Where("uid = ?", uid).Select()
	if err != nil {
		return nil, err
	}
	return user.Collections, nil
}

func GetCollectionsOfPost(pid int) ([]User, error) {
	var post Post
	err := db.Model(&post).Relation("Collections").Where("pid = ?", pid).Select()
	if err != nil {
		return nil, err
	}
	return post.Collections, nil
}

//加入板块

type JoinShip struct {
	Uid int
	Bid int
}

func InsertJoinShip(bid int, uid int) error {
	var jo []JoinShip
	if db.Model(&jo).Where("uid = ?", uid).Where("uid = ?", uid).Select(); jo != nil {
		return errors.New("already liked")
	}

	if err := CheckUserId(uid); err != nil {
		return err
	}
	if err := CheckBoardId(bid); err != nil {
		return err
	}
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

func DeleteJoinShip(bid int, uid int) error {
	joinShip := &JoinShip{}
	_, err := db.Model(joinShip).Where("bid = ?", bid).Where("uid = ?", uid).Delete()
	if err != nil {
		return err
	}
	return nil
}

func GetBoardsOfUser(uid int) ([]Board, error) {
	var joins []JoinShip
	if err := db.Model(&joins).Where("uid = ?", uid).Select(); err != nil {
		return nil, err
	}
	var bids []int
	for i, _ := range joins {
		bids = append(bids, joins[i].Bid)
	}
	boards, err := GetBoardsByBids(bids)
	if err != nil {
		return nil, err
	}
	return boards, nil
}

func GetMembersOfBoard(bid int) ([]User, error) {
	var joins []JoinShip
	if err := db.Model(&joins).Where("bid = ?", bid).Select(); err != nil {
		return nil, err
	}
	var uids []int
	for i, _ := range joins {
		uids = append(uids, joins[i].Uid)
	}
	users, err := GetUsersByUids(uids)
	if err != nil {
		return nil, err
	}
	return users, nil
}

//管理板块

type ManageShip struct {
	Bid int
	Uid int
}

func InsertManageShip(bid int, uid int) error {
	var fo []ManageShip
	if db.Model(&fo).Where("uid = ?", uid).Where("bid = ?", bid).Select(); fo != nil {
		return errors.New("already managed")
	}

	if err := CheckUserId(uid); err != nil {
		return err
	}
	if err := CheckBoardId(bid); err != nil {
		return err
	}

	f := &ManageShip{
		Bid: bid,
		Uid: uid,
	}
	_, err := db.Model(f).Insert()
	if err != nil {
		return err
	}
	return nil
}

func GetBoardsMngOfUser(uid int) ([]Board, error) {
	var user User
	err := db.Model(&user).Relation("Manages").Where("uid = ?", uid).Select()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func GetManagersOfBoard(bid int) ([]User, error) {
	var b Board
	err := db.Model(&b).Relation("Manages").Where("bid = ?", bid).Select()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//关注

type FollowShip struct {
	tableName struct{}
	Followee  int
	Follower  int
}

func InsertFollowShip(followee int, follower int) error {
	if err := CheckUserId(followee); err != nil {
		return err
	}
	if err := CheckUserId(follower); err != nil {
		return err
	}

	var fo []FollowShip
	if db.Model(&fo).Where("followee = ?", followee).Where("follower = ?", follower).Select(); fo != nil {
		return errors.New("already followed")
	}

	f := &FollowShip{
		Followee: followee,
		Follower: follower,
	}

	_, err := db.Model(f).Insert()
	if err != nil {
		return err
	}
	return nil
}

func DeleteFollowShip(followee int, follower int) error {
	followShip := &FollowShip{}
	_, err := db.Model(followShip).Where("follower = ?", follower).Where("followee = ?", followee).Delete()
	if err != nil {
		return err
	}
	return nil
}

func GetFollowingOfUser(uid int) ([]User, error) {
	var follow []FollowShip
	if err := db.Model(&follow).Where("follower = ?", uid).Select(); err != nil {
		return nil, err
	}
	var uids []int
	for i, _ := range follow {
		uids = append(uids, follow[i].Followee)
	}
	users, err := GetUsersByUids(uids)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetFollowersOfUser(uid int) ([]User, error) {
	var follow []FollowShip
	if err := db.Model(&follow).Where("followee = ?", uid).Select(); err != nil {
		return nil, err
	}
	var uids []int
	for i, _ := range follow {
		uids = append(uids, follow[i].Follower)
	}
	users, err := GetUsersByUids(uids)
	if err != nil {
		return nil, err
	}
	return users, nil
}