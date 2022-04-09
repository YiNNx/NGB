package util

import "github.com/labstack/echo/v4"

type Response struct {
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func ErrorResponse(c echo.Context, code int, msg string) error {
	return c.JSON(
		code,
		Response{
			Success: false,
			Msg:     msg,
			Data:    nil,
		})
}

func SuccessRespond(c echo.Context, code int, data interface{}) error {
	return c.JSON(
		code,
		Response{
			Success: true,
			Msg:     "",
			Data:    data,
		})
}
