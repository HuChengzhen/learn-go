package net

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

func Serve(network, addr string) error {
	listener, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			err2 := handleConn(conn)
			if err2 != nil {
				_ = conn.Close()
			}
		}()
	}
}

func handleConn(conn net.Conn) error {
	for {
		bs := make([]byte, 8)
		n, err := conn.Read(bs)

		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return err
		}

		if err != nil {
			return err
		}

		//if n != 8 {
		//	return errors.New("micro: 没读够数据")
		//}

		res := handleMsg(bs)

		n, err = conn.Write(res)
		if err != nil {
			return err
		}

		if n != len(res) {
			return errors.New("micro: 没写完数据")
		}
	}
}

func handleMsg(req []byte) []byte {
	res := make([]byte, 2*len(req))
	copy(res[:len(req)], req)
	copy(res[len(req):], req)
	return res
}

type Server struct {
	network string
	addr    string
}

func NewServer(network, addr string) *Server {
	return &Server{
		network: network,
		addr:    addr,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen(s.network, s.addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			err2 := s.handleConn(conn)
			if err2 != nil {
				_ = conn.Close()
			}
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	for {
		lenBs := make([]byte, 8)
		n, err := conn.Read(lenBs)

		if err != nil {
			return err
		}

		length := binary.BigEndian.Uint64(lenBs)

		resBs := make([]byte, length)

		resData := handleMsg(resBs)
		respLen := len(resData)

		res := make([]byte, 8+respLen)

		binary.BigEndian.PutUint64(res[:8], uint64(respLen))

		copy(res[8:], resData)

		n, err = conn.Write(res)
		if err != nil {
			return err
		}

		if n != len(res) {
			return errors.New("micro: 没写完数据")
		}
	}
}
