package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	myware "ngb/middleware"
	"ngb/model"
	"ngb/util"
)

func GetNewNotification(c echo.Context) error {
	tx := model.BeginTx()
	defer tx.Close()

	limit, offset, err := paginate(c)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	uid := c.(*myware.SessionContext).Uid

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

	return util.SuccessResponse(c, http.StatusOK, res)
}

func GetNotification(c echo.Context) error {
	tx := model.BeginTx()
	defer tx.Close()

	limit, offset, err := paginate(c)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	uid := c.(*myware.SessionContext).Uid

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

	return util.SuccessResponse(c, http.StatusOK, res)
}
