package server

import (
	"net/http"
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
	zeroctx "github.com/zerogo-hub/zero-api/context"

	"github.com/zerogo-hub/zero-helper/file"
	zerograceful "github.com/zerogo-hub/zero-helper/graceful/http"
)

type server struct {
	// app 应用实例
	app zeroapi.App

	// httpServer 实际使用  graceful.Server 替代 http.Server
	httpServer zerograceful.Server

	// tlsCertFile tls 证书路径
	tlsCertFile string

	// tlsKeyFile tls 私钥路径
	tlsKeyFile string
}

// NewServer 新建一个 http 服务器
func NewServer(app zeroapi.App) zeroapi.Server {
	s := &server{app: app}
	s.httpServer = zerograceful.NewServer(s, app.Logger())

	return s
}

// ServeHTTP 实现 http.Handler 接口
func (s *server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := s.app.Context()

	defer func() {
		if p := recover(); p != nil {
			s.app.Logger().Errorf("%+v", p)
		}

		go ctx.RunEnd()

		zeroctx.ReleaseWriter(ctx.Response())
	}()

	if s.app.MaxMemory() > 0 {
		req.Body = http.MaxBytesReader(res, req.Body, s.app.MaxMemory())
	}

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
		if handler == nil {
			continue
		}

		handler(ctx)
		if ctx.IsStopped() {
			return
		}
	}

	ctx.RunAfter()
}

// Start 根据配置调用 ListenAndServe 或者 ListenAndServeTLS，接收连接请求
// addr: host:port，例如: ":8080"，"192.168.1.8:80"
func (s *server) Start(addr string) error {
	logger := s.app.Logger()
	logger.Infof("Framework version: %s", s.app.Version())
	logger.Infof("PID: %d", os.Getpid())

	// tls
	if s.tlsCertFile != "" && s.tlsKeyFile != "" {
		if logger.IsInfoAble() {
			logger.Infof("TLS on, %s/%s", s.tlsCertFile, s.tlsKeyFile)
			logger.Infof("Listen on: https://%s", addr)
		}

		return s.httpServer.ListenAndServeTLS(addr, s.tlsCertFile, s.tlsKeyFile)
	}

	if logger.IsInfoAble() {
		logger.Infof("Listen on: http://%s", addr)
	}

	return s.httpServer.ListenAndServe(addr)
}

// HTTPServer 实际使用的 http 服务器
func (s *server) HTTPServer() zerograceful.Server {
	return s.httpServer
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

// SetShutdownTimeout 设置优雅退出超时时间
// 服务器会每隔500毫秒检查一次连接是否都断开处理完毕
// 如果超过超时时间，就不再检查，直接退出
//
// ms: 单位：毫秒，当 <= 0 时无效，直接退出
func (s *server) SetShutdownTimeout(ms int) {
	s.httpServer.SetShutdownTimeout(ms)
}

// RegisterShutdownHandler 注册关闭函数
// 按照注册的顺序调用这些函数
// 所有已经添加的服务器都会响应这个函数
func (s *server) RegisterShutdownHandler(f func()) {
	s.httpServer.RegisterShutdownHandler(f)
}
