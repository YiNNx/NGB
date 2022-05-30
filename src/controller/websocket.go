package controller

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	myware "ngb/middleware"
	"ngb/model"
	"ngb/util"
	"ngb/util/log"
	"strconv"
)

func Chat(c echo.Context) error {
	with, err := strconv.Atoi(c.QueryParam("with"))
	if err != nil {
		log.Logger.Error(err)
		return err
	}
	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id
	ws := c.(*myware.WsContext).Conn

	readClient := util.GetClient(uid, with, ws)
	writeClient := util.GetClient(with, uid, ws)

	wait := make(chan bool)
	go readClient.ReadMsg()
	go writeClient.WriteMsg()
	<-wait

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

	n := &util.Notification{
		Uid:       receiver,
		Type:      model.TypeMessage,
		ContentId: m.Mid,
	}
	if err := util.PublishToMQ(n); err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return util.SuccessRespond(c, http.StatusOK, m)
}
