package context

import (
	"context"
	"testing"
)

type myKey struct{}

func TestContext(t *testing.T) {
	// 一般是链路起点，或者调用的起点
	ctx := context.Background()
	//ctx := context.TODO()

	ctx = context.WithValue(ctx, myKey{}, "my-value")
	ctx, cancel := context.WithCancel(ctx)

	cancel()

	<-ctx.Done()
}
