package rpc

type Request struct {
	ServiceName string
	MethodName  string
	Args        any
}
