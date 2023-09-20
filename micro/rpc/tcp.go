package rpc

import (
	"encoding/binary"
	"net"
)

func ReadMsg(conn net.Conn) ([]byte, error) {
	lenBs := make([]byte, numOfLengthBytes)

	_, err := conn.Read(lenBs)
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint64(lenBs)

	respBs := make([]byte, length)

	_, err = conn.Read(respBs)
	return respBs, err
}

func EncodeMsg(data []byte) []byte {
	reqLen := len(data)

	req := make([]byte, reqLen+numOfLengthBytes)

	binary.BigEndian.PutUint64(req[:numOfLengthBytes], uint64(reqLen))

	copy(req[numOfLengthBytes:], data)

	return req
}
