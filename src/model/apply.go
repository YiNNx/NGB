package model

import "time"

//Apply Type:
//	1 - 创建板块申请
//	2 - 管理员申请
const (
	TypeApplyBoard = 1
	TypeApplyAdmin = 2
)

type Apply struct {
	tableName struct{}

	Apid   int       `pg:",pk"`
	Time   time.Time `pg:"default:now()"`
	Type   int
	Uid    int
	Bid    int
	Name   string
	Reason string
	Status int //0-未审核 1-通过 2-不通过
}

func InsertApply(a *Apply) error {
	_, err := tx.Model(a).Insert()
	if err != nil {
		return err
	}
	return nil
}

func SelectBoardApplies() ([]Apply, error) {
	var applies []Apply
	err := tx.Model(&applies).Where("type = 1").Select()
	if err != nil {
		return nil, err
	}
	return applies, nil
}

func SelectAdminApplies() ([]Apply, error) {
	var applies []Apply
	err := tx.Model(&applies).Where("type = 2").Select()
	if err != nil {
		return nil, err
	}
	return applies, nil
}

func SelectApplyByApid(apid int) (*Apply, error) {
	ap := &Apply{Apid: apid}
	err := tx.Model(ap).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return ap, nil
}

func UpdateApplyStatus(a *Apply) error {
	_, err := tx.Model(a).Column("status").WherePK().Update()
	if err != nil {
		return err
	}
	return nil
}
