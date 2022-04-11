package main

import (
	"ngb/model"
)

func main() {
	model.Connect()
	defer model.Close()
	model.CreateSchema()
	//following, err := model.GetMembersOfBoard(1)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//for i := range following {
	//	fmt.Println(following[i].Uid)
	//}
	//
	////e := echo.New()
	////e.Use(middleware.Logger())
	////e.Use(middleware.Recover())
	////
	////router.InitRouters(e)
	////
	////e.Logger.Fatal(e.Start(config.C.App.Addr))
}
