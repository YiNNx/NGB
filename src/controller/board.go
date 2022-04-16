package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ngb/model"
	"ngb/util"
	"strconv"
)

func GetAllBoards(c echo.Context) error {
	boards, err := model.SelectAllBoards()
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	res := NewBoardOutlines(boards)
	return util.SuccessRespond(c, http.StatusOK, res)
}

func GetBoard(c echo.Context) error {
	limit, offset, err := paginate(c)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	bid, err := strconv.Atoi(c.Param("bid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	b, err := model.GetBoardByBid(bid)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	p, err := model.GetPostsByBoard(bid, limit, offset)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	posts, err := NewPostOutlines(p)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	res := &responseBoardDetail{
		Bid:    bid,
		Name:   b.Name,
		Avatar: b.Avatar,
		Intro:  b.Intro,
		Posts:  posts,
	}
	return util.SuccessRespond(c, http.StatusOK, res)
}
