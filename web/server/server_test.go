package server

// func TestServer(t *testing.T) {
// 	h := NewHttpServer()
// 	//http.ListenAndServe(":8080", h)
// 	//http.ListenAndServeTLS(":443", "", "", h)

// 	handler1 := func(ctx *Context) {
// 		ctx.Resp.Write([]byte("Hello asdfsadf"))
// 	}

// 	h.addRoute(http.MethodGet, "/user", func(ctx *Context) {
// 		handler1(ctx)
// 	})

// 	h.addRoute(http.MethodGet, "/user/*", func(ctx *Context) {
// 		ctx.Resp.Write([]byte("/user/*"))
// 	})

// 	h.addRoute(http.MethodGet, "/user/order", func(ctx *Context) {
// 		ctx.Resp.Write([]byte("/user/order"))
// 	})

// 	h.Start(":8080")
// }
