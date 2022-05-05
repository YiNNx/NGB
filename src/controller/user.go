package controller

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"ngb/model"
	"ngb/util"
	"strconv"
)

func SignUP(c echo.Context) error {
	trans := model.BeginTx()
	defer model.CloseTx(trans)

	rec := new(receiveUserAccount)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := validate.Struct(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	pwdHash, err := util.PwdHash(rec.Pwd)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	u := &model.User{
		Email:    rec.Email,
		Username: rec.Username,
		PwdHash:  string(pwdHash),
	}
	if err := model.InsertUser(u); err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	res := &responseUserToken{
		Uid:   u.Uid,
		Token: util.GenerateToken(u.Uid, u.Role),
	}
	return util.SuccessRespond(c, http.StatusOK, res)
}

func LogIn(c echo.Context) error {
	trans := model.BeginTx()
	defer model.CloseTx(trans)

	email := c.QueryParam("email")
	pwd := c.QueryParam("pwd")

	u, err := model.ValidateUser(email, pwd)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	response := &responseUserToken{
		Uid:   u.Uid,
		Token: util.GenerateToken(u.Uid, u.Role),
	}

	return util.SuccessRespond(c, http.StatusOK, response)
}

func GetUserProfile(c echo.Context) error {
	trans := model.BeginTx()
	defer model.CloseTx(trans)

	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	u, err := model.GetUserByUid(uid)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}
	p, err := model.GetPostsByUid(uid)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	posts, err := NewPostOutlines(p)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	fr, err := model.GetFollowersOfUser(uid)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	followers := NewUserOutlines(fr)
	fi, err := model.GetFollowingOfUser(uid)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	following := NewUserOutlines(fi)
	l, err := model.GetLikesOfUser(uid)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	likes, err := NewPostOutlines(l)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	co, err := model.GetCollectionsOfUser(uid)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	collections, err := NewPostOutlines(co)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	b, err := model.GetBoardsOfUser(uid)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	boards := NewBoardOutlines(b)
	println("1")
	res := &responseUserProfile{
		Username:    u.Username,
		Nickname:    u.Nickname,
		Avatar:      u.Avatar,
		Gender:      u.Gender,
		Posts:       posts,
		Followers:   followers,
		Following:   following,
		Likes:       likes,
		Collections: collections,
		BoardsJoin:  boards,
	}

	return util.SuccessRespond(c, http.StatusOK, res)
}

func GetUserAccount(c echo.Context) error {
	trans := model.BeginTx()
	defer model.CloseTx(trans)

	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id

	u, err := model.GetUserByUid(uid)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	res := &userAccount{
		Email:    u.Email,
		Username: u.Username,
		Phone:    u.Phone,
		Avatar:   u.Avatar,
		Nickname: u.Nickname,
		Gender:   u.Gender,
		Intro:    u.Intro,
	}

	return util.SuccessRespond(c, http.StatusOK, res)
}

func ChangeUserInfo(c echo.Context) error {
	trans := model.BeginTx()
	defer model.CloseTx(trans)

	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id

	rec := new(userAccount)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := validate.Struct(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	u := &model.User{
		Uid:      uid,
		Email:    rec.Email,
		Username: rec.Username,
		Phone:    rec.Phone,
		Avatar:   rec.Avatar,
		Nickname: rec.Nickname,
		Gender:   rec.Gender,
		Intro:    rec.Intro,
	}
	err := model.UpdateUser(u)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return util.SuccessRespond(c, http.StatusOK, nil)
}

func ChangeUserPwd(c echo.Context) error {
	trans := model.BeginTx()
	defer model.CloseTx(trans)

	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id
	u, err := model.GetUserByUid(uid)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	change := new(receiveChangePwd)
	if err := c.Bind(change); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, "pwd_old error,please check again")
	}

	if err := validate.Struct(change); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if u.Email != change.Email {
		return util.ErrorResponse(c, http.StatusUnauthorized, "email error,please check again")
	}

	_, err = model.ValidateUser(change.Email, change.PwdOld)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	pwdHashNew, err := util.PwdHash(change.PwdNew)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	err = model.ChangePwd(string(pwdHashNew), uid)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return util.SuccessRespond(c, http.StatusOK, nil)
}

