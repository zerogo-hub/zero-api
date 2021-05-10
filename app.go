package zeroapi

import (
	"errors"
	"net/http"
	"net/url"
	_path "path"
	"path/filepath"
	"sync"

	"github.com/zerogo-hub/zero-helper/logger"
)

// App 应用
type App interface {

	// Router 获取路由管理示例
	Router() Router

	// Server http 服务器
	Server() Server

	// Context 从 pool 中获取一个 Context
	Context() Context

	// ReleaseContext 将 Context 释放到 pool 中
	ReleaseContext(ctx Context)

	// Version 获取框架版本号
	Version() string

	// Logger 获取日志实例
	Logger() logger.Logger

	// FileMaxMemory 文件系统使用的最大内存
	FileMaxMemory() int64

	// IsCookieEncode cookie 是否需要进行编码
	IsCookieEncode() bool

	// CookieEncodeHandler 获取 cookie 编码函数
	CookieEncodeHandler() CookieEncodeHandler

	// CookieDecodeHandler 获取 cookie 解码函数
	CookieDecodeHandler() CookieDecodeHandler

	// Use 添加 App 级别 中间件，每一次路由都会调用公共中间件
	Use(handlers ...Handler)

	// ExecuteMiddlewares 执行 App 级别的中间件
	ExecuteMiddlewares(ctx Context)

	// Run 启动服务，此方法会阻塞，直到应用关闭
	// addr: host:port，例如: ":8080"，"192.168.1.8:80"
	Run(addr string) error

	RouterRegister
}

// RouterRegister 路由注册相关接口
type RouterRegister interface {
	// Prefix 设置前缀，设置前就已添加的路由不会有该前缀
	// 例如: prefix = "/blog"，则 "/user" -> "/blog/user"
	Prefix(prefix string) App

	// Get method = "GET"
	// path: 路径，以 "/" 开头，不可以为空
	// handlers: 路由级别中间件和处理函数
	Get(path string, handlers ...Handler) App

	// Post method = "POST"
	// path: 路径，以 "/" 开头，不可以为空
	// handlers: 路由级别中间件和处理函数
	Post(path string, handlers ...Handler) App

	// Put method = "PUT"
	// path: 路径，以 "/" 开头，不可以为空
	// handlers: 路由级别中间件和处理函数
	Put(path string, handlers ...Handler) App

	// Delete method = "DELETE"
	// path: 路径，以 "/" 开头，不可以为空
	// handlers: 路由级别中间件和处理函数
	Delete(path string, handlers ...Handler) App

	// Head method = "HEAD"
	// path: 路径，以 "/" 开头，不可以为空
	// handlers: 路由级别中间件和处理函数
	Head(path string, handlers ...Handler) App

	// Patch method = "PATCH"
	// path: 路径，以 "/" 开头，不可以为空
	// handlers: 路由级别中间件和处理函数
	Patch(path string, handlers ...Handler) App

	// Options method = "OPTIONS"
	// path: 路径，以 "/" 开头，不可以为空
	// handlers: 路由级别中间件和处理函数
	Options(path string, handlers ...Handler) App

	// Group 创建组路由实例
	Group(path string) Group

	// Static 添加静态资源服务
	// prefix 静态资源路由前缀
	// path 资源真实位置(绝对路径，相对路径)
	Static(prefix, path string)
}

type app struct {
	// router 路由管理器
	router Router

	// server http 服务器
	server Server

	// context 对象池
	ctxPool *sync.Pool

	// config 应用配置
	config *config

	// middlewares App级别 中间件
	middlewares []Handler
}

// NewApp 生成一个应用实例
func NewApp(opts ...Option) App {
	a := &app{
		ctxPool: &sync.Pool{},
		config:  defaultConfig(),
	}

	a.router = NewRouter(a)
	a.server = NewServer(a)
	a.ctxPool.New = func() interface{} {
		return NewContext(a)
	}

	for _, opt := range opts {
		opt(a.config)
	}

	return a
}

