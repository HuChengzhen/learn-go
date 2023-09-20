package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	testCases := []struct {
		name string
		req  *Request
	}{
		{
			name: "normal",
			req: &Request{
				RequestId:  1,
				Version:    12,
				Compresser: 14,
				Serializer: 4,
				ServicName: "user-service",
				// MethodName: "GetById",
				// Meta: map[string]string{
				// 	"trace-id": "123456",
				// 	"a/b":      "a",
				// },
				// Data: []byte("hello, world"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.req.calculateHeaderLength()
			tc.req.calculateBodyLength()

			data := EncodeReq(tc.req)
			req := DecodeReq(data)

			assert.Equal(t, req, tc.req)
		})
	}
}

func (req *Request) calculateHeaderLength() {
	req.HeadLength = 15 + uint32(len(req.ServicName))
}

func (req *Request) calculateBodyLength() {
	req.BodyLength = uint32(len(req.Data))
}
