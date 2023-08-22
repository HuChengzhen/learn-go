package channel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBroker_Send(t *testing.T) {
	b := &Broker{}

	go func() {
		for {
			err := b.Send(Msg{Content: time.Now().String()})
			if err != nil {
				t.Log(err)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	var wg = sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			msgs, err := b.Subscribe(100)
			if err != nil {
				return
			}

			for msg := range msgs {
				fmt.Println(msg.Content)
			}
		}()
	}
	wg.Wait()
}
