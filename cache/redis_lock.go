package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"time"
)

var (
	ErrFailedToPreemptLock = errors.New("redis-lock: 抢锁失败")
	ErrLockNotHold         = errors.New("redis-lock: 你没有持有锁")

	//go:embed lua/unlock.lua
	luaUnlock string

	//go:embed lua/refresh.lua
	luaRefresh string

	//go:embed lua/lock.lua
	luaLock string
)

type Client struct {
	client redis.Cmdable
	g      singleflight.Group
}

func NewClient(client redis.Cmdable) *Client {
	return &Client{
		client: client,
		g:      singleflight.Group{},
	}
}
func (c *Client) SingleflightLock(ctx context.Context, key string, expiration time.Duration, timeout time.Duration, retry RetryStrategy) (*Lock, error) {
	for {
		var flag bool = false
		doChan := c.g.DoChan(key, func() (interface{}, error) {
			flag = true
			return c.Lock(ctx, key, expiration, timeout, retry)
		})

		select {
		case res := <-doChan:
			if flag {
				c.g.Forget(key)
				if res.Err != nil {
					return nil, res.Err
				}
				return res.Val.(*Lock), nil
			}
		case <-ctx.Done():
		}
	}
}

func (c *Client) Lock(ctx context.Context, key string, expiration time.Duration, timeout time.Duration, retry RetryStrategy) (*Lock, error) {
	var timer *time.Timer
	val := uuid.New().String()
	for {

		lctx, cancel := context.WithTimeout(ctx, timeout)
		result, err := c.client.Eval(lctx, luaLock, []string{key}, val, expiration.Seconds()).Result()
		cancel()
		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}

		if result == "OK" {
			return &Lock{
				client:     c.client,
				key:        key,
				value:      val,
				expiration: expiration,
			}, nil
		}

		interval, ok := retry.Next()
		if !ok {
			return nil, fmt.Errorf("超出重试限制, %w", ErrFailedToPreemptLock)
		}

		if timer == nil {
			timer = time.NewTimer(interval)
		} else {
			timer.Reset(interval)
		}

		select {
		case <-timer.C:
		case <-ctx.Done():
			return nil, ctx.Err()
		}

	}
}

func (c *Client) TryLock(ctx context.Context, key string, expiration time.Duration) (*Lock, error) {
	val := uuid.New().String()
	ok, err := c.client.SetNX(ctx, key, val, expiration).Result()
	if err != nil {
		return nil, err
	}

	if !ok {
		// 代表的是别人抢到了锁
		return nil, ErrFailedToPreemptLock
	}

	return &Lock{
		client:     c.client,
		key:        key,
		value:      val,
		expiration: expiration,
	}, nil
}

type Lock struct {
	client     redis.Cmdable
	key        string
	value      string
	expiration time.Duration
	unlockChan chan struct{}
}

func (l *Lock) AutoRefresh(interval time.Duration, timeout time.Duration) error {

	timeoutChan := make(chan struct{}, 1)

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err := l.Refresh(ctx)
			cancel()
			if errors.Is(err, context.DeadlineExceeded) {
				timeoutChan <- struct{}{}
				continue
			}
			if err != nil {

				return err
			}
		case <-timeoutChan:

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()
			if errors.Is(err, context.DeadlineExceeded) {
				timeoutChan <- struct{}{}
				continue
			}
			if err != nil {
				return err
			}
		case <-l.unlockChan:
			return nil
		}
	}

}

func (l *Lock) Refresh(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaRefresh, []string{l.key}, l.value, l.expiration).Int64()

	if err == redis.Nil {
		return ErrLockNotHold
	}

	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotHold
	}

	return nil
}

func (l *Lock) Unlock(ctx context.Context) error {
	defer func() {
		//close(l.unlockChan)
		l.unlockChan <- struct{}{}
	}()

	res, err := l.client.Eval(ctx, luaUnlock, []string{l.key}, l.value).Int64()

	if err == redis.Nil {
		return ErrLockNotHold
	}

	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotHold
	}

	return nil
}

//func (l *Lock) Unlock(ctx context.Context) error {
//
//	cnt, err := l.client.Del(ctx, l.key).Result()
//	if err != nil {
//		return err
//	}
//
//	if cnt != 1 {
//		//代表你加的锁过期了
//
//		return errors.New("redis-lock: 解锁失败，锁不存在")
//	}
//
//	return nil
//}
