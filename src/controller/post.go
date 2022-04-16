package controller

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"ngb/model"
	"ngb/util"
	"strconv"
)

func NewPost(c echo.Context) error {
	rec := new(receiveNewPost)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := validate.Struct(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	bid, err := strconv.Atoi(c.Param("bid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	p := &model.Post{
		Board:   bid,
		Author:  c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id,
		Tags:    rec.Tags,
		Title:   rec.Title,
		Content: rec.Content,
	}
	if err := model.InsertPost(p); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	res := &responseNewPost{
		Pid:  p.Pid,
		Time: p.Time,
	}
	return util.SuccessRespond(c, http.StatusOK, res)
}

func GetAllPosts(c echo.Context) error {
	limit, offset, err := paginate(c)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	posts, err := model.SelectAllPosts(limit, offset)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	res, err := NewPostOutlines(posts)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return util.SuccessRespond(c, http.StatusOK, res)
}

func GetPost(c echo.Context) error {
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	p, err := model.GetPostByPid(pid)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	likesCount, err := model.GetLikesCountOfPost(pid)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	commentsCount, err := model.GetCommentsCountOfPost(pid)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	author, err := NewUerOutline(p.Author)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	board, err := NewBoardOutline(p.Board)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	co, err := model.GetCommentsByPid(pid)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	comments, err := NewCommentDetails(co)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	res := &responsePostDetail{
		Pid:           pid,
		Title:         p.Title,
		Author:        *author,
		Time:          p.Time,
		Board:         *board,
		Tags:          p.Tags,
		Content:       p.Content,
		LikesCount:    likesCount,
		CommentsCount: commentsCount,
		Comments:      comments,
	}
	return util.SuccessRespond(c, http.StatusOK, res)
}

func GetPostsByTag(c echo.Context) error {
	limit, offset, err := paginate(c)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	tag := c.Param("tag")
	posts, err := model.GetPostsByTag(tag, limit, offset)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	res, err := NewPostOutlines(posts)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return util.SuccessRespond(c, http.StatusOK, res)
}

func CollectPost(c echo.Context) error {
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id
	rec := new(receiveNewStatus)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if rec.Status {
		if err := model.InsertCollection(pid, uid); err != nil {
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	} else {
		if err := model.DeleteCollection(pid, uid); err != nil {
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	return util.SuccessRespond(c, http.StatusOK, nil)
}

func LikePost(c echo.Context) error {
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id

	rec := new(receiveNewStatus)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if rec.Status {
		if err := model.InsertLike(pid, uid); err != nil {
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return util.SuccessRespond(c, http.StatusOK, nil)
	} else {
		if err := model.DeleteLike(pid, uid); err != nil {
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return util.SuccessRespond(c, http.StatusOK, nil)
	}
}

func CommentPost(c echo.Context) error {
	rec := new(receiveCommentPost)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := validate.Struct(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	p, err := model.GetPostByPid(pid)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id
	isAuthor := false
	if p.Author == uid {
		isAuthor = true
	}

	comment := &model.Comment{
		Post:     pid,
		IsAuthor: isAuthor,
		From:     uid,
		Content:  rec.Content,
	}

	err = model.InsertComment(comment)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return util.SuccessRespond(c, http.StatusOK, nil)
}

func SubCommentPost(c echo.Context) error {
	rec := new(receiveSubCommentPost)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := validate.Struct(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	cid, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	p, err := model.GetPostByPid(pid)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id
	isAuthor := false
	if p.Author == uid {
		isAuthor = true
	}

	comment := &model.Comment{
		Post:      pid,
		IsAuthor:  isAuthor,
		ParentCid: cid,
		From:      uid,
		To:        rec.To,
		Content:   rec.Content,
	}

	err = model.InsertComment(comment)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return util.SuccessRespond(c, http.StatusOK, nil)
}

func paginate(c echo.Context) (int, int, error) {
	a, p := c.QueryParam("amount"), c.QueryParam("page")
	if a == "" {
		a = "10"
	}
	if p == "" {
		p = "1"
	}
	limit, err := strconv.Atoi(a)
	if err != nil {
		return 0, 0, err
	}
	page, err := strconv.Atoi(p)
	if err != nil {
		return 0, 0, err
	}
	offset := limit * (page - 1)
	return limit, offset, nil
}
