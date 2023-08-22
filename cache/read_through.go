package cache

import (
	"context"
	"errors"
	"golang.org/x/sync/singleflight"
	"time"
)

type ReadThroughCache struct {
	Cache
	LoadFunc   func(ctx context.Context, key string) (any, error)
	Expiration time.Duration
	g          *singleflight.Group
}

func (r *ReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if errors.Is(err, ErrKeyNotFound) {
		val, err = r.LoadFunc(ctx, key)
		if err == nil {
			err := r.Set(ctx, key, val, r.Expiration)
			if err != nil {

			}
		}
	}

	return val, err
}

func (r *ReadThroughCache) GetV3(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if errors.Is(err, ErrKeyNotFound) {
		val, err, _ = r.g.Do(key, func() (interface{}, error) {
			v, er := r.LoadFunc(ctx, key)
			if er == nil {
				er := r.Set(ctx, key, val, r.Expiration)
				if er != nil {

				}
			}
			return v, er
		})
	}

	return val, err
}
