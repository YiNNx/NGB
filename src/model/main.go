package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"ngb/config"
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
		User:     config.C.Postgresql.User,
		Password: config.C.Postgresql.Password,
		Database: config.C.Postgresql.Dbname,
	})

	var n int
	if _, err := db.QueryOne(pg.Scan(&n), "SELECT 1"); err != nil {
		panic(err)
	}

	return db
}

// Close database
func Close() {
	db.Close()
}

// CreateSchema creates database schema for User model
func CreateSchema() error {
	models := []interface{}{
		(*Board)(nil), (*User)(nil), (*Post)(nil), (*Comment)(nil),
		//(*FollowShip)(nil), (*Like)(nil), (*Manage)(nil), (*Join)(nil), (*Board)(nil), (*Comment)(nil),
		//(*Join)(nil), (*Manage)(nil),
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
