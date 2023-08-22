package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"learn_geektime_go/cache/mocks"
	"log"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) redis.Cmdable
		key      string
		wantErr  error
		wantLock *Lock
	}{
		{
			name: "set nx error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				res := redis.NewBoolResult(false, context.DeadlineExceeded)
				cmd.EXPECT().SetNX(context.Background(), "key1", gomock.Any(), time.Minute).Return(res)
				return cmd
			},
			key: "key1",
			wantLock: &Lock{
				key: "key1",
			},
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			client := NewClient(tc.mock(ctrl))

			lock, err := client.TryLock(context.Background(), tc.key, time.Minute)

			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantLock.key, lock.key)
			assert.NotEmpty(t, lock.value)

		})
	}
}

func TestLock_Unlock(t *testing.T) {

}

func ExampleLock_Refresh() {
	var l *Lock
	stopChan := make(chan struct{})
	errChan := make(chan error)
	timeoutChan := make(chan struct{}, 1)
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		var timeoutRetry = 0
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
					errChan <- err
					close(stopChan)
					close(errChan)
					return
				}
				timeoutRetry = 0
			case <-timeoutChan:
				timeoutRetry++

				if timeoutRetry > 10 {
					errChan <- context.DeadlineExceeded
					return
				}
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				err := l.Refresh(ctx)
				cancel()
				if errors.Is(err, context.DeadlineExceeded) {
					timeoutChan <- struct{}{}
					continue
				}
				if err != nil {
					errChan <- err
					close(stopChan)
					close(errChan)
					return
				}
			case <-stopChan:
				return
			}
		}

	}()

	for i := 0; i < 10; i++ {
		select {
		case <-errChan:
			break
		default:
			// 正常的业务逻辑
		}
	}

	select {
	case err := <-errChan:
		log.Printf(" %v", err)
		break
	default:
		// 正常的业务逻辑1
	}

	select {
	case err := <-errChan:
		log.Printf(" %v", err)
		break
	default:
		// 正常的业务逻辑2
	}

	stopChan <- struct{}{}

	l.Unlock(context.Background())
	// Output:
	// Hello
}

func ExampleLock_AutoRefresh() {
	var l *Lock
	go func() {
		// 这里返回error 需要终端业务
		l.AutoRefresh(time.Second*10, time.Second)
	}()

	fmt.Println("Hello")
}
