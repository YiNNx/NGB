package model

import (
	"ngb/util"
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

	Followers   []int
	Following   []int
	Posts       []int
	Comments    []int
	Likes       []int
	Collections []int
	BoardsJoin  []int
	BoardsMng   []int
}

func InsertUser(u *User) error {
	_, err := db.Model(u).Insert()
	if err != nil {
		return err
	}
	return nil
}

// Validate user's email & password
func ValidateUser(email string, pwd string) (*User, error) {
	u := new(User)
	if err := db.Model(u).Where("email = ?", email).Select(); err != nil {
		return nil, err
	}
	err := util.ValidatePwd(pwd, u.PwdHash)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Update user's info by id
func UpdateUser(u *User) error {
	_, err := db.Model(u).
		Column("email").
		Column("username").
		Column("phone").
		Column("avatar").
		Column("nickname").
		Column("gender").
		Column("intro").
		Where("uid", u.Uid).
		Update()
	if err != nil {
		return err
	}
	return nil
}

func ChangePwd(pwdHashNew string, uid int) error {
	u := new(User)
	_, err := db.Model(u).
		Set("pwd_hash = ?", pwdHashNew).
		Where("uid = ?", uid).
		Update()
	if err != nil {
		return err
	}
	return nil
}

// GetUser returns user info by id.
func GetUserByUid(uid int) (*User, error) {
	u := &User{Uid: uid}
	err := db.Model(u).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return u, nil
}

// SelectAllUser returns all users' info
func SelectAllUser() ([]User, error) {
	var users []User
	err := db.Model(&users).Select()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUser(uid int) error {
	u := &User{Uid: uid}
	_, err := db.Model(u).WherePK().Delete()
	if err != nil {
		return err
	}
	return nil
}

func CheckUserId(uid int) error {
	u := &User{Uid: uid}
	err := db.Model(u).WherePK().Select()
	if err != nil {
		return err
	}
	return nil
}

func FollowUser(follower int, followee int) error {
	u := &User{Uid: follower}
	if err := db.Model(u).WherePK().Select(); err != nil {
		return err
	}
	u.Following = append(u.Following, followee)
	if _, err := db.Model(u).Column("following").WherePK().Update(); err != nil {
		return err
	}

	u = &User{Uid: followee}
	if err := db.Model(u).WherePK().Select(); err != nil {
		return err
	}
	u.Following = append(u.Following, follower)
	if _, err := db.Model(u).Column("followers").WherePK().Update(); err != nil {
		return err
	}
	return nil
}
