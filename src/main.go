package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ngb/config"
	"ngb/model"
	"ngb/router"
	"ngb/util"
)

func main() {
	model.Connect()
	defer model.Close()

	util.InitLogger()

	util.Logger.Info("提示信息")

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	router.InitRouters(e)

	e.Logger.Fatal(e.Start(config.C.App.Addr))

}
