package model

import "time"

//Notification Type:
//	1 - 私信
//	2 - 评论
//	3 - @
//	4 - 关注人发帖

const (
	TypeMessage   = 1
	TypeComment   = 2
	TypeMentioned = 3
	TypeNewPost   = 4
)

type Notification struct {
	tableName struct{}

	Nid       int       `pg:",pk"`
	Time      time.Time `pg:"default:now()"`
	Uid       int
	Type      int
	ContentId int //私信为mid,关注人发帖和@为pid,评论为cid
	Status    int `pg:"default:0"` //0未读 1已读
}

type Message struct {
	Mid      int       `pg:",pk"`
	Time     time.Time `pg:"default:now()"`
	Sender   int
	Receiver int
	Content  string
}

func InsertNotification(notiType int, uid int, contentId int) error {
	tx, err := tx.Begin()
	if err != nil {
		return err
	}
	defer tx.Close()
	n := &Notification{
		Uid:       uid,
		Type:      notiType,
		ContentId: contentId,
	}
	_, err = tx.Model(n).Insert()
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func InsertMessage(n *Message) error {
	tx, err := tx.Begin()
	if err != nil {
		return err
	}
	defer tx.Close()

	_, err = tx.Model(n).Insert()
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func GetMessageByMid(mid int) (*Message, error) {
	m := &Message{Mid: mid}
	err := tx.Model(m).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func GetNotificationsByUid(uid int, limit int, offset int) ([]Notification, error) {
	var notifications []Notification
	err := tx.Model(&notifications).
		Where("uid = ?", uid).
		Limit(limit).Offset(offset).Select()
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func UpdateNotificationStatus(a *Notification) error {
	_, err := tx.Model(a).Column("status").WherePK().Update()
	if err != nil {
		return err
	}
	return nil
}
