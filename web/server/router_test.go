package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func Test_router_AddRouter(t *testing.T) {
	testRouters := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user/order",
		},
		{
			method: http.MethodGet,
			path:   "/user/order/:id",
		},
		{
			method: http.MethodGet,
			path:   "/user/*",
		},
		{
			method: http.MethodPost,
			path:   "/",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
	}

	r := newRouter()
	mockHandlerFunc := func(ctx *Context) {}
	for _, testRouter := range testRouters {
		r.addRoute(testRouter.method, testRouter.path, mockHandlerFunc)
	}

	// 断言路由树和你预期的一模一样

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodPost: &node{
				path:    "/",
				handler: mockHandlerFunc,
				children: map[string]*node{
					"login": &node{
						path:    "login",
						handler: mockHandlerFunc,
					},
				},
			},
			http.MethodGet: &node{
				path:    "/",
				handler: mockHandlerFunc,
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: mockHandlerFunc,
						starChild: &node{
							path:    "*",
							handler: mockHandlerFunc,
						},
						children: map[string]*node{
							"home": &node{
								path:    "home",
								handler: mockHandlerFunc,
							},
							"order": &node{
								path:    "order",
								handler: mockHandlerFunc,
								paramChild: &node{
									path:    ":id",
									handler: mockHandlerFunc,
								},
							},
						},
					},
				},
			},
		},
	}

	msg, ok := r.equal(wantRouter)
	assert.True(t, ok, msg)

	newr := newRouter()
	assert.Panics(t, func() {
		//r.addRoute(http.MethodGet, "/", mockHandlerFunc)
		newr.addRoute(http.MethodGet, "", mockHandlerFunc)
	}, "\"\" 没有panic")

	assert.Panics(t, func() {
		//r.addRoute(http.MethodGet, "/", mockHandlerFunc)
		newr.addRoute(http.MethodGet, "/asdf/", mockHandlerFunc)
	}, "/asdf/ 没有panic")

	assert.Panics(t, func() {
		//r.addRoute(http.MethodGet, "/", mockHandlerFunc)
		newr.addRoute(http.MethodGet, "/asdf//asdf", mockHandlerFunc)
	}, "/asdf//asdf 没有panic")

	newr = newRouter()

	assert.Panics(t, func() {
		//r.addRoute(http.MethodGet, "/", mockHandlerFunc)
		newr.addRoute(http.MethodGet, "/asdf", mockHandlerFunc)
		newr.addRoute(http.MethodGet, "/asdf", mockHandlerFunc)
	}, "重复注册 没有panic")

	newr = newRouter()
	assert.Panics(t, func() {
		//r.addRoute(http.MethodGet, "/", mockHandlerFunc)
		newr.addRoute(http.MethodGet, "/", mockHandlerFunc)
		newr.addRoute(http.MethodGet, "/", mockHandlerFunc)
	}, "重复注册 没有panic")

	newr = newRouter()
	assert.Panics(t, func() {
		//r.addRoute(http.MethodGet, "/", mockHandlerFunc)
		newr.addRoute(http.MethodGet, "/a/:id", mockHandlerFunc)
		newr.addRoute(http.MethodGet, "/a/*", mockHandlerFunc)
	}, "重复注册通配符和路径参数")

	newr = newRouter()
	assert.Panics(t, func() {
		//r.addRoute(http.MethodGet, "/", mockHandlerFunc)
		newr.addRoute(http.MethodGet, "/a/*", mockHandlerFunc)
		newr.addRoute(http.MethodGet, "/a/:id", mockHandlerFunc)
	}, "重复注册通配符和路径参数")
}

func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return "找不到对应的 http method", false
		}

		msg, ok := v.equal(dst)
		if !ok {
			return msg, false
		}
	}

	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if y.path != n.path {
		return "节点路径不匹配", false
	}
	if len(n.children) != len(y.children) {
		return "子节点数量不相等", false
	}

	if n.starChild != nil {
		msg, ok := n.starChild.equal(y.starChild)
		if !ok {
			return msg, false
		}
	}

	if n.paramChild != nil {
		msg, ok := n.paramChild.equal(y.paramChild)
		if !ok {
			return msg, false
		}
	}

	nHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if nHandler != yHandler {
		return "handler 不相等", false
	}

	for path, c := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点 %s 不存在", path), false
		}

		msg, ok := c.equal(dst)

		if !ok {
			return msg, false
		}
	}

	return "", true
}

func TestRouter_findRoute(t *testing.T) {
	testRouter := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user/order",
		},
		{
			method: http.MethodGet,
			path:   "/user/*",
		},
		{
			method: http.MethodPost,
			path:   "/",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		{
			method: http.MethodPost,
			path:   "/login/:username",
		},
		{
			method: http.MethodGet,
			path:   "/user/order/:id/test",
		},
	}

	r := newRouter()

	var mockHandler HandlerFunc = func(ctx *Context) {

	}
	for _, route := range testRouter {
		r.addRoute(route.method, route.path, mockHandler)
	}

	testCase := []struct {
		name       string
		method     string
		path       string
		wantFound  bool
		wantNode   *node
		pathParams map[string]string
	}{
		{
			name:      "user abc",
			method:    http.MethodGet,
			path:      "/user/abc",
			wantFound: true,
			wantNode: &node{
				handler: mockHandler,
				path:    "*",
			},
		},
		{
			name:      "user order",
			method:    http.MethodGet,
			path:      "/user/order",
			wantFound: true,
			wantNode: &node{
				handler: mockHandler,
				path:    "order",
			},
		},
		{
			name:      "user",
			method:    http.MethodGet,
			path:      "/user",
			wantFound: true,
			wantNode: &node{
				handler: mockHandler,
				path:    "user",
				children: map[string]*node{
					"order": &node{
						handler: mockHandler,
						path:    "order",
					},
					"home": &node{
						handler: mockHandler,
						path:    "home",
					},
				},
				starChild: &node{
					path:    "*",
					handler: mockHandler,
				},
			},
		},
		{
			name:      "user/order/order123/test",
			method:    http.MethodGet,
			path:      "/user/order/order123/test",
			wantFound: true,
			wantNode: &node{
				handler: mockHandler,
				path:    "test",
			},
			pathParams: map[string]string{
				"id": "order123",
			},
		},
		{
			name:      "root",
			method:    http.MethodGet,
			path:      "/",
			wantFound: true,
			wantNode: &node{
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: mockHandler,
						children: map[string]*node{
							"home": &node{
								path:    "home",
								handler: mockHandler,
							},
							"order": &node{
								path:    "order",
								handler: mockHandler,
							},
						},
					},
				},
			},
		},
		{
			name:      "login",
			method:    http.MethodPost,
			path:      "/login",
			wantFound: true,
			wantNode: &node{
				path:    "login",
				handler: mockHandler,
				paramChild: &node{
					path:    ":username",
					handler: mockHandler,
				},
			},
		},
		{
			name:      "login username",
			method:    http.MethodPost,
			path:      "/login/aaaa",
			wantFound: true,
			wantNode: &node{
				path:    ":username",
				handler: mockHandler,
			},

			pathParams: map[string]string{
				"username": "aaaa",
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			matchInfo, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}

			if tc.pathParams != nil {
				for k, v := range tc.pathParams {
					pathParam, ok := matchInfo.pathParams[k]
					assert.True(t, ok, "路径参数不匹配")
					assert.Equal(t, v, pathParam, "路径参数不匹配")
				}
			}

			msg, ok := tc.wantNode.equal(matchInfo.n)

			assert.True(t, ok, msg)

		})
	}

}
