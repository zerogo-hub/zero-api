package app

import (
	"errors"
	"net/http"
	"net/url"
	_path "path"
	"path/filepath"
	"sync"

	zeroapi "github.com/zerogo-hub/zero-api"
	zeroctx "github.com/zerogo-hub/zero-api/context"
	zerorouter "github.com/zerogo-hub/zero-api/router"
	zeroserver "github.com/zerogo-hub/zero-api/server"

	zerologger "github.com/zerogo-hub/zero-helper/logger"

	zamcors "github.com/zerogo-hub/zero-api-middleware/cors"
	zamlogger "github.com/zerogo-hub/zero-api-middleware/logger"
)

type app struct {
	// router 路由管理器
	router zeroapi.Router

	// server http 服务器
	server zeroapi.Server

	// context 对象池
	ctxPool *sync.Pool

	// config 应用配置
	config *config

	// middlewares App级别 中间件
	middlewares []zeroapi.Handler
}

// New 生成一个应用实例
func New() zeroapi.App {
	a := NewApp()
	return a
}

// Default 生成默认的应用实例
func Default() zeroapi.App {
	a := NewApp()

	// 请求日志
	a.Use(zamlogger.New())
	// 跨域
	a.Use(zamcors.New(nil))

	return a
}

// NewApp 生成一个应用实例
func NewApp(opts ...Option) zeroapi.App {
	a := &app{
		ctxPool: &sync.Pool{},
		config:  defaultConfig(),
	}

	a.router = zerorouter.NewRouter(a)
	a.server = zeroserver.NewServer(a)
	a.ctxPool.New = func() interface{} {
		return zeroctx.NewContext(a)
	}

	for _, opt := range opts {
		opt(a.config)
	}

	return a
}

// Router 获取路由管理示例
func (a *app) Router() zeroapi.Router {
	return a.router
}

// Server http 服务器
func (a *app) Server() zeroapi.Server {
	return a.server
}

// Context 从 pool 中获取一个 Context
func (a *app) Context() zeroapi.Context {
	return a.ctxPool.Get().(zeroapi.Context)
}

// ReleaseContext 将 Context 释放到 pool 中
func (a *app) ReleaseContext(ctx zeroapi.Context) {
	a.ctxPool.Put(ctx)
}

// Version 获取框架版本号
func (a *app) Version() string {
	return a.config.version
}

// Logger 获取日志实例
func (a *app) Logger() zerologger.Logger {
	return a.config.logger
}

// FileMaxMemory 文件系统使用的最大内存
func (a *app) FileMaxMemory() int64 {
	return a.config.fileMaxMemory
}

// IsCookieEncode cookie 是否需要进行编码
func (a *app) IsCookieEncode() bool {
	return a.config.cookieEncode != nil && a.config.cookieDecode != nil
}

// CookieEncodeHandler 获取 cookie 编码函数
func (a *app) CookieEncodeHandler() zeroapi.CookieEncodeHandler {
	return a.config.cookieEncode
}

// CookieDecodeHandler 获取 cookie 解码函数
func (a *app) CookieDecodeHandler() zeroapi.CookieDecodeHandler {
	return a.config.cookieDecode
}

// Use 添加 App 级别 中间件，每一次路由都会调用公共中间件
func (a *app) Use(handlers ...zeroapi.Handler) {
	for _, handler := range handlers {
		if handler != nil {
			a.middlewares = append(a.middlewares, handler)
		}
	}
}

// ExecuteMiddlewares 执行 App 级别的中间件
func (a *app) ExecuteMiddlewares(ctx zeroapi.Context) {
	for _, handler := range a.middlewares {
		handler(ctx)
		if ctx.IsStopped() {
			return
		}
	}
}

// Run 启动服务，此方法会阻塞，直到应用关闭
// addr: host:port，例如: ":8080"，"192.168.1.8:80"
func (a *app) Run(addr string) error {
	if !a.Router().Build() {
		return errors.New("router build failed")
	}

	if err := a.server.Start(addr); err != nil {
		if err == http.ErrServerClosed {
			a.Logger().Info(http.ErrServerClosed.Error())
		} else {
			a.Logger().Error(err.Error())
		}
	}

	return nil
}

// Prefix 设置前缀，设置前就已添加的路由不会有该前缀
// 例如: prefix = "/blog"，则 "/user" -> "/blog/user"
func (a *app) Prefix(prefix string) zeroapi.App {
	a.router.Prefix(prefix)
	return a
}

// Get method = "GET"
// path: 路径，以 "/" 开头，不可以为空
// handlers: 路由级别中间件和处理函数
func (a *app) Get(path string, handlers ...zeroapi.Handler) zeroapi.App {
	a.router.Register(zeroapi.MethodGet, path, handlers...)
	return a
}

// Post method = "POST"
// path: 路径，以 "/" 开头，不可以为空
// handlers: 路由级别中间件和处理函数
func (a *app) Post(path string, handlers ...zeroapi.Handler) zeroapi.App {
	a.router.Register(zeroapi.MethodPost, path, handlers...)
	return a
}

// Put method = "PUT"
// path: 路径，以 "/" 开头，不可以为空
// handlers: 路由级别中间件和处理函数
func (a *app) Put(path string, handlers ...zeroapi.Handler) zeroapi.App {
	a.router.Register(zeroapi.MethodPost, path, handlers...)
	return a
}

// Delete method = "DELETE"
// path: 路径，以 "/" 开头，不可以为空
// handlers: 路由级别中间件和处理函数
func (a *app) Delete(path string, handlers ...zeroapi.Handler) zeroapi.App {
	a.router.Register(zeroapi.MethodDelete, path, handlers...)
	return a
}

// Head method = "HEAD"
// path: 路径，以 "/" 开头，不可以为空
// handlers: 路由级别中间件和处理函数
func (a *app) Head(path string, handlers ...zeroapi.Handler) zeroapi.App {
	a.router.Register(zeroapi.MethodHead, path, handlers...)
	return a
}

// Patch method = "PATCH"
// path: 路径，以 "/" 开头，不可以为空
// handlers: 路由级别中间件和处理函数
func (a *app) Patch(path string, handlers ...zeroapi.Handler) zeroapi.App {
	a.router.Register(zeroapi.MethodPatch, path, handlers...)
	return a
}

// Options method = "OPTIONS"
// path: 路径，以 "/" 开头，不可以为空
// handlers: 路由级别中间件和处理函数
func (a *app) Options(path string, handlers ...zeroapi.Handler) zeroapi.App {
	a.router.Register(zeroapi.MethodOptions, path, handlers...)
	return a
}

// Group 创建组路由实例
func (a *app) Group(path string) zeroapi.Group {
	return zerorouter.NewGroup(a, path)
}

// Static 添加静态资源服务
// prefix 静态资源路由前缀
// path 资源真实位置(绝对路径，相对路径)
func (a *app) Static(prefix, path string) {
	if path == "" {
		path = "."
	}

	f := func(ctx zeroapi.Context) {
		fileName, err := url.PathUnescape(ctx.Dynamic("*"))
		if fileName == "" || err != nil {
			ctx.NotFound()
			return
		}
		path := filepath.Join(path, _path.Clean("/"+fileName))
		ctx.DownloadFile(path, fileName)
	}

	if prefix == "/" {
		a.Get(prefix+"*", f)
		return
	}

	a.Get(prefix+"/*", f)
}
