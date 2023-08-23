package net

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

type Pool struct {
	// 空闲链接队列
	idlesConns chan *idleConn
	// 请求队列
	reqQueue []connReq

	maxCnt int

	cnt int

	maxIdleTime time.Duration

	initCnt int

	factory func() (net.Conn, error)

	lock sync.Mutex
}

func NewPool(
	initCnt int,
	maxIdleCnt int,
	maxCnt int,
	maxIdletime time.Duration,
	factory func() (net.Conn, error),
) (*Pool, error) {
	if initCnt > maxIdleCnt {
		return nil, errors.New("micro: 初始连接数量不能大于最大空闲连接数")
	}

	idlesConns := make(chan *idleConn, maxIdleCnt)

	for i := 0; i < initCnt; i++ {
		c, err := factory()
		if err != nil {
			return nil, err
		}
		idlesConns <- &idleConn{c: c, lastActiveTime: time.Now()}
	}

	res := &Pool{
		idlesConns:  idlesConns,
		maxCnt:      maxCnt,
		cnt:         0,
		maxIdleTime: maxIdletime,
		initCnt:     initCnt,
		factory:     factory,
	}

	return res, nil
}

func (p *Pool) Get(ctx context.Context) (net.Conn, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	for {
		select {
		case ic := <-p.idlesConns:

			if ic.lastActiveTime.Add(p.maxIdleTime).Before(time.Now()) {
				_ = ic.c.Close()
				continue
			}

			return ic.c, nil
		default:
			p.lock.Lock()
			if p.cnt >= p.maxCnt {
				req := connReq{connChan: make(chan net.Conn)}
				p.reqQueue = append(p.reqQueue, req)
				p.lock.Unlock()
				select {
				case <-ctx.Done():
					go func() {
						c := <-req.connChan
						_ = p.Put(context.Background(), c)
					}()

					return nil, ctx.Err()
				case c := <-req.connChan:
					return c, nil
				}

			}

			c, err := p.factory()
			if err != nil {
				return nil, err
			}

			p.cnt++
			p.lock.Unlock()
			return c, nil

		}
	}
}

func (p *Pool) Put(ctx context.Context, conn net.Conn) error {
	p.lock.Lock()
	if len(p.reqQueue) > 0 {
		p.reqQueue = p.reqQueue[1:]
		p.lock.Unlock()
		p.reqQueue[0].connChan <- conn
		return nil
	}

	defer p.lock.Unlock()

	ic := &idleConn{
		c:              conn,
		lastActiveTime: time.Now(),
	}

	select {
	case p.idlesConns <- ic:

	default:
		_ = conn.Close()
		// p.lock.Lock()
		p.cnt--
		// p.lock.Unlock()
	}

	return nil
}

type idleConn struct {
	c              net.Conn
	lastActiveTime time.Time
}

type connReq struct {
	connChan chan net.Conn
}
