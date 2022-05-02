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

func SetBoardApply(c echo.Context) error {
	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id

	rec := new(receiveBoardApply)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	apply := &model.Apply{
		Type:   1,
		Uid:    uid,
		Name:   rec.Name,
		Reason: rec.Reason,
	}
	if err := model.InsertApply(apply); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return util.SuccessRespond(c, http.StatusOK, nil)
}

func SetAdminApply(c echo.Context) error {
	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id

	rec := new(receiveAdminApply)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	apply := &model.Apply{
		Type:   2,
		Uid:    uid,
		Bid:    rec.Bid,
		Reason: rec.Reason,
	}
	if err := model.InsertApply(apply); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return util.SuccessRespond(c, http.StatusOK, nil)
}

func GetBoardApply(c echo.Context) error {
	ap, err := model.SelectBoardApplies()
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	var applies []responseBoardApply
	for i, _ := range ap {
		user, err := NewUserOutline(ap[i].Uid)
		if err != nil {
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		applies = append(applies, responseBoardApply{
			Type:      ap[i].Type,
			Apid:      ap[i].Apid,
			Time:      ap[i].Time,
			Applicant: *user,
			Name:      ap[i].Name,
			Reason:    ap[i].Reason,
			Status:    ap[i].Status,
		})
	}
	return util.SuccessRespond(c, http.StatusOK, applies)
}

func GetAdminApply(c echo.Context) error {
	ap, err := model.SelectAdminApplies()
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	var applies []responseAdminApply
	for i, _ := range ap {
		print("2")
		user, err := NewUserOutline(ap[i].Uid)
		if err != nil {
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		board, err := NewBoardOutline(ap[i].Bid)
		fmt.Println(ap[i])
		if err != nil {
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		print("1")
		applies = append(applies, responseAdminApply{
			Type:      ap[i].Type,
			Apid:      ap[i].Apid,
			Time:      ap[i].Time,
			Applicant: *user,
			Board:     *board,
			Reason:    ap[i].Reason,
			Status:    ap[i].Status,
		})
	}
	return util.SuccessRespond(c, http.StatusOK, applies)
}

func PassApply(c echo.Context) error {
	rec := new(receiveNewStatus)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	apid, err := strconv.Atoi(c.Param("apid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	apply, err := model.SelectApplyByApid(apid)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	if apply.Status != 0 {
		return util.ErrorResponse(c, http.StatusInternalServerError, "apply has already checked")
	}

	if !rec.Status {
		apply.Status = 2
		err := model.UpdateApplyStatus(apply)
		if err != nil {
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return util.SuccessRespond(c, http.StatusOK, nil)
	} else {
		apply.Status = 1
		err := model.UpdateApplyStatus(apply)
		if err != nil {
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	if apply.Type == 1 {
		b := &model.Board{
			Name: apply.Name,
		}
		err := model.InsertBoard(b)
		if err != nil {
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return util.SuccessRespond(c, http.StatusOK, nil)
	} else {

		err := model.InsertManageShip(apply.Bid, apply.Uid)
		if err != nil {
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return util.SuccessRespond(c, http.StatusOK, nil)
	}
}
