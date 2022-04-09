package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ngb/model"
	"ngb/util"
	"strconv"
)

func SignUP(c echo.Context) error {
	receive := new(receiveSignUp)
	if err := c.Bind(receive); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := validate.Struct(receive); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	pwdHash, err := util.PwdHash(receive.Pwd)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	u := &model.User{
		Email:    receive.Email,
		Username: receive.Username,
		PwdHash:  string(pwdHash),
	}
	if err := model.InsertUser(u); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	response := &responseToken{
		Uid:   u.Uid,
		Token: util.GenerateToken(u.Uid, u.Role),
	}
	return util.SuccessRespond(c, http.StatusOK, response)
}

func LogIn(c echo.Context) error {
	email := c.QueryParam("email")
	pwd := c.QueryParam("pwd")
	u, err := model.ValidateUser(email, pwd)
	if err != nil {
		return util.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}
	response := &responseToken{
		Uid:   u.Uid,
		Token: util.GenerateToken(u.Uid, u.Role),
	}
	return util.SuccessRespond(c, http.StatusOK, response)
}

func GetUser(c echo.Context) error {
	uid, _ := strconv.Atoi(c.Param("uid"))
	u, err := model.GetUserByUid(uid)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	res := responseUserInfo{
		Uid:      u.Uid,
		Email:    u.Email,
		Username: u.Username,
	}
	return util.SuccessRespond(c, http.StatusOK, res)
}

func ChangeInfo(c echo.Context) error {
	uid, _ := strconv.Atoi(c.Param("uid"))
	info := new(receiveChangeInfo)
	if err := c.Bind(info); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := validate.Struct(info); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	pwdHash, err := util.PwdHash(info.Pwd)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	u := &model.User{Uid: uid, PwdHash: string(pwdHash)}
	err = model.UpdateUser(u)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return util.SuccessRespond(c, http.StatusOK, nil)
}

func GetAllUser(c echo.Context) error {
	users, err := model.SelectAllUser()
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	usersInfo := make([]responseAllUser, len(users))
	for i := range users {
		usersInfo[i].Uid = users[i].Uid
		usersInfo[i].Email = users[i].Email
		usersInfo[i].Username = users[i].Username
		usersInfo[i].CreateTime = users[i].CreateTime
		usersInfo[i].Role = users[i].Role
	}
	return util.SuccessRespond(c, http.StatusOK, usersInfo)
}

func DeleteUser(c echo.Context) error {
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := model.CheckUserId(uid); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	err = model.DeleteUser(uid)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return util.SuccessRespond(c, http.StatusOK, nil)
}
