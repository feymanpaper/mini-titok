package model

const (
	prefixFollowList      = "followList:%d"
	prefixFollowCount     = "followCount:%d"
	prefixFollowCountHash = "followCountBucket:%d"
	prefixFanCountHash    = "fanCountBucket:%d"
	FollowListExpire      = 3600
	FollowPairExpire      = 3600

	BucketNum = 3000 //预设3000个桶存储100w数据s
)
