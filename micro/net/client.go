package net

import (
	"fmt"
	"net"
	"time"
)

func Connect(network, addr string) error {
	conn, err := net.DialTimeout(network, addr, time.Minute*3)
	if err != nil {
		return err
	}

	defer conn.Close()

	for {
		_, err = conn.Write([]byte("hello"))
		if err != nil {
			return err
		}
		res := make([]byte, 128)

		_, err := conn.Read(res)

		if err != nil {
		}
		fmt.Println(string(res))
	}
}

type Client struct {
	network string
	addr    string
}

func NewClient(network string, addr string) *Client {
	return &Client{
		network: network,
		addr:    addr,
	}
}

func (c *Client) Send(data string) (string, error) {
	conn, err := net.DialTimeout(c.network, c.addr, time.Minute*3)
	if err != nil {
		return "", err
	}

	defer conn.Close()

	for {
		_, err = conn.Write([]byte("hello"))
		if err != nil {
			return "", err
		}
		res := make([]byte, 128)

		_, err := conn.Read(res)

		if err != nil {
		}
		fmt.Println(string(res))
	}
}