func FollowUser(c echo.Context) error {
	trans := model.BeginTx()
	defer model.CloseTx(trans)

	followee, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	follower := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id

	rec := new(receiveNewStatus)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	fmt.Println(rec.Status)
	if rec.Status {
		if err := model.InsertFollowShip(followee, follower); err != nil {
			model.RollbackTx(trans)
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	} else {
		if err := model.DeleteFollowShip(followee, follower); err != nil {
			model.RollbackTx(trans)
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	return util.SuccessRespond(c, http.StatusOK, nil)
}

func GetAllUsers(c echo.Context) error {
	trans := model.BeginTx()
	defer model.CloseTx(trans)

	users, err := model.SelectAllUser()
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	usersInfo := NewUserInfos(users)
	return util.SuccessRespond(c, http.StatusOK, usersInfo)
}

func DeleteUser(c echo.Context) error {
	trans := model.BeginTx()
	defer model.CloseTx(trans)

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

func GetAdmins(c echo.Context) error {
	trans := model.BeginTx()
	defer model.CloseTx(trans)

	boards, err := model.SelectAllBoards()
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	res := make([]responseAllAdmins, len(boards))
	for i := range boards {
		users, err := model.GetManagersOfBoard(boards[i].Bid)
		if err != nil {
			model.RollbackTx(trans)
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		res[i].Admins = NewUserInfos(users)
		res[i].Bid = boards[i].Bid
		res[i].Name = boards[i].Name
		res[i].Bid = boards[i].Bid
		res[i].Intro = boards[i].Intro
	}
	return util.SuccessRespond(c, http.StatusOK, res)
}

func GetNewNotification(c echo.Context) error {
	trans := model.BeginTx()
	defer model.CloseTx(trans)

	limit, offset, err := paginate(c)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id

	noti, err := model.GetNotificationsByUid(uid, limit, offset)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	var res []responseNotification

	for i := range noti {
		if noti[i].Status == 1 {
			continue
		}
		noti[i].Status = 1
		err := model.UpdateNotificationStatus(&noti[i])
		if err != nil {
			model.RollbackTx(trans)
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}

		if noti[i].Type == 1 {
			m, err := model.GetMessageByMid(noti[i].ContentId)
			if err != nil {
				model.RollbackTx(trans)
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			res = append(res, responseNotification{
				Nid:     noti[i].Nid,
				Type:    noti[i].Type,
				Time:    noti[i].Time,
				Content: m,
			})
		} else if noti[i].Type == 2 {
			m, err := model.GetCommentByCid(noti[i].ContentId)
			if err != nil {
				model.RollbackTx(trans)
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			res = append(res, responseNotification{
				Nid:     noti[i].Nid,
				Type:    noti[i].Type,
				Time:    noti[i].Time,
				Content: m,
			})
		} else {
			m, err := model.GetPostByPid(noti[i].ContentId)
			if err != nil {
				model.RollbackTx(trans)
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			res = append(res, responseNotification{
				Nid:     noti[i].Nid,
				Type:    noti[i].Type,
				Time:    noti[i].Time,
				Content: m,
			})
		}
	}

	return util.SuccessRespond(c, http.StatusOK, res)
}

func GetNotification(c echo.Context) error {
	trans := model.BeginTx()
	defer model.CloseTx(trans)

	limit, offset, err := paginate(c)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id

	noti, err := model.GetNotificationsByUid(uid, limit, offset)
	if err != nil {
		model.RollbackTx(trans)
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	var res []responseNotification

	for i := range noti {
		if noti[i].Status != 1 {
			noti[i].Status = 1
			err := model.UpdateNotificationStatus(&noti[i])
			if err != nil {
				model.RollbackTx(trans)
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
		}
		if noti[i].Type == 1 {
			m, err := model.GetMessageByMid(noti[i].ContentId)
			if err != nil {
				model.RollbackTx(trans)
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			res = append(res, responseNotification{
				Nid:     noti[i].Nid,
				Type:    noti[i].Type,
				Time:    noti[i].Time,
				Content: m,
			})
		} else if noti[i].Type == 2 {
			m, err := model.GetCommentByCid(noti[i].ContentId)
			if err != nil {
				model.RollbackTx(trans)
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			res = append(res, responseNotification{
				Nid:     noti[i].Nid,
				Type:    noti[i].Type,
				Time:    noti[i].Time,
				Content: m,
			})
		} else {
			m, err := model.GetPostByPid(noti[i].ContentId)
			if err != nil {
				model.RollbackTx(trans)
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			res = append(res, responseNotification{
				Nid:     noti[i].Nid,
				Type:    noti[i].Type,
				Time:    noti[i].Time,
				Content: m,
			})
		}
	}

	return util.SuccessRespond(c, http.StatusOK, res)
}
