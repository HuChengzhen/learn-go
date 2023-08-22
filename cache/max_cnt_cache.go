package cache

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var (
	errOverCapacity = errors.New("cache: 超过容量限制")
)

type MaxCntCache struct {
	*BuildInMapCache
	cnt    int32
	maxCnt int32
}

func NewMaxCntCache(c *BuildInMapCache, maxCnt int32) *MaxCntCache {
	res := &MaxCntCache{
		BuildInMapCache: c,
		maxCnt:          maxCnt,
	}

	origin := c.onEvicted

	res.onEvicted = func(key string, val any) {
		atomic.AddInt32(&res.cnt, -1)
		if origin != nil {
			origin(key, val)
		}
	}

	return res
}

func (c *MaxCntCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	//// 这种写法，如果key已经存在这个计数就不准了
	//cnt := atomic.AddInt32(&c.cnt, 1)
	//if cnt > c.maxCnt {
	//	atomic.AddInt32(&c.cnt, -1)
	//	return errOverCapacity
	//}
	//return c.BuildInMapCache.Set(ctx, key, val, expiration)

	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, ok := c.data[key]
	if !ok {
		if c.cnt+1 > c.maxCnt {
			// 可以在这里设置复杂的淘汰策略
			return errOverCapacity
		}
		c.cnt++
	}
	return c.set(key, val, expiration)

}
