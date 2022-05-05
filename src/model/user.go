package model

import (
	"github.com/go-pg/pg/v10"
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

	Likes       []Post `pg:"many2many:likes"`
	Collections []Post `pg:"many2many:collections"`
}

func InsertUser(u *User) error {
	_, err := tx.Model(u).Insert()
	if err != nil {
		return err
	}
	return nil
}

// Validate user's email & password
func ValidateUser(email string, pwd string) (*User, error) {
	u := new(User)
	if err := tx.Model(u).Where("email = ?", email).Select(); err != nil {
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

func GetUserByUid(uid int) (*User, error) {
	u := &User{Uid: uid}
	err := tx.Model(u).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return u, nil
}

func GetUsersByUids(uids []int) ([]User, error) {
	if uids != nil {
		var users []User
		err := tx.Model(&users).
			Where("uid in (?)", pg.In(uids)).
			Select()
		if err != nil {
			return nil, err
		}
		return users, nil
	} else {
		return nil, nil
	}
}

// SelectAllUser returns all users' info
func SelectAllUser() ([]User, error) {
	var users []User
	err := tx.Model(&users).Select()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUser(uid int) error {
	u := &User{Uid: uid}
	_, err := tx.Model(u).WherePK().Delete()
	if err != nil {
		return err
	}
	return nil
}

func CheckUserId(uid int) error {
	u := &User{Uid: uid}
	err := tx.Model(u).WherePK().Select()
	if err != nil {
		return err
	}
	return nil
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
