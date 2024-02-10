package model

import (
	"fmt"
	"hash/crc32"
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

func formatFanListKey(uid int64) string {
	return fmt.Sprintf(prefixFanList, uid)
}

func formatFollowCountKey(uid int64) string {
	return fmt.Sprintf(prefixFollowCount, uid)
}

func formatFollowCountKeyHash(uid uint32) string {
	return fmt.Sprintf(prefixFollowCountHash, uid)
}

func formatFanCountKeyHash(uid uint32) string {
	return fmt.Sprintf(prefixFanCountHash, uid)
}

func hashBucketNum(key string) uint32 {
	num := crc32.ChecksumIEEE([]byte(key)) % BucketNum
	return num
}
