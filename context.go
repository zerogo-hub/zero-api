package zeroapi

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zerogo-hub/zero-helper/bytes"
	"github.com/zerogo-hub/zero-helper/file"
	"google.golang.org/protobuf/proto"
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

	// Response 获取原始 http 响应
	Response() http.ResponseWriter

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
}

// Header ..
type Header interface {
	// Header 获取 http.Request 中 Header 指定参数的值
	Header(key string) string

	// AddHeader 添加响应头
	// header 该响应头有多个值 key = []string{value1, value2, ...}
	AddHeader(key, value string)

	// SetHeader 设置响应头
	// header 该响应头只有一个值 key = []string{value}
	SetHeader(key, value string)

	// DelHeader 移除响应中的 header
	DelHeader(key string)
}

// Query 包括 GET, POST, PUT
type Query interface {
	// Query 获取指定参数的值
	// 多个参数同名，只获取第一个
	// 例如 /user?id=Yaha&id=Gama
	// Query("id") 的结果为 "Yaha"
	Query(key string) string

	// QueryAll 获取所有参数的值
	// 例如 /user?id=Yaha&id=Gama&age=18
	// QueryAll() 的结果为 {"id": ["Yaha", "Gama"], "age": ["18"]}
	QueryAll() map[string][]string

	// QueryStrings 获取指定参数的值
	// 例如 /user?id=Yaha&id=Gama
	// QueryStrings("id") 的结果为 ["Yaha", "Gama"]
	QueryStrings(key string) []string

	// QueryEscape 获取指定参数的值，并对被转码的结果进行还原
	QueryEscape(key string) string

	// QueryBool 获取指定参数的值，并将结果转为 bool
	QueryBool(key string) bool

	// QueryInt32 获取指定参数的值，并将结果转为 int32
	QueryInt32(key string) int32

	// QueryInt64 获取指定参数的值，并将结果转为 int64
	QueryInt64(key string) int64

	// QueryFloat32 获取指定参数的值，并将结果转为 float32
	QueryFloat32(key string) float32

	// QueryFloat64 获取指定参数的值，并将结果转为 float64
	QueryFloat64(key string) float64

	// QueryDefault 获取指定参数的值，如果不存在，则返回默认值 def
	QueryDefault(key, def string) string

	// QueryBoolDefault 获取指定参数的值(结果转为 bool)，如果不存在，则返回默认值 def
	QueryBoolDefault(key string, def bool) bool

	// QueryInt32Default 获取指定参数的值(结果转为 int32)，如果不存在，则返回默认值 def
	QueryInt32Default(key string, def int32) int32

	// QueryInt64Default 获取指定参数的值(结果转为 int64)，如果不存在，则返回默认值 def
	QueryInt64Default(key string, def int64) int64

	// QueryFloat32Default 获取指定参数的值(结果转为 float32)，如果不存在，则返回默认值 def
	QueryFloat32Default(key string, def float32) float32

	// QueryFloat64Default 获取指定参数的值(结果转为 float64)，如果不存在，则返回默认值 def
	QueryFloat64Default(key string, def float64) float64
}

// Get 只包括 GET
type Get interface {
	// Get 获取指定参数的值
	// 多个参数同名，只获取第一个
	// 例如 /user?id=Yaha&id=Gama
	// Get("id") 的结果为 "Yaha"
	Get(key string) string

	// GetEscape 获取指定参数的值，并对被编码的结果进行解码
	GetEscape(key string) string

	// GetBool 获取指定参数的值，并将结果转为 bool
	GetBool(key string) bool

	// GetInt32 获取指定参数的值，并将结果转为 int32
	GetInt32(key string) int32

	// GetInt64 获取指定参数的值，并将结果转为 int64
	GetInt64(key string) int64

	// GetFloat32 获取指定参数的值，并将结果转为 float32
	GetFloat32(key string) float32

	// GetFloat64 获取指定参数的值，并将结果转为 float64
	GetFloat64(key string) float64

	// GetDefault 获取指定参数的值，如果不存在，则返回默认值 def
	GetDefault(key, def string) string

	// GetBoolDefault 获取指定参数的值(结果转为 bool)，如果不存在，则返回默认值 def
	GetBoolDefault(key string, def bool) bool

	// GetInt32Default 获取指定参数的值(结果转为 int32)，如果不存在，则返回默认值 def
	GetInt32Default(key string, def int32) int32

	// GetInt64Default 获取指定参数的值(结果转为 int64)，如果不存在，则返回默认值 def
	GetInt64Default(key string, def int64) int64

	// GetFloat32Default 获取指定参数的值(结果转为 float32)，如果不存在，则返回默认值 def
	GetFloat32Default(key string, def float32) float32

	// GetFloat64Default 获取指定参数的值(结果转为 float64)，如果不存在，则返回默认值 def
	GetFloat64Default(key string, def float64) float64
}

