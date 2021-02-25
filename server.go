package zeroapi

import (
	"net/http"
	"os"
	"sync"

	"github.com/zerogo-hub/zero-helper/file"
	graceful "github.com/zerogo-hub/zero-helper/graceful/http"
)

// Server http 服务器
type Server interface {
	http.Handler

	// App 返回应用实例
	App() App

	// SetTLS 指定 tls 证书，密钥路径
	// certFile: 证书路径
	// keyFile: 私钥路径
	SetTLS(certFile, keyFile string) bool

	// IsTLS 是否启用 tls
	IsTLS() bool

	// Start 根据配置调用 ListenAndServe 或者 ListenAndServeTLS，接收连接请求
	// addr: host:port，例如: ":8080"，"192.168.1.8:80"
	Start(addr string) error

	// Context 从 pool 中获取一个 Context
	Context() Context

	// ReleaseContext 将 Context 释放到 pool 中
	ReleaseContext(ctx Context)

	// HTTPServer 实际使用的 http 服务器
	HTTPServer() graceful.Server
}

type server struct {
	// app 应用实例
	app App

	// httpServer 实际使用  graceful.Server 替代 http.Server
	httpServer graceful.Server

	// tlsCertFile tls 证书路径
	tlsCertFile string

	// tlsKeyFile tls 私钥路径
	tlsKeyFile string

	// context 对象池
	ctxPool *sync.Pool
}

// NewServer 新建一个 http 服务器
func NewServer(app App) Server {
	s := &server{app: app, ctxPool: &sync.Pool{}}
	s.httpServer = graceful.NewServer(s, app.Logger())
	s.ctxPool.New = func() interface{} {
		return NewContext(app)
	}

	return s
}

// ServeHTTP 实现 http.Handler 接口
func (s *server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := s.Context()

	defer func() {
		if p := recover(); p != nil {
			s.app.Logger().Errorf("%+v", p)
		}

		go ctx.RunEnd()
	}()

	ctx.Reset(res, req)

	// 执行应用级别中间件
	s.app.ExecuteMiddlewares(ctx)
	if ctx.IsStopped() {
		return
	}

	// 匹配路由
	method := ctx.Method()
	path := ctx.Request().URL.Path
	handlers, dynamic := s.app.Router().Lookup(method, path)
	if handlers == nil {
		ctx.NotFound()
		return
	}

	if dynamic != nil {
		ctx.SetDynamics(dynamic)
	}

	// 执行路由处理函数和路由级别中间件
	for _, handler := range handlers {
		handler(ctx)
		if ctx.IsStopped() {
			return
		}
	}

	ctx.RunAfter()
}

// App 返回应用实例
func (s *server) App() App {
	return s.app
}

// SetTLS 指定 tls 证书，密钥路径
// certFile: 证书路径
// keyFile: 私钥路径
func (s *server) SetTLS(certFile, keyFile string) bool {
	if !file.IsExist(s.tlsCertFile) {
		s.app.Logger().Errorf("Cert file: \"%s\" is not exist", s.tlsCertFile)
		return false
	}

	if !file.IsExist(s.tlsKeyFile) {
		s.app.Logger().Errorf("Key file: \"%s\" is not exist", s.tlsKeyFile)
		return false
	}

	s.tlsCertFile = certFile
	s.tlsKeyFile = keyFile

	return true
}

// IsTLS 是否启用 tls
func (s *server) IsTLS() bool {
	return s.tlsCertFile != "" && s.tlsKeyFile != ""
}

// Start 根据配置调用 ListenAndServe 或者 ListenAndServeTLS，接收连接请求
// addr: host:port，例如: ":8080"，"192.168.1.8:80"
func (s *server) Start(addr string) error {
	logger := s.app.Logger()
	logger.Infof("Framework version: %s", s.app.Version())
	logger.Infof("PID: %d", os.Getpid())

	// tls
	if s.IsTLS() {
		if logger.IsDebugAble() {
			logger.Debugf("TLS on, %s/%s", s.tlsCertFile, s.tlsKeyFile)
			logger.Debugf("Listen on: https://%s", addr)
		}

		return s.httpServer.ListenAndServeTLS(addr, s.tlsCertFile, s.tlsKeyFile)
	}

	if logger.IsDebugAble() {
		logger.Debugf("Listen on: http://%s", addr)
	}
	return s.httpServer.ListenAndServe(addr)
}

// Context 从 pool 中获取一个 Context
func (s *server) Context() Context {
	return s.ctxPool.Get().(Context)
}

// ReleaseContext 将 Context 释放到 pool 中
func (s *server) ReleaseContext(ctx Context) {
	s.ctxPool.Put(ctx)
}

// HTTPServer 实际使用的 http 服务器
func (s *server) HTTPServer() graceful.Server {
	return s.httpServer
}
