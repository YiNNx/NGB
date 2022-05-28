package model

import (
	"errors"
	"time"
)

type Board struct {
	tableName struct{}

	Bid    int    `pg:",pk"`
	Name   string `pg:",notnull"`
	Avatar string
	Time   time.Time `pg:"default:now()"`
	Intro  string
}

type ManageShip struct {
	Bid int
	Uid int
}

// ----------- Board -------------

func UpdateBoard(b *Board) error {
	_, err := tx.Model(b).
		Column("name", "avatar", "intro").
		Where("bid = ?", b.Bid).
		Update()
	if err != nil {
		return err
	}
	return nil
}

// --------- ManageShip ---------

func InsertManageShip(bid int, uid int) error {
	if err := CheckPK(&Board{Bid: bid}); err != nil {
		return err
	}
	f := &ManageShip{
		Bid: bid,
		Uid: uid,
	}
	if res, err := tx.Model(f).Where("uid = ?", uid).Where("bid = ?", bid).SelectOrInsert(); err != nil {
		return err
	} else if res == false {
		return errors.New("already exist")
	}
	return nil
}

func GetBoardsMngOfUser(uid int) ([]Board, error) {
	var user User
	err := tx.Model(&user).Relation("Manages").Where("uid = ?", uid).Select()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func GetManagersOfBoard(bid int) ([]User, error) {
	var ms []ManageShip
	if err := tx.Model(&ms).Where("bid = ?", bid).Select(); err != nil {
		return nil, err
	}
	var uids []int
	for i, _ := range ms {
		uids = append(uids, ms[i].Uid)
	}
	users, err := GetUsersByUidList(uids)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func CheckAdmin(bid int, uid int) (bool, error) {
	m := new(ManageShip)
	if err := tx.Model(m).Where("bid = ?", bid).Where("uid = ?", uid).Select(); err != nil {
		return false, err
	}
	if m == nil {
		return false, nil
	}
	return true, nil
}
