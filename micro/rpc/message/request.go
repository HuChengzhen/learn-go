package message

import "encoding/binary"

type Request struct {
	HeadLength uint32
	BodyLength uint32
	RequestId  uint32
	Version    uint8
	Compresser uint8
	Serializer uint8

	ServicName string
	MethodName string

	Meta map[string]string

	Data []byte
}

func EncodeReq(req *Request) []byte {
	bs := make([]byte, req.HeadLength+req.BodyLength)
	binary.BigEndian.PutUint32(bs, req.HeadLength)
	binary.BigEndian.PutUint32(bs[4:], req.BodyLength)

	binary.BigEndian.PutUint32(bs[8:12], req.RequestId)

	bs[12] = req.Version
	bs[13] = req.Compresser
	bs[14] = req.Serializer

	copy(bs[15:15+len(req.ServicName)], req.ServicName)

	return bs
}

func DecodeReq(data []byte) *Request {
	req := &Request{}
	req.HeadLength = binary.BigEndian.Uint32(data[:4])
	req.BodyLength = binary.BigEndian.Uint32(data[4:8])
	req.RequestId = binary.BigEndian.Uint32(data[8:12])
	req.Version = data[12]
	req.Compresser = data[13]
	req.Serializer = data[14]
	req.ServicName = string(data[15:])
	return req
}
