package cache

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"learn_geektime_go/cache/mocks"
	"testing"
	"time"
)

func TestRedisCache_Get(t *testing.T) {
}

func TestRedisCache_Set(t *testing.T) {

	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) redis.Cmdable
		key        string
		value      string
		expiration time.Duration

		wantErr error
	}{
		{
			name: "",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdable := mocks.NewMockCmdable(ctrl)

				status := redis.NewStatusCmd(context.Background())
				status.SetVal("OK")

				cmdable.EXPECT().
					Set(context.Background(), "key1", "value1", time.Second).
					Return(status)
				return cmdable
			},
			key:        "key1",
			value:      "value1",
			expiration: time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			c := NewRedisCache(tc.mock(ctrl))
			err := c.Set(context.Background(), tc.key, tc.value, tc.expiration)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
		})
	}
}