// Router 获取路由管理示例
func (a *app) Router() Router {
	return a.router
}

// Server http 服务器
func (a *app) Server() Server {
	return a.server
}

// Context 从 pool 中获取一个 Context
func (a *app) Context() Context {
	return a.ctxPool.Get().(Context)
}

// ReleaseContext 将 Context 释放到 pool 中
func (a *app) ReleaseContext(ctx Context) {
	a.ctxPool.Put(ctx)
}

// Version 获取框架版本号
func (a *app) Version() string {
	return a.config.version
}

// Logger 获取日志实例
func (a *app) Logger() logger.Logger {
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
func (a *app) CookieEncodeHandler() CookieEncodeHandler {
	return a.config.cookieEncode
}

// CookieDecodeHandler 获取 cookie 解码函数
func (a *app) CookieDecodeHandler() CookieDecodeHandler {
	return a.config.cookieDecode
}

// Use 添加 App 级别 中间件，每一次路由都会调用公共中间件
func (a *app) Use(handlers ...Handler) {
	for _, handler := range handlers {
		if handler != nil {
			a.middlewares = append(a.middlewares, handler)
		}
	}
}

// ExecuteMiddlewares 执行 App 级别的中间件
func (a *app) ExecuteMiddlewares(ctx Context) {
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
func (a *app) Prefix(prefix string) App {
	a.router.Prefix(prefix)
	return a
}

// Get method = "GET"
// path: 路径，以 "/" 开头，不可以为空
// handlers: 路由级别中间件和处理函数
func (a *app) Get(path string, handlers ...Handler) App {
	a.router.Register(MethodGet, path, handlers...)
	return a
}

// Post method = "POST"
// path: 路径，以 "/" 开头，不可以为空
// handlers: 路由级别中间件和处理函数
func (a *app) Post(path string, handlers ...Handler) App {
	a.router.Register(MethodPost, path, handlers...)
	return a
}

// Put method = "PUT"
// path: 路径，以 "/" 开头，不可以为空
// handlers: 路由级别中间件和处理函数
func (a *app) Put(path string, handlers ...Handler) App {
	a.router.Register(MethodPost, path, handlers...)
	return a
}

// Delete method = "DELETE"
// path: 路径，以 "/" 开头，不可以为空
// handlers: 路由级别中间件和处理函数
func (a *app) Delete(path string, handlers ...Handler) App {
	a.router.Register(MethodDelete, path, handlers...)
	return a
}

// Head method = "HEAD"
// path: 路径，以 "/" 开头，不可以为空
// handlers: 路由级别中间件和处理函数
func (a *app) Head(path string, handlers ...Handler) App {
	a.router.Register(MethodHead, path, handlers...)
	return a
}

// Patch method = "PATCH"
// path: 路径，以 "/" 开头，不可以为空
// handlers: 路由级别中间件和处理函数
func (a *app) Patch(path string, handlers ...Handler) App {
	a.router.Register(MethodPatch, path, handlers...)
	return a
}

// Options method = "OPTIONS"
// path: 路径，以 "/" 开头，不可以为空
// handlers: 路由级别中间件和处理函数
func (a *app) Options(path string, handlers ...Handler) App {
	a.router.Register(MethodOptions, path, handlers...)
	return a
}

// Group 创建组路由实例
func (a *app) Group(path string) Group {
	return NewGroup(a, path)
}

// Static 添加静态资源服务
// prefix 静态资源路由前缀
// path 资源真实位置(绝对路径，相对路径)
func (a *app) Static(prefix, path string) {
	if path == "" {
		path = "."
	}

	f := func(ctx Context) {
		fileName, err := url.PathUnescape(ctx.Dynamic("*"))
		if fileName == "" || err != nil {
			ctx.NotFound()
			return
		}
		path := filepath.Join(path, _path.Clean("/"+fileName))
		ctx.DownloadFile(path, fileName)
		return
	}

	if prefix == "/" {
		a.Get(prefix+"*", f)
		return
	}

	a.Get(prefix+"/*", f)
}
