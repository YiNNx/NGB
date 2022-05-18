package controller

import (
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
	"ngb/model"
	"ngb/util"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Chat(c echo.Context) error {

	with, err := strconv.Atoi(c.QueryParam("with"))
	if err != nil {
		util.Logger.Error(err)
		return err
	}
	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		util.Logger.Error(err)
		return err
	}

	readClient := util.GetClient(uid, with, ws)
	writeClient := util.GetClient(with, uid, ws)

	go readClient.ReadMsg()
	go writeClient.WriteMsg()

	wg.Add(1)
	wg.Wait()
	return nil
}

//--------- Old Http Api ---------

func SendMessage(c echo.Context) error {
	tx := model.BeginTx()
	defer tx.Close()

	receiver, err := strconv.Atoi(c.QueryParam("receiver"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	sender := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id
	rec := new(receiveMessage)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	m := &model.Message{
		Sender:   sender,
		Receiver: receiver,
		Content:  rec.Content,
	}
	if err := model.Insert(m); err != nil {
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	n := &model.Notification{
		Uid:       receiver,
		Type:      model.TypeMessage,
		ContentId: m.Mid,
	}
	if err := model.Insert(n); err != nil {
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return util.SuccessRespond(c, http.StatusOK, m)
}
