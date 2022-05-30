package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ngb/config"
	"ngb/model"
	"ngb/router"
	"ngb/util/log"
)

func main() {
	model.Connect()
	defer model.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	router.InitRouters(e)

	log.Logger.Fatal(e.Start(config.C.App.Addr))
}