// Post 包括 POST, PUT, PATCH
type Post interface {
	// Post 获取指定参数的值
	// 多个参数同名，只获取第一个
	// 例如 /user?id=Yaha&id=Gama
	// Post("id") 的结果为 "Yaha"
	Post(key string) string

	// PostStrings 获取指定参数的值
	// 例如 /user?id=Yaha&id=Gama
	// PostStrings("id") 的结果为 ["Yaha", "Gama"]
	PostStrings(key string) []string

	// PostEscape 获取指定参数的值，并对被编码的结果进行解码
	PostEscape(key string) string

	// PostBool 获取指定参数的值，并将结果转为 bool
	PostBool(key string) bool

	// PostInt32 获取指定参数的值，并将结果转为 int32
	PostInt32(key string) int32

	// PostInt64 获取指定参数的值，并将结果转为 int64
	PostInt64(key string) int64

	// PostFloat32 获取指定参数的值，并将结果转为 float32
	PostFloat32(key string) float32

	// PostFloat64 获取指定参数的值，并将结果转为 float64
	PostFloat64(key string) float64

	// PostDefault 获取指定参数的值，如果不存在，则返回默认值 def
	PostDefault(key, def string) string

	// PostBoolDefault 获取指定参数的值(结果转为 bool)，如果不存在，则返回默认值 def
	PostBoolDefault(key string, def bool) bool

	// PostInt32Default 获取指定参数的值(结果转为 int32)，如果不存在，则返回默认值 def
	PostInt32Default(key string, def int32) int32

	// PostInt64Default 获取指定参数的值(结果转为 int64)，如果不存在，则返回默认值 def
	PostInt64Default(key string, def int64) int64

	// PostFloat32Default 获取指定参数的值(结果转为 float32)，如果不存在，则返回默认值 def
	PostFloat32Default(key string, def float32) float32

	// PostFloat64Default 获取指定参数的值(结果转为 float64)，如果不存在，则返回默认值 def
	PostFloat64Default(key string, def float64) float64
}

// Dynamic 动态参数
type Dynamic interface {
	// Dynamic 获取动态参数的值
	// 示例:
	// 定义路由: /blog/:id
	// 调用路由: /blog/1001
	// params = {id: "1001"}
	// Dynamic("id") -> 1001
	Dynamic(key string) string

	// SetDynamic 设置动态参数，key 的格式为 "param" 或者 ":param"
	SetDynamic(key string, value string) error

	// SetDynamics 替换动态参数
	SetDynamics(dynamics map[string]string)
}

// File 文件相关
type File interface {
	// File 获取上传文件信息
	File(key string) (multipart.File, *multipart.FileHeader, error)

	// Files 获取上传文件信息，可能有多个文件
	// cbs 在存盘前修改 multipart.FileHeader
	Files(destDirectory string, cbs ...func(Context, *multipart.FileHeader)) (int64, error)

	// DownloadFile 下载文件
	// path 文件路径
	// filename 文件名称
	DownloadFile(path string, filename ...string)
}

