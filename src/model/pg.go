package model

import (
	"errors"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"ngb/config"
	"ngb/util"
)

var db *pg.DB

func init() {
	// Register many-to-many model so ORM can better recognize m2m relation.
	// This should be done before dependant models are used.
	orm.RegisterTable((*Collection)(nil))
	orm.RegisterTable((*Like)(nil))
}

// Connect database
func Connect() *pg.DB {
	db = pg.Connect(&pg.Options{
		Addr:     config.C.Postgresql.Host + ":" + config.C.Postgresql.Port,
		User:     config.C.Postgresql.User,
		Password: config.C.Postgresql.Password,
		Database: config.C.Postgresql.Dbname,
	})
	var n int
	if _, err := db.QueryOne(pg.Scan(&n), "SELECT 1"); err != nil {
		util.Logger.Panic("Postgresql-connection failed")
	}
	util.Logger.Info("Postgresql connected")
	return db
}

// Close database
func Close() {
	if err := db.Close(); err != nil {
		util.Logger.Panic("Postgresql-close failed")
	}
}

// CreateSchema creates database schema for User model
func CreateSchema() error {
	models := []interface{}{
		//(*Board)(nil), (*User)(nil), (*Post)(nil), (*Comment)(nil),
		//(*FollowShip)(nil), (*Like)(nil), (*Manage)(nil), (*Join)(nil), (*Board)(nil), (*Comment)(nil),
		(*Apply)(nil), (*Notification)(nil), (*Message)(nil),
	}
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Transaction

var tx *pg.Tx

type Transaction struct {
	Tx    *pg.Tx
	abort bool
}

func BeginTx() *Transaction {
	var err error
	tx, err = db.Begin()
	if err != nil {
		util.Logger.Error("tx-begin failed:" + err.Error())
	}
	trans := &Transaction{
		Tx:    tx,
		abort: false,
	}
	return trans
}

func (trans *Transaction) Rollback() {
	err := trans.Tx.Rollback()
	if err != nil {
		util.Logger.Error("tx-close failed:" + err.Error())
	}
	trans.abort = true
	tx = nil
}

func (trans *Transaction) Close() {
	if trans.abort == false {
		err := trans.Tx.Commit()
		if err != nil {
			util.Logger.Error("tx-commit failed:" + err.Error())
		}
	}
	err := trans.Tx.Close()
	if err != nil {
		util.Logger.Error("tx-close failed:" + err.Error())
	}
	tx = nil
}

// Some shared model functions

func Insert(m interface{}) error {
	_, err := tx.Model(m).Insert()
	return err
}

// CheckPK , when doesn't exist return false
func CheckPK(m interface{}) (err error) {
	res, err := tx.Model(m).WherePK().Exists()
	if err != nil {
		return err
	}
	if res == false {
		return errors.New("PK doesn't exist")
	}
	return nil
}

func GetByPK(m interface{}) error {
	err := tx.Model(m).WherePK().Select()
	return err
}

func GetAll(m interface{}) error {
	//m *[]model.Xx
	err := tx.Model(m).Select()
	return err
}

func Delete(m interface{}) error {
	_, err := tx.Model(m).WherePK().Delete()
	if err != nil {
		return err
	}
	return nil
}
