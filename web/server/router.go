package server

import "strings"

type router struct {
	// http method => 路由树
	trees map[string]*node
}

func newRouter() *router {
	return &router{
		trees: make(map[string]*node),
	}
}

type node struct {
	path string

	// 子path到子节点的映射
	children map[string]*node

	// 缺一个代表用户注册的业务逻辑

	starChild *node

	paramChild *node

	handler HandlerFunc
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {

	n, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return &matchInfo{
			n: n,
		}, true
	}

	path = strings.Trim(path, "/")

	segs := strings.Split(path, "/")

	var pathParams map[string]string
	for _, seg := range segs {
		child, param, found := n.childOf(seg)
		if !found {
			return nil, false
		}

		if param {
			if pathParams == nil {
				pathParams = map[string]string{}
				pathParams[child.path[1:]] = seg
			}
		}

		n = child
	}

	return &matchInfo{
		n:          n,
		pathParams: pathParams,
	}, true
}

func (r *router) addRoute(method string, path string, handler HandlerFunc) {
	if path == "" {
		panic("web 路径不能为空字符串")
	}

	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	if path[0] != '/' {
		panic("web: 路径必须以 / 开头")
	}

	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路径不能以 / 结尾")
	}

	if path == "/" {

		if root.handler != nil {
			panic("web: 路由冲突")
		}

		root.handler = handler
		return
	}

	split := strings.Split(path[1:], "/")
	child := root
	for _, seg := range split {
		if seg == "" {
			panic("web: 中间不能有连续//")
		}
		child = child.childOrCreate(seg)
	}

	if child.handler != nil {
		panic("web: 路由冲突")
	}
	child.handler = handler
}

func (n *node) childOrCreate(seg string) *node {
	if seg[0] == ':' {
		if n.starChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配，已有通配符匹配")
		}

		n.paramChild = &node{
			path: seg,
		}
		return n.paramChild
	}

	if seg == "*" {
		if n.paramChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配，已有路径参数匹配")
		}

		if n.starChild == nil {
			n.starChild = &node{
				path: seg,
			}
		}
		return n.starChild
	}

	if n.children == nil {
		n.children = map[string]*node{}
	}

	res, ok := n.children[seg]
	if !ok {
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}

	return res
}

// childOf 优先考虑静态匹配，匹配不上，才考虑通配符匹配
// 第一个返回值是子节点，第二个返回值是标记是否是路径参数，第三个标记是否命中
func (n *node) childOf(seg string) (*node, bool, bool) {

	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}

		return n.starChild, false, n.starChild != nil
	}

	child, ok := n.children[seg]

	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}

		return n.starChild, false, n.starChild != nil
	}
	return child, false, ok
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}
