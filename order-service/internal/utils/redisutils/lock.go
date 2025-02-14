package redisutils

import (
	"fmt"
)

func NewLockKey(requestID string, userID int64) string {
	return fmt.Sprintf("lock:%v-%v", requestID, userID)
}