// Write 响应相关
type Write interface {
	// Bytes 将数据写入响应
	Bytes(bytes []byte) (int, error)

	// Text 将数据写入响应中
	Text(value string) (int, error)

	// Textf 将数据写入响应中
	Textf(format string, a ...interface{})

	// Map map 转 text
	Map(obj interface{}) (int, error)

	// JSON 将数据转为 JSON 格式写入响应
	JSON(obj interface{}) (int, error)

	// XML 将数据转为 XML 格式写入响应
	XML(obj interface{}) (int, error)

	// HTML 发送 html 响应
	HTML(html string) (int, error)

	// HTMLf 发送 html 响应
	HTMLf(format string, a ...interface{}) (int, error)

	// Protobuf 将数据装为 google protobuf 格式，写入响应
	Protobuf(obj interface{}) (int, error)

	// Size 响应的数据大小
	Size() int64

	// Redirect 重定向
	// httpCode: HTTP Code， 需要在 3xx 范围内, 比如 301, 302, 303 ... 308
	// url: 重定向后的地址
	Redirect(httpCode int, url string) error

	// Flush 将数据推向客户端
	Flush()

	// Push HTTP/2 服务器推送
	Push(value string, opts *http.PushOptions) error

	// AutoContentType 根据给定的文件类型，自动设置 Content-Type
	// .json -> app/json
	// fileExt: 文件后缀名，例如 .json
	AutoContentType(fileExt string)

	// Message 传递 {"code": xx, "message": xxx}
	Message(code int, message ...string) (int, error)
}

// CookieOption cookie 选项
type CookieOption func(cookie *http.Cookie)

// Cookie cookie 相关
type Cookie interface {

	// Cookie 获取 cookie 值
	Cookie(key string, opts ...CookieOption) (string, error)

	// SetCookie 设置 cookie，见 https://tools.ietf.org/html/rfc6265
	// key: cookie 参数名称
	// value: cookie 值
	// maxAge: 见 https://tools.ietf.org/html/rfc6265#section-4.1.2.2
	// 		 = 0: 表示不指定
	// 		 < 0: 表示立即删除
	// 		 > 0: cookie 生存时间，单位秒
	// domain: 见 https://tools.ietf.org/html/rfc6265#section-4.1.2.3
	// path: 见 https://tools.ietf.org/html/rfc6265#section-4.1.2.4
	// secure: 见 https://tools.ietf.org/html/rfc6265#section-4.1.2.5
	// httpOnly: 见 https://tools.ietf.org/html/rfc6265#section-4.1.2.6
	SetCookie(key, value string, opts ...CookieOption)

	// RemoveCookie 移除指定的 cookie
	RemoveCookie(key string, opts ...CookieOption)

	// SetHTTPCookie 设置原始的 cookie
	SetHTTPCookie(cookie *http.Cookie)

	// HTTPCookies 获取所有原始的 cookie
	HTTPCookies() []*http.Cookie
}

// Hook 钩子，一般用于中间件中
type Hook interface {
	// AppendAfter 添加处理函数，这些函数 在中间件和路由函数执行完毕之后执行，如果中间有中断，则不会执行
	// 先入后出顺序执行，越先加入的函数越后执行
	AppendAfter(hook HookHandler)

	// AppendEnd 添加处理函数，这些函数 在处理都完成之后执行，无论中间是否有中断和异常，都会执行
	// 先入后出顺序执行，越先加入的函数越后执行
	AppendEnd(hook HookHandler)

	// RunAfter 执行通过 AppendAfter 加入的处理函数
	RunAfter()

	// RunEnd 执行通过 AppendEnd 加入的处理函数
	RunEnd()
}

