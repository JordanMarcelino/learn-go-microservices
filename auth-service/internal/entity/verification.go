package entity

import "time"

type Verification struct {
	ID       int64
	UserID   int64
	Token    string
	ExpireAt time.Time
}
