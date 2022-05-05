package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"ngb/config"
	"ngb/util"
)

var db *pg.DB

func init() {
	// Register many to many model so ORM can better recognize m2m relation.
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
		panic(err)
	}
	util.Logger.Info("Postgresql connected")
	return db
}

// Close database
func Close() {
	db.Close()
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

var tx *pg.Tx

type Transaction struct {
	Tx    *pg.Tx
	abort bool
}

func BeginTx() *Transaction {
	var err error
	tx, err = db.Begin()
	if err != nil {
		util.Logger.Panic("tx-begin failed")
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
		util.Logger.Panic("tx-close failed")
	}
	trans.abort = true
	tx = nil
}

func (trans *Transaction) Close() {
	if trans.abort == false {
		err := trans.Tx.Commit()
		if err != nil {
			util.Logger.Panic("tx-commit failed")
		}
	}
	err := trans.Tx.Close()
	if err != nil {
		util.Logger.Panic("tx-close failed")
	}
	tx = nil
}