type context struct {
	// app 应用实例
	app App
	// status 上下文状态，默认 ContextStatusNormal
	status int

	// req http 请求
	req *http.Request
	// res http 响应
	res http.ResponseWriter
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

func (ctx *context) Reset(res http.ResponseWriter, req *http.Request) {
	ctx.res = res
	ctx.req = req
	ctx.status = ContextStatusNormal
	ctx.httpCode = http.StatusOK

	ctx.afters = nil
	ctx.ends = nil
}

func (ctx *context) Request() *http.Request {
	return ctx.req
}

func (ctx *context) Response() http.ResponseWriter {
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
	ctx.res.WriteHeader(httpCode)
}

func (ctx *context) App() App {
	return ctx.app
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

func (ctx *context) Header(key string) string {
	return ctx.req.Header.Get(key)
}

func (ctx *context) AddHeader(key, value string) {
	if key != "" && value != "" {
		ctx.res.Header().Add(key, value)
	}
}

func (ctx *context) SetHeader(key, value string) {
	if key != "" && value != "" {
		ctx.res.Header().Set(key, value)
	}
}

func (ctx *context) DelHeader(key string) {
	if key != "" {
		ctx.res.Header().Del(key)
	}
}

// queryAll 获取所有的参数值，内部使用
func (ctx *context) queryAll() (map[string][]string, bool) {
	if err := ctx.req.ParseForm(); err != nil {
		return nil, false
	}
	// Form: 需要先调用 ParseForm()
	// including both the URL field's query parameters and the POST or PUT form data
	if form := ctx.req.Form; len(form) > 0 {
		return form, true
	}
	// Form 中没有数据，则使用了 PATCH ?
	//
	// PostForm: 需要先调用 ParseForm()
	// contains the parsed form data from POST, PATCH, or PUT body parameters
	if form := ctx.req.PostForm; len(form) > 0 {
		return form, true
	}

	return nil, false
}

func (ctx *context) Query(key string) string {
	if form, exist := ctx.queryAll(); exist {
		if values := form[key]; len(values) > 0 {
			return values[0]
		}
	}

	return ""
}

func (ctx *context) QueryAll() map[string][]string {
	form, _ := ctx.queryAll()
	return form
}

func (ctx *context) QueryStrings(key string) []string {
	if form, exist := ctx.queryAll(); exist {
		if values := form[key]; len(values) > 0 {
			return values
		}
	}

	return nil
}

// Unescape 解码
// 解码以下编码的结果:
// js: escape, encodeURI, encodeURIComponent
// go: QueryEscape
func unescape(value string) string {
	if value == "" {
		return ""
	}
	result, err := url.QueryUnescape(value)
	if err != nil {
		return value
	}
	return result
}

func (ctx *context) QueryEscape(key string) string {
	return unescape(ctx.Query(key))
}

func (ctx *context) QueryBool(key string) bool {
	v, _ := strconv.ParseBool(ctx.Query(key))
	return v
}

func (ctx *context) QueryInt32(key string) int32 {
	v, _ := strconv.ParseInt(ctx.Query(key), 10, 32)
	return int32(v)
}

func (ctx *context) QueryInt64(key string) int64 {
	v, _ := strconv.ParseInt(ctx.Query(key), 10, 64)
	return v
}

func (ctx *context) QueryFloat32(key string) float32 {
	v, _ := strconv.ParseFloat(ctx.Query(key), 32)
	return float32(v)
}

func (ctx *context) QueryFloat64(key string) float64 {
	v, _ := strconv.ParseFloat(ctx.Query(key), 64)
	return v
}

func (ctx *context) QueryDefault(key, def string) string {
	if value := ctx.Query(key); value != "" {
		return value
	}

	return def
}

func (ctx *context) QueryBoolDefault(key string, def bool) bool {
	if value := ctx.Query(key); value != "" {
		result, err := strconv.ParseBool(value)
		if err != nil {
			return def
		}
		return result
	}

	return def
}

func (ctx *context) QueryInt32Default(key string, def int32) int32 {
	if value := ctx.Query(key); value != "" {
		result, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return def
		}
		return int32(result)
	}

	return def
}

func (ctx *context) QueryInt64Default(key string, def int64) int64 {
	if value := ctx.Query(key); value != "" {
		result, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return def
		}
		return result
	}

	return def
}

func (ctx *context) QueryFloat32Default(key string, def float32) float32 {
	if value := ctx.Query(key); value != "" {
		result, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return def
		}
		return float32(result)
	}

	return def
}

