package model

import (
	"fmt"
	"strconv"
)

func formatFollowKey(fromId int64, toId int64) string {
	fromIdStr := strconv.FormatInt(fromId, 10)
	toIdStr := strconv.FormatInt(toId, 10)
	key := fromIdStr + ":" + toIdStr
	return key
}

func formatFollowListKey(uid int64) string {
	return fmt.Sprintf(prefixFollowList, uid)
}

func formatFollowCountKey(uid int64) string {
	return fmt.Sprintf(prefixFollowCount, uid)
}

func formatFollowCountKeyHash(uid uint32) string {
	return fmt.Sprintf(prefixFollowCountHash, uid)
}
