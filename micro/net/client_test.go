package net

import (
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	go func() {
		Serve("tcp", ":8082")
	}()

	time.Sleep(time.Second * 3)

	err := Connect("tcp", "localhost:8082")

	t.Log(err)
}
