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
	receiver, err := strconv.Atoi(c.QueryParam("receiver"))
	sender := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
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
	if err := model.InsertMessage(m); err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	if err := model.InsertNotification(model.TypeMessage, receiver, m.Mid); err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return util.SuccessRespond(c, http.StatusOK, m)
}
