package zeroapi

import (
	"net"
	"net/http"
	"strings"
)

const (
	// ContextStatusNormal 正常状态
	ContextStatusNormal = 1
	// ContextStatusStopped 终止状态
	ContextStatusStopped = 2
)

// Context 上下文
type Context interface {
	Base
	Header
	Query
	Get
	Post
	Dynamic
	File
	Write
	Cookie
	Hook
}

// Base 基础
type Base interface {
	// App 返回应用实例
	App() App

	// Reset 重置
	Reset(res http.ResponseWriter, req *http.Request)

	// Request 获取原始 http 请求
	Request() *http.Request

	// Response 获取 http 响应
	Response() Writer

	// Method 获取当前响应的 HTTP Method，比如 GET, POST ...
	Method() string

	// Path 请求路径
	Path() string

	// HTTPCode 获取 http 状态码，默认 200
	HTTPCode() int

	// HTTPCode 设置 http 状态码
	SetHTTPCode(code int)

	// IP 获取真实IP
	IP() string

	// IPs 获取 IP 数组，每经过一级代理(匿名代理除外)，代理服务器都会把这次请求的来源IP放到数组中
	IPs() []string

	// Protocol 获取协议版本，HTTP/1.1 or HTTP/2
	Protocol() string

	// Host ..
	Host() string

	// IsAjax 判断是否是 ajax 请求
	IsAjax() bool

	// Referer HTTP Referer
	Referer() string

	// UserAgent ..
	UserAgent() string

	// NotFound 路由未找到，设置 404
	NotFound()

	// IsStopped 判断是否处于停止状态
	// 比如 auth中间件判断未通过验证，就会调用 Stopped() 来停止继续向下调用
	IsStopped() bool

	// Stopped 设置停止状态
	Stopped()

	// Value 获取对应的自定义值
	Value(key string) interface{}

	// SetValue 设置对应的自定义值
	SetValue(key string, value interface{})
}

type context struct {
	// app 应用实例
	app App
	// status 上下文状态，默认 ContextStatusNormal
	status int

	// req http 请求
	req *http.Request
	// res http 响应
	res Writer
	// httpCode
	httpCode int
	// responseSize 响应内容大小
	responseSize int64

	// dynamic 存储动态参数的值
	// 示例:
	// 定义路由: /blog/:id
	// 调用路由: /blog/1001
	// dynamics = {id: "1001"}
	dynamics map[string]string

	// values 玩家自定义数据
	values map[string]interface{}

	// afters 存储钩子函数，路由执行成功后才会执行
	afters []HookHandler
	// ends 存储钩子函数，无论路由是否执行成功，无论是否发生异常，都会在最终处执行 ends，后进先出
	ends []HookHandler

	// handlers 存储路由处理函数和中间件
	handlers []Handler
}

// NewContext 创建一个 Context 实例
func NewContext(app App) Context {
	return &context{
		app: app,
	}
}

// App 返回应用实例
func (ctx *context) App() App {
	return ctx.app
}

func (ctx *context) Reset(res http.ResponseWriter, req *http.Request) {
	ctx.res = acquireWriter()
	ctx.res.SetWriter(res)
	ctx.req = req
	ctx.status = ContextStatusNormal
	ctx.httpCode = http.StatusOK

	ctx.afters = nil
	ctx.ends = nil
}

func (ctx *context) Request() *http.Request {
	return ctx.req
}

func (ctx *context) Response() Writer {
	return ctx.res
}

func (ctx *context) Method() string {
	return ctx.req.Method
}

func (ctx *context) Path() string {
	return ctx.req.RequestURI
}

func (ctx *context) HTTPCode() int {
	return ctx.httpCode
}

func (ctx *context) SetHTTPCode(httpCode int) {
	ctx.httpCode = httpCode
	ctx.res.Writer().WriteHeader(httpCode)
}

func (ctx *context) IP() string {
	if ctx.req == nil {
		panic("Please initialize before use context")
	}

	addr := ctx.req.Header.Get("X-Real-IP")
	if addr != "" {
		return addr
	}

	ips := ctx.IPs()
	if len(ips) > 0 && ips[0] != "" {
		ip, _, err := net.SplitHostPort(ips[0])
		if err != nil {
			ip = ips[0]
		}
		return ip
	}

	addr = ctx.req.RemoteAddr
	if ip, _, err := net.SplitHostPort(addr); err == nil {
		return ip
	}

	return ctx.req.RemoteAddr
}

func (ctx *context) IPs() []string {
	if ips := ctx.Header("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}

	return []string{}
}

func (ctx *context) Protocol() string {
	return ctx.req.Proto
}

func (ctx *context) Host() string {
	if ctx.req.Host != "" {
		if host, _, err := net.SplitHostPort(ctx.req.Host); err == nil {
			return host
		}
		return ctx.req.Host
	}

	return "localhost"
}

func (ctx *context) IsAjax() bool {
	return ctx.Header("X-Requested-With") == "XMLHttpRequest"
}

func (ctx *context) Referer() string {
	return ctx.Header("Referer")
}

func (ctx *context) UserAgent() string {
	return ctx.Header("User-Agent")
}

func (ctx *context) NotFound() {
	ctx.SetHTTPCode(http.StatusNotFound)
	ctx.Message(http.StatusNotFound, "PAGE NOT FOUND")
}

func (ctx *context) IsStopped() bool {
	return ctx.status == ContextStatusStopped
}

func (ctx *context) Stopped() {
	ctx.status = ContextStatusStopped
}

// Value 获取对应的自定义值
func (ctx *context) Value(key string) interface{} {
	if ctx.values == nil {
		return nil
	}

	return ctx.values[key]
}

// SetValue 设置对应的自定义值
func (ctx *context) SetValue(key string, value interface{}) {
	if ctx.values == nil {
		ctx.values = make(map[string]interface{})
	}

	ctx.values[key] = value
}
