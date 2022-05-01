package model

import "time"

//Notification Type:
//	1 - 私信
//	2 - 评论
//	3 - @
//	4 - 关注人发帖

type Notification struct {
	tableName struct{}

	Nid       int       `pg:",pk"`
	Time      time.Time `pg:"default:now()"`
	Uid       int
	Type      int
	ContentId int //私信为mid,关注人发帖和@为pid,评论为cid
}

type Message struct {
	Mid      int       `pg:",pk"`
	Time     time.Time `pg:"default:now()"`
	Sender   int
	Receiver int
	Content  string
}