func (ctx *context) QueryFloat64Default(key string, def float64) float64 {
	if value := ctx.Query(key); value != "" {
		result, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return def
		}
		return result
	}

	return def
}

func (ctx *context) Get(key string) string {
	return ctx.req.URL.Query().Get(key)
}

func (ctx *context) GetEscape(key string) string {
	return unescape(ctx.Get(key))
}

func (ctx *context) GetBool(key string) bool {
	v, _ := strconv.ParseBool(ctx.Get(key))
	return v
}

func (ctx *context) GetInt32(key string) int32 {
	v, _ := strconv.ParseInt(ctx.Get(key), 10, 32)
	return int32(v)
}

func (ctx *context) GetInt64(key string) int64 {
	v, _ := strconv.ParseInt(ctx.Get(key), 10, 64)
	return v
}

func (ctx *context) GetFloat32(key string) float32 {
	v, _ := strconv.ParseFloat(ctx.Get(key), 32)
	return float32(v)
}

func (ctx *context) GetFloat64(key string) float64 {
	v, _ := strconv.ParseFloat(ctx.Get(key), 64)
	return v
}

func (ctx *context) GetDefault(key, def string) string {
	if value := ctx.Get(key); value != "" {
		return value
	}

	return def
}

func (ctx *context) GetBoolDefault(key string, def bool) bool {
	if value := ctx.Get(key); value != "" {
		result, err := strconv.ParseBool(value)
		if err != nil {
			return def
		}
		return result
	}

	return def
}

func (ctx *context) GetInt32Default(key string, def int32) int32 {
	if value := ctx.Get(key); value != "" {
		result, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return def
		}
		return int32(result)
	}

	return def
}

func (ctx *context) GetInt64Default(key string, def int64) int64 {
	if value := ctx.Get(key); value != "" {
		result, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return def
		}
		return result
	}

	return def
}

func (ctx *context) GetFloat32Default(key string, def float32) float32 {
	if value := ctx.Get(key); value != "" {
		result, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return def
		}
		return float32(result)
	}

	return def
}

func (ctx *context) GetFloat64Default(key string, def float64) float64 {
	if value := ctx.Get(key); value != "" {
		result, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return def
		}
		return result
	}

	return def
}

func (ctx *context) Post(key string) string {
	// PostForm: 需要先调用 ParseForm()
	// contains the parsed form data from POST, PATCH, or PUT body parameters

	if err := ctx.req.ParseForm(); err != nil {
		return ""
	}

	return ctx.req.PostForm.Get(key)
}

func (ctx *context) PostStrings(key string) []string {
	// PostForm: 需要先调用 ParseForm()
	// contains the parsed form data from POST, PATCH, or PUT body parameters

	if err := ctx.req.ParseForm(); err != nil {
		return nil
	}

	if values, exist := ctx.req.PostForm[key]; exist {
		return values
	}

	return nil
}

func (ctx *context) PostEscape(key string) string {
	return unescape(ctx.Post(key))
}

func (ctx *context) PostBool(key string) bool {
	v, _ := strconv.ParseBool(ctx.Post(key))
	return v
}

func (ctx *context) PostInt32(key string) int32 {
	v, _ := strconv.ParseInt(ctx.Post(key), 10, 32)
	return int32(v)
}

func (ctx *context) PostInt64(key string) int64 {
	v, _ := strconv.ParseInt(ctx.Post(key), 10, 64)
	return v
}

func (ctx *context) PostFloat32(key string) float32 {
	v, _ := strconv.ParseFloat(ctx.Post(key), 32)
	return float32(v)
}

func (ctx *context) PostFloat64(key string) float64 {
	v, _ := strconv.ParseFloat(ctx.Post(key), 64)
	return v
}

func (ctx *context) PostDefault(key, def string) string {
	if value := ctx.Post(key); value != "" {
		return value
	}

	return def
}

