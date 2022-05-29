package controller

import (
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"ngb/model"
	"ngb/util"
	"time"
)

var switchType = map[int]string{
	model.TypeMessage:   "message",
	model.TypeComment:   "comment",
	model.TypeMentioned: "mentioned",
	model.TypeNewPost:   "new_post",
}

func publicToMQ(n *Notification) error {
	n.Time = time.Now()
	nBytes, err := json.Marshal(n)
	if err != nil {
		return err
	}

	err = util.Public(nBytes, switchType[n.Type])
	if err != nil {
		return err
	}

	return nil
}

func GetNewNotification(c echo.Context) error {
	tx := model.BeginTx()
	defer tx.Close()

	limit, offset, err := paginate(c)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id

	noti, err := model.GetNotificationsByUid(uid, limit, offset)
	if err != nil {
		tx.Rollback()
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
			tx.Rollback()
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}

		if noti[i].Type == 1 {
			m := &model.Message{Mid: noti[i].ContentId}
			err := model.GetByPK(m)
			if err != nil {
				tx.Rollback()
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			res = append(res, responseNotification{
				Nid:     noti[i].Nid,
				Type:    noti[i].Type,
				Time:    noti[i].Time,
				Content: m,
			})
		} else if noti[i].Type == 2 {
			com := &model.Comment{Cid: noti[i].ContentId}
			err := model.GetByPK(com)
			if err != nil {
				tx.Rollback()
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			res = append(res, responseNotification{
				Nid:     noti[i].Nid,
				Type:    noti[i].Type,
				Time:    noti[i].Time,
				Content: com,
			})
		} else {
			p := &model.Post{Pid: noti[i].ContentId}
			err = model.GetByPK(p)
			if err != nil {
				tx.Rollback()
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			res = append(res, responseNotification{
				Nid:     noti[i].Nid,
				Type:    noti[i].Type,
				Time:    noti[i].Time,
				Content: p,
			})
		}
	}

	return util.SuccessRespond(c, http.StatusOK, res)
}

func GetNotification(c echo.Context) error {
	tx := model.BeginTx()
	defer tx.Close()

	limit, offset, err := paginate(c)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id

	noti, err := model.GetNotificationsByUid(uid, limit, offset)
	if err != nil {
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	var res []responseNotification

	for i := range noti {
		if noti[i].Status != 1 {
			noti[i].Status = 1
			err := model.UpdateNotificationStatus(&noti[i])
			if err != nil {
				tx.Rollback()
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
		}
		if noti[i].Type == 1 {

			m := &model.Message{Mid: noti[i].ContentId}
			err := model.GetByPK(m)
			if err != nil {
				tx.Rollback()
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			res = append(res, responseNotification{
				Nid:     noti[i].Nid,
				Type:    noti[i].Type,
				Time:    noti[i].Time,
				Content: m,
			})
		} else if noti[i].Type == 2 {
			com := &model.Comment{Cid: noti[i].ContentId}
			err := model.GetByPK(com)
			if err != nil {
				tx.Rollback()
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			res = append(res, responseNotification{
				Nid:     noti[i].Nid,
				Type:    noti[i].Type,
				Time:    noti[i].Time,
				Content: com,
			})
		} else {
			p := &model.Post{Pid: noti[i].ContentId}
			err = model.GetByPK(p)
			if err != nil {
				tx.Rollback()
				return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			res = append(res, responseNotification{
				Nid:     noti[i].Nid,
				Type:    noti[i].Type,
				Time:    noti[i].Time,
				Content: p,
			})
		}
	}

	return util.SuccessRespond(c, http.StatusOK, res)
}
