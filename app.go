package zeroapi

import (
	"github.com/zerogo-hub/zero-helper/logger"
)

// App 应用
type App interface {

	// Version 获取框架版本号
	Version() string

	// Logger 获取日志实例
	Logger() logger.Logger

	// Router 获取路由管理示例
	Router() Router

	// FileMaxMemory 文件系统使用的最大内存
	FileMaxMemory() int64

	// IsCookieEncode cookie 是否需要进行编码
	IsCookieEncode() bool

	// CookieEncodeHandler 获取 cookie 编码函数
	CookieEncodeHandler() CookieEncodeHandler

	// CookieDecodeHandler 获取 cookie 解码函数
	CookieDecodeHandler() CookieDecodeHandler

	// SetCookieSignHandler 设置 cookie 编码与解码函数
	SetCookieHandler(encoder CookieEncodeHandler, decoder CookieDecodeHandler)
}

type app struct {
	// router 路由管理器
	router Router
}

// NewApp 生成一个应用实例
func NewApp() App {
	a := new(app)
	r := NewRouter(a)

	a.router = r

	return a
}

// Version 获取框架版本号
func (a *app) Version() string {
	return ""
}

// Logger 获取日志实例
func (a *app) Logger() logger.Logger {
	return nil
}

// Router 获取路由管理示例
func (a *app) Router() Router {
	return a.router
}

// FileMaxMemory 文件系统使用的最大内存
func (a *app) FileMaxMemory() int64 {
	return 0
}

// IsCookieEncode cookie 是否需要进行编码
func (a *app) IsCookieEncode() bool {
	return false
}

// CookieEncodeHandler 获取 cookie 编码函数
func (a *app) CookieEncodeHandler() CookieEncodeHandler {
	return nil
}

// CookieDecodeHandler 获取 cookie 解码函数
func (a *app) CookieDecodeHandler() CookieDecodeHandler {
	return nil
}

// SetCookieSignHandler 设置 cookie 编码与解码函数
func (a *app) SetCookieHandler(encoder CookieEncodeHandler, decoder CookieDecodeHandler) {

}
