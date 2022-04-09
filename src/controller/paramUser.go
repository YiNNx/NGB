package controller

import "time"

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
	Follower    []userOutline  `json:"follower"`
	Following   []userOutline  `json:"following"`
	Likes       []postOutline  `json:"likes"`
	Collections []postOutline  `json:"collections"`
	BoardsJoin  []boardOutline `json:"boards_join"`
}

type userAccount struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Gender   int    `json:"gender"`
	Intro    string `json:"intro"`
}

type receiveChangePwd struct {
	Email  string `json:"email"`
	PwdOld string `json:"pwdOld"`
	PwdNew string `json:"pwdNew"`
}

type responseAllUser struct {
	Uid        int       `json:"uid"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	CreateTime time.Time `json:"createTime"`
	Role       bool      `json:"role"`
}
