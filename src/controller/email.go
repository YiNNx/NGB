package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ngb/util"
)

func SendEmail(c echo.Context) error {
	rec := new(receiveSendEmail)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	res, err := util.EmailPool(rec.To, rec.Subject, rec.Content)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return util.SuccessRespond(c, http.StatusOK, res)
}
