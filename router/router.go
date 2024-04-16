package router

import (
	zeroapi "github.com/zerogo-hub/zero-api"
)

type router struct {
	// app 应用实例
	app zeroapi.App

	// prefix 路由前缀
	prefix string

	// routes 按照 Method 存储路由
	routes map[string]Route

	// validators 存储验证函数
	validators map[string]zeroapi.RouterValidator
}

// NewRouter 创建一个 zeroapi.Router 实例
func NewRouter(app zeroapi.App) zeroapi.Router {
	return &router{
		app:        app,
		routes:     make(map[string]Route, len(zeroapi.AllMethods())),
		validators: make(map[string]zeroapi.RouterValidator),
	}
}

// Prefix 设置前缀，设置前就已添加的路由不会有该前缀
// 例如: prefix = "/blog"，则 "/user" -> "/blog/user"
func (r *router) Prefix(prefix string) {
	if len(prefix) == 0 {
		return
	}

	if prefix[0] != '/' {
		prefix = "/" + prefix
	}

	r.prefix = prefix
}

// Register 注册路由处理函数，以及中间件
// method: HTTP Method，见 core/const.go Methodxxxx
// path: 路径，以 "/" 开头，不可以为空
// handles: 处理函数和路由级别中间件，匹配成功后会调用该函数
func (r *router) Register(method, path string, handlers ...zeroapi.Handler) bool {
	if len(path) == 0 {
		return false
	} else if len(handlers) == 0 {
		return false
	}

	if path[0] != '/' {
		path = "/" + path
	}

	if r.prefix != "" {
		path = r.prefix + "/" + path
	}

	re := r.routes[method]
	if re == nil {
		re = NewRoute()
		r.routes[method] = re
	}

	re.Insert(path, handlers...)

	return true
}

// Build 解析路由，包括动态参数，正则表达式，验证函数的解析，路由路径查找优化
func (r *router) Build() bool {
	for _, re := range r.routes {
		if !re.Build(r) {
			return false
		}
	}

	return true
}

// Lookup 查找路由
func (r *router) Lookup(method, path string) ([]zeroapi.Handler, map[string]string) {
	if re := r.routes[method]; re != nil {
		return re.Lookup(path)
	}

	return nil, nil
}

// RegisterRouterValidator 注册路由验证函数
func (r *router) RegisterRouterValidator(name string, validator zeroapi.RouterValidator) {
	if _, exist := r.validators[name]; exist {
		return
	}

	r.validators[name] = validator
}

// Validator 获取路由验证函数
func (r *router) Validator(name string) zeroapi.RouterValidator {
	if f, exist := r.validators[name]; exist {
		return f
	}

	return nil
}
