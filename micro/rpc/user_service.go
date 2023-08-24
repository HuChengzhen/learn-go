package rpc

import "context"

type UserService struct {
	// 用反射来赋值
	//； 类型是函数的字段，他不是方法。
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

func (u *UserService) Name() string {
	return "user-service"
}

type GetByIdReq struct {
	Id int
}

type GetByIdResp struct {
}
