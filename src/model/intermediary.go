package model

type Like struct {
	UserUid int
	PostPid int
}

func GetLikesOfUser(uid int) ([]Post, error) {
	var user User
	err := db.Model(&user).Relation("Likes").Where("uid = ?", uid).Select()
	if err != nil {
		panic(err)
	}
	return user.Likes, nil
}

func GetLikesOfPost(pid int) ([]User, error) {
	var post Post
	err := db.Model(&post).Relation("Likes").Where("pid = ?", pid).Select()
	if err != nil {
		return nil, err
	}
	return post.Likes, nil
}

type Collection struct {
	UserUid int
	PostPid int
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

type Join struct {
	UserUid  int
	BoardBid int
}

func GetBoardsOfUser(uid int) ([]Board, error) {
	var user User
	err := db.Model(&user).Relation("Join").Where("uid = ?", uid).Select()
	if err != nil {
		return nil, err
	}
	return user.BoardsJoin, nil
}

func GetMembersOfBoard(bid int) ([]User, error) {
	var b Board
	err := db.Model(&b).Relation("Join").Where("bid = ?", bid).Select()
	if err != nil {
		return nil, err
	}
	return b.Members, nil
}

type Manage struct {
	BoardBid int
	UserUid  int
}

func GetBoardsMngOfUser(uid int) ([]Board, error) {
	var user User
	err := db.Model(&user).Relation("Manages").Where("uid = ?", uid).Select()
	if err != nil {
		return nil, err
	}
	return user.BoardsMng, nil
}

func GetManagersOfBoard(bid int) ([]User, error) {
	var b Board
	err := db.Model(&b).Relation("Manages").Where("bid = ?", bid).Select()
	if err != nil {
		return nil, err
	}
	return b.Managers, nil
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
//
//func InsertLike(pid int, uid int) error {
//	l := &Like{
//		Pid: pid,
//		Uid: uid,
//	}
//	_, err := db.Model(l).Insert()
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func InsertManageShip(bid int, uid int) error {
//	f := ManageShip{
//		Bid: bid,
//		Uid: uid,
//	}
//	_, err := db.Model(f).Insert()
//	if err != nil {
//		return err
//	}
//	return nil
//}

type FollowShip struct {
	tableName struct{}
	Followee  int
	Follower  int
}

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

//
//func GetLikesByPid(pid int) ([]Post, error) {
//	var likes []Like
//	if err := db.Model(&likes).Where("pid = ?", pid).Select(); err != nil {
//		return nil, err
//	}
//	var uids []int
//	for i, _ := range likes {
//		uids = append(uids, likes[i].Uid)
//	}
//	posts, err := GetPostsByUids(uids)
//	if err != nil {
//		return nil, err
//	}
//	return posts, nil
//}
