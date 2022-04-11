package main

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"ngb/model"
	"time"
)

func init() {
	// Register many to many model so ORM can better recognize m2m relation.
	// This should be done before dependant models are used.
	orm.RegisterTable((*Collection)(nil))
}

type User struct {
	tableName struct{}

	Uid        int       `pg:",pk"`
	Email      string    //`pg:",unique,notnull"`
	Username   string    //`pg:",unique,notnull"`
	Phone      string    `pg:",unique"`
	PwdHash    string    //`pg:",notnull"`
	Role       bool      `pg:"default:false"` //0:default 1:super_admin
	CreateTime time.Time `pg:"default:now()"`

	Avatar   string
	Nickname string
	Gender   int //0:secret 1:female 2:male 3:third gender
	Intro    string

	//Followers   []User    `pg:"many2many:follow_ships"`
	//Following   []User    `pg:"many2many:follow_ships"`
	//Posts       []Post    `pg:"rel:has-many"`
	//Comments    []Comment `pg:"rel:has-many"`
	//Likes []Post `pg:"many2many:likes"`
	Collections []*Post `pg:"many2many:collections"`
	//BoardsJoin  []Board   `pg:"many2many:join_ships"`
	//BoardsMng   []Board   `pg:"many2many:manage_ships"`
}

type Post struct {
	tableName struct{}

	Pid int `pg:",pk"`
	//Board  int       `pg:",notnull"`
	Time time.Time `pg:"default:now()"`
	//Author int       `pg:",notnull"`
	Tags []string

	//Title   string `pg:",notnull"`
	//Content string `pg:",notnull"`

	//Comments    []Comment `pg:"rel:has-many"`
	//Likes []User `pg:"many2many:likes"`
	Collections []User `pg:"many2many:collections"`
}

type Collection struct {
	UserUid int
	PostPid int
}

func ExampleDB_Model_manyToMany() {
	db := model.Connect()
	defer db.Close()

	if err := createManyToManyTables(db); err != nil {
		panic(err)
	}

	values := []interface{}{
		&Post{Pid: 1},
		&Post{Pid: 2},
		&User{Uid: 1},
		&Collection{UserUid: 1, PostPid: 1},
		&Collection{UserUid: 1, PostPid: 2},
	}
	for _, v := range values {
		_, err := db.Model(v).Insert()
		if err != nil {
			panic(err)
		}
	}

	user := new(User)
	err := db.Model(user).Relation("Collections").First()
	if err != nil {
		panic(err)
	}
	fmt.Println(user.Collections[1].Pid)

	p := new(Post)
	err = db.Model(p).Relation("Collections").First()
	if err != nil {
		panic(err)
	}
	fmt.Println(p.Collections[0].Uid)

}

func createManyToManyTables(db *pg.DB) error {
	models := []interface{}{
		(*User)(nil),
		(*Post)(nil),
		(*Collection)(nil),
	}
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	ExampleDB_Model_manyToMany()
	//u := 2
	//posts, err := model.GetLikesByUid(u)
	//fmt.Println(posts)
	//fmt.Println(err)
	////
	////e := echo.New()
	////e.Use(middleware.Logger())
	////e.Use(middleware.Recover())
	////
	////router.InitRouters(e)
	////
	////e.Logger.Fatal(e.Start(config.C.App.Addr))
}
