package controller

import "time"

type userOutline struct {
	Uid      int    `json:"uid"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
}

type boardOutline struct {
	Bid    int    `json:"bid"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Intro  string `json:"intro"`
}

type postOutline struct {
	Pid           int          `json:"pid"`
	Title         string       `json:"title"`
	Author        string       `json:"author"`
	Time          time.Time    `json:"time"`
	Board         boardOutline `json:"board"`
	LikesCount    int          `json:"likes_count"`
	CommentsCount int          `json:"comments_count"`
}

type commentDetail struct {
	Cid         int             `json:"cid"`
	From        userOutline     `json:"from"`
	To          userOutline     `json:"to"`
	Time        time.Time       `json:"time"`
	IsAuthor    bool            `json:"is_author"`
	Content     string          `json:"content"`
	SubComments []commentDetail `json:"sub_comments"`
}

type responsePostDetail struct {
	Pid              int             `json:"uid"`
	Title            string          `json:"title"`
	Author           userOutline     `json:"author"`
	Time             time.Time       `json:"time"`
	Board            boardOutline    `json:"board"`
	Tags             []string        `json:"tags"`
	Content          string          `json:"content"`
	LikesCount       int             `json:"likes_count"`
	IsLike           bool            `json:"is_like"`
	Likes            []userOutline   `json:"likes"`
	CollectionsCount int             `json:"collections_count"`
	IsCollect        bool            `json:"is_collect"`
	CommentsCount    int             `json:"comments_count"`
	Comments         []commentDetail `json:"comments"`
}

type responsePosts struct {
	posts []responsePostDetail `json:"posts"`
}

type receiveNewPost struct {
	Bid     int    `json:"bid"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type responseNewPost struct {
	Pid  int       `json:"pid"`
	Time time.Time `json:"time"`
}

type receiveNewCollection struct {
	Collect bool `json:"collect"`
}

type receiveNewStatus struct {
	Status bool `json:"status"`
}

type responseBoardDetail struct {
	Bid    int           `json:"bid"`
	Name   string        `json:"name"`
	Avatar string        `json:"avatar"`
	Intro  string        `json:"intro"`
	Posts  []postOutline `json:"posts"`
}

type responseAllBoards struct {
	Boards []boardOutline
}