func (ctx *context) PostBoolDefault(key string, def bool) bool {
	if value := ctx.Post(key); value != "" {
		result, err := strconv.ParseBool(value)
		if err != nil {
			return def
		}
		return result
	}

	return def
}

func (ctx *context) PostInt32Default(key string, def int32) int32 {
	if value := ctx.Post(key); value != "" {
		result, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return def
		}
		return int32(result)
	}

	return def
}

func (ctx *context) PostInt64Default(key string, def int64) int64 {
	if value := ctx.Post(key); value != "" {
		result, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return def
		}
		return result
	}

	return def
}

func (ctx *context) PostFloat32Default(key string, def float32) float32 {
	if value := ctx.Post(key); value != "" {
		result, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return def
		}
		return float32(result)
	}

	return def
}

func (ctx *context) PostFloat64Default(key string, def float64) float64 {
	if value := ctx.Post(key); value != "" {
		result, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return def
		}
		return result
	}

	return def
}

func (ctx *context) Dynamic(key string) string {
	if len(key) == 0 {
		return ""
	}

	// 兼容判断
	if key[0] == ':' {
		key = key[1:]
	}

	if ctx.dynamics != nil {
		if value, exist := ctx.dynamics[key]; exist {
			return value
		}
	}

	return ""
}

func (ctx *context) SetDynamic(key string, value string) error {
	if len(key) == 0 {
		return errors.New("Parameter key cannot be empty")
	}

	// 兼容判断
	if key[0] == ':' {
		key = key[1:]
	}

	if ctx.dynamics == nil {
		ctx.dynamics = make(map[string]string)
	}

	ctx.dynamics[key] = value

	return nil
}

func (ctx *context) SetDynamics(dynamics map[string]string) {
	ctx.dynamics = dynamics
}

// upload 从临时文件夹或者内存中写入到指定位置的文件夹中
func upload(dest string, header *multipart.FileHeader) (int64, error) {
	// 打开临时文件或者内存中的文件内容
	src, err := header.Open()
	if err != nil {
		return 0, err
	}
	defer src.Close()

	file, err := os.OpenFile(filepath.Join(dest, header.Filename), os.O_WRONLY|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return io.Copy(file, src)
}

func (ctx *context) File(key string) (multipart.File, *multipart.FileHeader, error) {
	if err := ctx.req.ParseMultipartForm(ctx.app.FileMaxMemory()); err != nil {
		return nil, nil, err
	}

	return ctx.req.FormFile(key)
}

func (ctx *context) Files(destDirectory string, cbs ...func(Context, *multipart.FileHeader)) (int64, error) {
	if err := ctx.req.ParseMultipartForm(ctx.app.FileMaxMemory()); err != nil {
		return 0, err
	}

	// MultipartForm: 需要先调用 ParseMultipartForm，
	// including file uploads

	if ctx.req.MultipartForm != nil {
		if f := ctx.req.MultipartForm.File; f != nil {
			var l int64
			for _, files := range f {
				for _, file := range files {
					for _, cb := range cbs {
						cb(ctx, file)
					}
					length, err := upload(destDirectory, file)
					if err != nil {
						return 0, err
					}
					l += length
				}
			}
			return l, nil
		}
	}

	return 0, http.ErrMissingFile
}

func (ctx *context) DownloadFile(path string, filename ...string) {
	if !file.IsExist(path) {
		http.ServeFile(ctx.res, ctx.req, path)
		return
	}

	fname := ""
	if len(filename) > 0 && filename[0] != "" {
		fname = filename[0]
	} else {
		fname = file.NameRand(path, 8)
	}

	ctx.AddHeader("Content-Disposition", "attachment; filename="+fname)
	ctx.AddHeader("Content-Description", "File Transfer")
	ctx.AddHeader("Content-Type", "app/octet-stream")
	ctx.AddHeader("Content-Transfer-Encoding", "binary")
	ctx.AddHeader("Expires", "0")
	ctx.AddHeader("Cache-Control", "must-revalidate")
	ctx.AddHeader("Pragma", "public")
	http.ServeFile(ctx.res, ctx.req, path)
}

func (ctx *context) Bytes(bytes []byte) (int, error) {
	var size int
	var err error

	size, err = ctx.res.Write(bytes)

	if err != nil {
		return 0, err
	}

	ctx.responseSize += int64(size)

	return size, nil
}

func (ctx *context) Text(value string) (int, error) {
	var size int
	var err error

	size, err = ctx.res.Write(bytes.StringToBytes(value))

	if err != nil {
		return 0, err
	}

	ctx.responseSize += int64(size)
	ctx.SetHeader("Content-Type", "text/plain;charset=utf-8")

	return size, nil
}

func (ctx *context) Textf(format string, a ...interface{}) {
	ctx.Text(fmt.Sprintf(format, a...))
}

func (ctx *context) Map(obj interface{}) (int, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return 0, err
	}

	return ctx.Bytes(bytes)
}

