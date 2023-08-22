package channel

import (
	"errors"
	"sync"
)

type Broker struct {
	mutex sync.RWMutex
	chans []chan Msg
}

func (b *Broker) Send(m Msg) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for _, msgs := range b.chans {
		select {
		case msgs <- m:
		default:
			return errors.New("消息队列已满")
		}
	}
	return nil
}

func (b *Broker) Subscribe(cap int) (<-chan Msg, error) {
	res := make(chan Msg, cap)
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.chans = append(b.chans, res)
	return res, nil
}

func (b *Broker) Close() error {
	b.mutex.Lock()
	chans := b.chans
	b.chans = nil
	b.mutex.Unlock()

	// 避免了重复close channel的问题
	for _, msgs := range chans {
		close(msgs)
	}

	return nil
}

type Msg struct {
	Content string
}

type BrokerV2 struct {
	mutex     sync.RWMutex
	consumers []func(msg Msg)
	//	consumers map[string][]func(msg Msg)
}

func (b *BrokerV2) Send(m Msg) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for _, consumer := range b.consumers {
		consumer(m)
	}

	return nil
}

func (b *BrokerV2) Subscribe(cb func(m Msg)) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.consumers = append(b.consumers, cb)
	return nil
}
