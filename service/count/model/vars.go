package model

import (
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound

// 不存在cache缓存空值的过期时间 5分钟
const notFoundExpiry = 300

// cache缓存过期时间 1小时
const CacheExpiry = 3600
const notFoundCachePlaceholder = "*"

const NotFoundCount = -1

// 缓存了空值
var errCachePlaceNotExist = errors.New("errCacheNotExist")

var errCacheNotFound = errors.New("errorCacheNotFound")

var errDBNotFound = errors.New("errorDBNotFound")
