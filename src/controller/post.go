package controller

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"ngb/model"
	"ngb/util"
	"ngb/util/log"
	"regexp"
	"strconv"
)

func NewPost(c echo.Context) error {
	tx := model.BeginTx()
	defer tx.Close()

	rec := new(receiveNewPost)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	author := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id
	bid, err := strconv.Atoi(c.Param("bid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	p := &model.Post{
		Board:   bid,
		Author:  author,
		Tags:    rec.Tags,
		Title:   rec.Title,
		Content: rec.Content,
	}
	if err := model.Insert(p); err != nil {
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err = util.InsertES(p.Pid, p.Title, p.Content)
	if err != nil {
		return err
	}

	users, err := GetUsersMentioned(rec.Content)
	for i := range users {
		n := &Notification{
			Uid:       users[i].Uid,
			Type:      model.TypeMentioned,
			ContentId: p.Pid,
		}
		if err := publicToMQ(n); err != nil {
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	followers, err := model.GetFollowersOfUser(author)
	if err != nil {
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	for i := range followers {
		n := &Notification{
			Uid:       users[i].Uid,
			Type:      model.TypeNewPost,
			ContentId: p.Pid,
		}
		if err := publicToMQ(n); err != nil {
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	res := &responseNewPost{
		Pid:  p.Pid,
		Time: p.Time,
	}
	return util.SuccessRespond(c, http.StatusOK, res)
}

// GetUsersMentioned 匹配在贴文中被提及的用户
func GetUsersMentioned(content string) ([]model.User, error) {
	var usernames []string

	reg := regexp.MustCompile(`@(\S+)`)
	if reg == nil {
		return nil, errors.New("regexp err")
	}
	res := reg.FindAllStringSubmatch(content, -1)
	for i := range res {
		usernames = append(usernames, res[i][1])
	}

	users, err := model.GetUsersByUsernames(usernames)
	if err != nil {
		return nil, err
	}
	return users, err
}

func GetAllPosts(c echo.Context) error {
	tx := model.BeginTx()
	defer tx.Close()

	limit, offset, err := paginate(c)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	posts, err := model.SelectAllPosts(limit, offset)
	if err != nil {
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	res, err := NewPostOutlines(posts)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return util.SuccessRespond(c, http.StatusOK, res)
}

func GetPost(c echo.Context) error {
	tx := model.BeginTx()
	defer tx.Close()

	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	p := &model.Post{Pid: pid}
	err = model.GetByPK(p)
	if err != nil {
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	log.Logger.Error("here")

	_, likesCount, err := model.GetLikesOfPost(pid)
	if err != nil {
		log.Logger.Error("er")
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	_, commentsCount, err := model.GetCommentsOfPost(pid)
	if err != nil {
		tx.Rollback()
		log.Logger.Error("er!")

		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	author, err := NewUserOutline(p.Author)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	board, err := NewBoardOutline(p.Board)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	co, _, err := model.GetCommentsOfPost(pid)
	if err != nil {
		tx.Rollback()
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
	tx := model.BeginTx()
	defer tx.Close()

	limit, offset, err := paginate(c)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	tag := c.Param("tag")
	if tag == "" {
		return util.ErrorResponse(c, http.StatusBadRequest, "")

	}

	posts, err := model.GetPostsByTag(tag, limit, offset)

	if err != nil {
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	res, err := NewPostOutlines(posts)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return util.SuccessRespond(c, http.StatusOK, res)
}

func SearchPost(c echo.Context) error {

	tx := model.BeginTx()
	defer tx.Close()

	keyword := c.QueryParam("keyword")
	pidList, err := util.Search(keyword)
	if err != nil {
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if err != nil {
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	posts, err := model.GetPostsByPidList(pidList)
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
	tx := model.BeginTx()
	defer tx.Close()

	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id
	rec := new(receiveNewStatus)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if rec.Status {
		co := &model.Collection{
			PostPid: pid,
			UserUid: uid,
		}
		if err := model.InsertLikeOrCollection(co, pid, uid); err != nil {
			tx.Rollback()
			log.Logger.Info("http-response:" + err.Error())
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	} else {
		if err := model.DeleteLikeOrCollection(&model.Collection{}, pid, uid); err != nil {
			tx.Rollback()
			log.Logger.Info("http-response:" + err.Error())
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	return util.SuccessRespond(c, http.StatusOK, nil)
}

func LikePost(c echo.Context) error {
	tx := model.BeginTx()
	defer tx.Close()

	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	uid := c.Get("user").(*jwt.Token).Claims.(*util.JwtUserClaims).Id

	rec := new(receiveNewStatus)
	if err := c.Bind(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(rec); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if rec.Status {
		l := &model.Like{
			PostPid: pid,
			UserUid: uid,
		}
		if err := model.InsertLikeOrCollection(l, pid, uid); err != nil {
			tx.Rollback()
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return util.SuccessRespond(c, http.StatusOK, nil)
	} else {
		if err := model.DeleteLikeOrCollection(&model.Like{}, pid, uid); err != nil {
			tx.Rollback()
			return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return util.SuccessRespond(c, http.StatusOK, nil)
	}
}

func CommentPost(c echo.Context) error {
	tx := model.BeginTx()
	defer tx.Close()

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

	p := &model.Post{Pid: pid}
	err = model.GetByPK(p)
	if err != nil {
		tx.Rollback()
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

	err = model.Insert(comment)
	if err != nil {
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	n := &Notification{
		Uid:       p.Author,
		Type:      model.TypeComment,
		ContentId: comment.Cid,
	}
	if err := publicToMQ(n); err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return util.SuccessRespond(c, http.StatusOK, nil)
}

func SubCommentPost(c echo.Context) error {
	tx := model.BeginTx()
	defer tx.Close()

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

	p := &model.Post{Pid: pid}
	err = model.GetByPK(p)
	if err != nil {
		tx.Rollback()
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

	err = model.Insert(comment)
	if err != nil {
		tx.Rollback()
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

func DeletePost(c echo.Context) error {
	tx := model.BeginTx()
	defer tx.Close()

	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	err = model.Delete(&model.Post{Pid: pid})
	if err != nil {
		tx.Rollback()
		return util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return util.SuccessRespond(c, http.StatusOK, nil)
}
