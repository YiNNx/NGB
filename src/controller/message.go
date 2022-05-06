package controller

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"ngb/model"
	"ngb/util"
	"strconv"
)

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
		Uid:       model.TypeMessage,
		Type:      receiver,
		ContentId: m.Mid,
	}
	if err := model.Insert(n); err != nil {
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return util.SuccessRespond(c, http.StatusOK, m)
}
