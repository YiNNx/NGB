package controller

import (
	"ngb/model"
	"time"
)

//user

type receiveUserAccount struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,max=20,min=4"`
	Pwd      string `json:"pwd" validate:"required,max=20,min=6"`
}

type responseUserToken struct {
	Uid   int    `json:"uid"`
	Token string `json:"token"`
}

type responseUserProfile struct {
	Username    string         `json:"username"`
	Nickname    string         `json:"nickname"`
	Avatar      string         `json:"avatar"`
	Gender      int            `json:"gender"`
	Posts       []postOutline  `json:"posts"`
	Followers   []userOutline  `json:"followers"`
	Following   []userOutline  `json:"following"`
	Likes       []postOutline  `json:"likes"`
	Collections []postOutline  `json:"collections"`
	BoardsJoin  []boardOutline `json:"boards_join"`
}

type userAccount struct {
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Avatar   string `json:"avatar" validate:"required"`
	Nickname string `json:"nickname" validate:"required"`
	Gender   int    `json:"gender" validate:"required"`
	Intro    string `json:"intro" validate:"required"`
}

type receiveChangePwd struct {
	Email  string `json:"email"  validate:"required,email"`
	PwdOld string `json:"pwd_old"  validate:"required,max=20,min=6"`
	PwdNew string `json:"pwd_new"  validate:"required,max=20,min=6"`
}

type responseAllUser struct {
	Uid        int       `json:"uid"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	CreateTime time.Time `json:"createTime"`
	Role       bool      `json:"role"`
}

//post

type userOutline struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

func NewUerOutline(uid int) (*userOutline, error) {
	u, err := model.GetUserByUid(uid)
	if err != nil {
		return nil, err
	}

	outline := &userOutline{
		Uid:      u.Uid,
		Username: u.Username,
		Avatar:   u.Avatar,
	}

	return outline, nil
}

func NewUerOutlines(u []model.User) []userOutline {
	var outlines []userOutline
	for i, _ := range u {
		outlines = append(outlines, userOutline{
			Uid:      u[i].Uid,
			Username: u[i].Username,
			Avatar:   u[i].Avatar,
		})
	}
	return outlines
}

type boardOutline struct {
	Bid    int    `json:"bid"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Intro  string `json:"intro"`
}

func NewBoardOutline(bid int) (*boardOutline, error) {
	b, err := model.GetBoardByBid(bid)
	if err != nil {
		return nil, err
	}

	outline := &boardOutline{
		Bid:    b.Bid,
		Name:   b.Name,
		Avatar: b.Avatar,
		Intro:  b.Intro,
	}

	return outline, nil
}

func NewBoardOutlines(b []model.Board) []boardOutline {
	var outlines []boardOutline
	for i, _ := range b {
		outlines = append(outlines, boardOutline{
			Bid:    b[i].Bid,
			Avatar: b[i].Avatar,
			Name:   b[i].Name,
			Intro:  b[i].Intro,
		})
	}
	return outlines
}

type postOutline struct {
	Pid        int       `json:"pid"`
	Title      string    `json:"title"`
	Author     int       `json:"author"`
	Time       time.Time `json:"time"`
	Board      int       `json:"board"`
	LikesCount int       `json:"likes_count"`
}

func NewPostOutlines(p []model.Post) ([]postOutline, error) {
	var outlines []postOutline
	for i, _ := range p {
		count, err := model.GetLikesCountOfPost(p[i].Pid)
		if err != nil {
			return nil, err
		}
		outlines = append(outlines, postOutline{
			Pid:        p[i].Pid,
			Title:      p[i].Title,
			Author:     p[i].Author,
			Time:       p[i].Time,
			Board:      p[i].Board,
			LikesCount: count,
		})
	}
	return outlines, nil
}

type commentDetail struct {
	Cid       int         `json:"cid"`
	ParentCid int         `json:"parent_cid"`
	IsAuthor  bool        `json:"is_author"`
	From      userOutline `json:"from"`
	To        int         `json:"to"`
	Time      time.Time   `json:"time"`
	Content   string      `json:"content"`
	//SubComments []commentDetail `json:"sub_comments"`
}

func NewCommentDetails(p []model.Comment) ([]commentDetail, error) {
	var details []commentDetail
	for i, _ := range p {
		user, err := NewUerOutline(p[i].From)
		if err != nil {
			return nil, err
		}
		details = append(details, commentDetail{
			Cid:       p[i].Cid,
			ParentCid: p[i].ParentCid,
			IsAuthor:  p[i].IsAuthor,
			From:      *user,
			To:        p[i].To,
			Time:      p[i].Time,
			Content:   p[i].Content,
		})
	}
	return details, nil
}

type responsePostDetail struct {
	Pid        int          `json:"uid"`
	Title      string       `json:"title"`
	Author     userOutline  `json:"author"`
	Time       time.Time    `json:"time"`
	Board      boardOutline `json:"board"`
	Tags       []string     `json:"tags"`
	Content    string       `json:"content"`
	LikesCount int          `json:"likes_count"`
	//IsLike           bool            `json:"is_like"`
	//CollectionsCount int             `json:"collections_count"`
	//IsCollect        bool            `json:"is_collect"`
	CommentsCount int             `json:"comments_count"`
	Comments      []commentDetail `json:"comments"`
}

type receiveNewPost struct {
	Title   string   `json:"title"  validate:"required"`
	Content string   `json:"content"  validate:"required"`
	Tags    []string `json:"tags"`
}

type responseNewPost struct {
	Pid  int       `json:"pid"`
	Time time.Time `json:"time"`
}

type receiveNewStatus struct {
	Status bool `json:"status" validate:"required" `
}

type responseBoardDetail struct {
	Bid    int           `json:"bid"`
	Name   string        `json:"name"`
	Avatar string        `json:"avatar"`
	Intro  string        `json:"intro"`
	Posts  []postOutline `json:"posts"`
}

type receiveCommentPost struct {
	Content string `json:"content"  validate:"required"`
}

type receiveSubCommentPost struct {
	To      int    `json:"to"  validate:"required"`
	Content string `json:"content"  validate:"required"`
}
