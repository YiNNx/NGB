package model

import (
	"errors"
	"github.com/go-pg/pg/v10"
	"ngb/util/bcrypt"

	"time"
)

type User struct {
	tableName struct{}

	Uid        int       `pg:",pk"`
	Email      string    `pg:",unique,notnull"`
	Username   string    `pg:",unique,notnull"`
	Phone      string    `pg:",unique"`
	PwdHash    string    `pg:",notnull"`
	Role       bool      `pg:"default:false"` //0:default 1:super_admin
	CreateTime time.Time `pg:"default:now()"`

	Avatar   string
	Nickname string
	Gender   int //0:secret 1:female 2:male 3:third gender
	Intro    string

	Likes       []Post `pg:"many2many:likes"`
	Collections []Post `pg:"many2many:collections"`
}

type FollowShip struct {
	tableName struct{}
	Followee  int
	Follower  int
}

// ----------- User -------------

func ValidateUser(email string, pwd string) (*User, error) {
	u := new(User)
	if err := tx.Model(u).Where("email = ?", email).Select(); err != nil {
		return nil, err
	}
	err := bcrypt.ValidatePwd(pwd, u.PwdHash)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UpdateUser(u *User) error {
	_, err := tx.Model(u).
		Column("email", "username", "phone", "avatar", "nickname", "gender", "intro").
		Where("uid = ?", u.Uid).
		Update()
	if err != nil {
		return err
	}
	return nil
}

func ChangePwd(pwdHashNew string, uid int) error {
	u := new(User)
	_, err := tx.Model(u).
		Set("pwd_hash = ?", pwdHashNew).
		Where("uid = ?", uid).
		Update()
	if err != nil {
		return err
	}
	return nil
}

func GetUsersByUidList(uids []int) ([]User, error) {
	if uids == nil {
		return nil, nil
	}
	var users []User
	err := tx.Model(&users).
		Where("uid in (?)", pg.In(uids)).
		Select()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUsersByUsernames(usernames []string) ([]User, error) {
	if usernames != nil {
		var users []User
		err := tx.Model(&users).Column("uid").
			Where("username in (?)", pg.In(usernames)).
			Select()
		if err != nil {
			return nil, err
		}
		return users, nil
	} else {
		return nil, nil
	}
}

// ----------- FollowShip -----------

func InsertFollowShip(followee int, follower int) error {
	if err := CheckPK(&User{Uid: followee}); err != nil {
		return err
	}
	if err := CheckPK(&User{Uid: follower}); err != nil {
		return err
	}
	f := &FollowShip{
		Followee: followee,
		Follower: follower,
	}
	if res, err := tx.Model(f).Where("followee = ?", followee).Where("follower = ?", follower).SelectOrInsert(); err != nil {
		return err
	} else if res == false {
		return errors.New("already exist")
	}
	return nil
}

func DeleteFollowShip(followee int, follower int) error {
	followShip := &FollowShip{}
	_, err := tx.Model(followShip).Where("follower = ?", follower).Where("followee = ?", followee).Delete()
	return err
}

func GetFollowingOfUser(uid int) ([]User, error) {
	var follow []FollowShip
	if err := tx.Model(&follow).Where("follower = ?", uid).Select(); err != nil {
		return nil, err
	}
	var uids []int
	for i, _ := range follow {
		uids = append(uids, follow[i].Followee)
	}
	users, err := GetUsersByUidList(uids)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetFollowersOfUser(uid int) ([]User, error) {
	var follow []FollowShip
	if err := tx.Model(&follow).Where("followee = ?", uid).Select(); err != nil {
		return nil, err
	}
	var uids []int
	for i, _ := range follow {
		uids = append(uids, follow[i].Follower)
	}
	users, err := GetUsersByUidList(uids)
	if err != nil {
		return nil, err
	}
	return users, nil
}
