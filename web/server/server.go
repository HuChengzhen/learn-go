package server

import (
	"net"
	"net/http"
)

type HandlerFunc func(ctx *Context)

var _ Server = &HTTPServer{}

type Server interface {
	http.Handler

	Start(addr string) error
	//Start() error

	addRoute(method string, path string, f HandlerFunc)
}

type HTTPServer struct {
	// Addr string 创建的时候传递， 而不是接收参数
	*router
}

func NewHttpServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

func (h *HTTPServer) AddRoute1(method string, path string, handler ...HandlerFunc) {

}

func (h *HTTPServer) Get(path string, handler HandlerFunc) {
	h.addRoute(http.MethodGet, path, handler)
}

func (h *HTTPServer) Post(path string, handler HandlerFunc) {
	h.addRoute(http.MethodPost, path, handler)
}

func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	// 查找路由 执行命中的逻辑
	h.serve(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {
	route, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || route.n.handler == nil {
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("NOT FOUND"))
		return
	}
	ctx.pathParams = route.pathParams
	route.n.handler(ctx)
}

// Start 处理请求
func (h *HTTPServer) Start(addr string) error {
	// 可以自己创建server
	//http.Server{}
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// 在这里让用户注册回调
	// 生命周期回调

	return http.Serve(listen, h)
}

func (h *HTTPServer) Start2(addr string) error {
	return http.ListenAndServe(addr, h)
}

func NewServer() {

}
