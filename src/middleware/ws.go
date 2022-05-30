package middleware

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"ngb/util/log"
)

type WsContext struct {
	echo.Context
	*websocket.Conn
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WsUpgrade(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ws, err := upgrade.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			log.Logger.Error(err)
			return err
		}
		cc := &WsContext{c, ws}
		return next(cc)
	}
}
