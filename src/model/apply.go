package model

import "time"

type Apply struct {
	tableName struct{}

	Apid   int       `pg:",pk"`
	Time   time.Time `pg:"default:now()"`
	Type   int
	Uid    int
	Bid    int
	Name   int
	Reason string
}