func (ctx *context) JSON(obj interface{}) (int, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return 0, err
	}

	ctx.SetHeader("Content-Type", "application/json;charset=utf-8")

	return ctx.Bytes(bytes)
}

func (ctx *context) XML(obj interface{}) (int, error) {
	bytes, err := xml.Marshal(obj)
	if err != nil {
		return 0, err
	}

	ctx.SetHeader("Content-Type", "application/xml;charset=utf-8")

	return ctx.Bytes(bytes)
}

func (ctx *context) HTML(html string) (int, error) {
	ctx.SetHeader("Content-Type", "text/html;charset=utf-8")

	return ctx.Bytes(bytes.StringToBytes(template.HTMLEscapeString(html)))
}

func (ctx *context) HTMLf(format string, a ...interface{}) (int, error) {
	return ctx.HTML(fmt.Sprintf(format, a...))
}

func (ctx *context) Protobuf(obj interface{}) (int, error) {
	bytes, err := proto.Marshal(obj.(proto.Message))
	if err != nil {
		return 0, err
	}

	ctx.SetHeader("Content-Type", "application/x-protobuf;charset=utf-8")

	return ctx.Bytes(bytes)
}

func (ctx *context) Size() int64 {
	return ctx.responseSize
}

func (ctx *context) Redirect(httpCode int, url string) error {
	if httpCode < http.StatusMultipleChoices || httpCode > http.StatusPermanentRedirect {
		return errors.New("httpCode should be in the 3xx, like 301, 302 etc")
	}

	ctx.Stopped()
	http.Redirect(ctx.res, ctx.req, url, httpCode)

	return nil
}

func (ctx *context) Flush() {
	if flusher, ok := ctx.res.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (ctx *context) Push(value string, opts *http.PushOptions) error {
	if push, ok := ctx.res.(http.Pusher); ok {
		push.Push(value, opts)
	}

	return http.ErrNotSupported
}

func (ctx *context) AutoContentType(fileExt string) {
	if !strings.HasPrefix(fileExt, ".") {
		fileExt = "." + fileExt
	}

	if ct := mime.TypeByExtension(fileExt); ct != "" {
		ctx.AddHeader("Content-Type", ct)
	}
}

func (ctx *context) Message(code int, message ...string) (int, error) {
	result := make(map[string]string)

	result["code"] = strconv.Itoa(code)
	if len(message) > 0 {
		result["message"] = message[0]
	}

	return ctx.Map(result)
}

func (ctx *context) AppendAfter(hook HookHandler) {
	ctx.afters = append(ctx.afters, hook)
}

func (ctx *context) AppendEnd(hook HookHandler) {
	ctx.ends = append(ctx.ends, hook)
}

func (ctx *context) RunAfter() {
	run(ctx.afters)
}

func (ctx *context) RunEnd() {
	// TODO 回收 ctx
	run(ctx.ends)
}

func run(hooks []HookHandler) {
	if len(hooks) == 0 {
		return
	}

	// 从尾巴开始执行
	for i := len(hooks) - 1; i >= 0; i-- {
		if err := hooks[i](); err != nil {
			return
		}
	}
}
