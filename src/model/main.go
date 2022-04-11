package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"ngb/config"
)

var db *pg.DB

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
		(*User)(nil), (*Board)(nil), (*Comment)(nil), (*Post)(nil),
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
