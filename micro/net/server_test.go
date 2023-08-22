package net

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"learn_geektime_go/micro/net/mocks"
	"net"
	"testing"
)

func Test_handleConn(t *testing.T) {
	type args struct {
		conn net.Conn
	}
	tests := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) net.Conn
		wantErr error
	}{
		{
			name: "read error",
			mock: func(ctrl *gomock.Controller) net.Conn {
				conn := mocks.NewMockConn(ctrl)
				conn.EXPECT().Read(gomock.Any()).Return(0, errors.New("read error"))
				return conn
			},
			wantErr: errors.New("read error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			err := handleConn(tt.mock(ctrl))

			assert.Equal(t, tt.wantErr, err)

			if err != nil {
				return
			}
		})
	}
}
