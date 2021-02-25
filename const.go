package zeroapi

import (
	"net/http"
)

const (
	// MethodGet 请求指定的页面信息，并返回实体主体
	MethodGet = http.MethodGet

	// MethodPost 向指定资源提交数据处理请求，数据包含在请求体中
	MethodPost = http.MethodPost

	// MethodPut 从客户端向服务器传送数据，取代指定文档的内容。替换
	MethodPut = http.MethodPut

	// MethodDelete 请求服务器删除指定数据
	MethodDelete = http.MethodDelete

	// MethodConn 用于代理
	MethodConn = http.MethodConnect

	// MethodHead 类似于 GET 请求，但响应没有具体的内容，用于获取报头
	MethodHead = http.MethodHead

	// MethodPatch 类似于 PUT，但可能只包含部分数据。修改部分数据
	MethodPatch = http.MethodPatch

	// MethodOptions ..
	MethodOptions = http.MethodOptions

	// MethodTrace 回显服务器收到的请求
	MethodTrace = http.MethodTrace

	// MethodAny 同时注册 GET, POST, PUT, DELETE, HEAD, PATCH, OPTIONS
	MethodAny = "ANY"
)

// AllMethods 所有 HTTP Method
func AllMethods() []string {
	return []string{
		MethodGet,
		MethodPost,
		MethodPut,
		MethodDelete,
		MethodConn,
		MethodHead,
		MethodPatch,
		MethodOptions,
		MethodTrace,
	}
}

type (
	// Handler 处理函数
	Handler func(ctx Context)

	// HookHandler 钩子处理函数，用于中间件开发，响应 ctx.afters, ctx.ends
	HookHandler func() error

	// RouterValidator 验证函数
	RouterValidator func(s string) bool

	// CookieEncodeHandler cookie 编码与解码函数
	CookieEncodeHandler func(s string) string

	// CookieDecodeHandler cookie 编码与解码函数
	CookieDecodeHandler func(s string) (string, error)
)
