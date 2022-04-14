package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ngb/config"
	"ngb/controller"
	"ngb/model"
	"ngb/router"
)

func main() {
	model.Connect()
	defer model.Close()
	uid := make([]int, 5)
	uid = []int{1, 2, 3}
	following, err := model.GetUsersByUids(uid)
	if err != nil {
		fmt.Println(err)
	}
	for i := range following {
		fmt.Println(following[i].Uid)
	}
	f := controller.NewUerOutlines(following)
	fmt.Println(err)
	for i := range f {
		fmt.Println(f[i].Uid)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	router.InitRouters(e)

	e.Logger.Fatal(e.Start(config.C.App.Addr))
}
