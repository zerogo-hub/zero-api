package zeroapi

import (
	"mime/multipart"
	"net/http"

	zerograceful "github.com/zerogo-hub/zero-helper/graceful/http"
	zerologger "github.com/zerogo-hub/zero-helper/logger"
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
	Logger() zerologger.Logger

	// MaxMemory 使用的最大内存
	MaxMemory() int64

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

// Context 上下文
type Context interface {
	ContextBase
	ContextHeader
	ContextQuery
	ContextGet
	ContextPost
	ContextDynamic
	ContextFile
	ContextWrite
	ContextCookie
	ContextHook
}

// ContextBase 基础
type ContextBase interface {
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

// ContextHeader ..
type ContextHeader interface {
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

// ContextQuery 包括 GET, POST, PUT
type ContextQuery interface {
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

	// QueryInt8 获取指定参数的值，并将结果转为 int8
	QueryInt8(key string) int8

	// QueryUint8 获取指定参数的值，并将结果转为 uint8
	QueryUint8(key string) uint8

	// QueryInt16 获取指定参数的值，并将结果转为 int16
	QueryInt16(key string) int16

	// QueryUint16 获取指定参数的值，并将结果转为 uint16
	QueryUint16(key string) uint16

	// QueryInt32 获取指定参数的值，并将结果转为 int32
	QueryInt32(key string) int32

	// QueryUint32 获取指定参数的值，并将结果转为 uint32
	QueryUint32(key string) uint32

	// QueryInt64 获取指定参数的值，并将结果转为 int64
	QueryInt64(key string) int64

	// QueryUint64 获取指定参数的值，并将结果转为 uint64
	QueryUint64(key string) uint64

	// QueryFloat32 获取指定参数的值，并将结果转为 float32
	QueryFloat32(key string) float32

	// QueryFloat64 获取指定参数的值，并将结果转为 float64
	QueryFloat64(key string) float64

	// QueryDefault 获取指定参数的值，如果不存在，则返回默认值 def
	QueryDefault(key, def string) string

	// QueryBoolDefault 获取指定参数的值(结果转为 bool)，如果不存在，则返回默认值 def
	QueryBoolDefault(key string, def bool) bool

	// QueryInt8Default 获取指定参数的值(结果转为 int8)，如果不存在，则返回默认值 def
	QueryInt8Default(key string, def int8) int8

	// QueryUint8Default 获取指定参数的值(结果转为 uint8)，如果不存在，则返回默认值 def
	QueryUint8Default(key string, def uint8) uint8

	// QueryInt16Default 获取指定参数的值(结果转为 int16)，如果不存在，则返回默认值 def
	QueryInt16Default(key string, def int16) int16

	// QueryUint16Default 获取指定参数的值(结果转为 uint16)，如果不存在，则返回默认值 def
	QueryUint16Default(key string, def uint16) uint16

	// QueryInt32Default 获取指定参数的值(结果转为 int32)，如果不存在，则返回默认值 def
	QueryInt32Default(key string, def int32) int32

	// QueryUint32Default 获取指定参数的值(结果转为 uint32)，如果不存在，则返回默认值 def
	QueryUint32Default(key string, def uint32) uint32

	// QueryInt64Default 获取指定参数的值(结果转为 int64)，如果不存在，则返回默认值 def
	QueryInt64Default(key string, def int64) int64

	// QueryUint64Default 获取指定参数的值(结果转为 uint64)，如果不存在，则返回默认值 def
	QueryUint64Default(key string, def uint64) uint64

	// QueryFloat32Default 获取指定参数的值(结果转为 float32)，如果不存在，则返回默认值 def
	QueryFloat32Default(key string, def float32) float32

	// QueryFloat64Default 获取指定参数的值(结果转为 float64)，如果不存在，则返回默认值 def
	QueryFloat64Default(key string, def float64) float64
}

// ContextGet 只包括 GET
type ContextGet interface {
	// Get 获取指定参数的值
	// 多个参数同名，只获取第一个
	// 例如 /user?id=Yaha&id=Gama
	// Get("id") 的结果为 "Yaha"
	Get(key string) string

	// Gets 获取指定参数的所有值
	Gets(key string) []string

	// GetEscape 获取指定参数的值，并对被编码的结果进行解码
	GetEscape(key string) string

	// GetBool 获取指定参数的值，并将结果转为 bool
	GetBool(key string) bool

	// GetInt8 获取指定参数的值，并将结果转为 int8
	GetInt8(key string) int8

	// GetUint8 获取指定参数的值，并将结果转为 uint8
	GetUint8(key string) uint8

	// GetInt16 获取指定参数的值，并将结果转为 int16
	GetInt16(key string) int16

	// GetUint16 获取指定参数的值，并将结果转为 uint16
	GetUint16(key string) uint16

	// GetInt32 获取指定参数的值，并将结果转为 int32
	GetInt32(key string) int32

	// GetUint32 获取指定参数的值，并将结果转为 uint32
	GetUint32(key string) uint32

	// GetInt64 获取指定参数的值，并将结果转为 int64
	GetInt64(key string) int64

	// GetUint64 获取指定参数的值，并将结果转为 uint64
	GetUint64(key string) uint64

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

// ContextPost 包括 POST, PUT, PATCH
type ContextPost interface {
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

	// PostInt8 获取指定参数的值，并将结果转为 int8
	PostInt8(key string) int8

	// PostUint8 获取指定参数的值，并将结果转为 uint8
	PostUint8(key string) uint8

	// PostInt16 获取指定参数的值，并将结果转为 int16
	PostInt16(key string) int16

	// PostUint16 获取指定参数的值，并将结果转为 uint16
	PostUint16(key string) uint16

	// PostInt32 获取指定参数的值，并将结果转为 int32
	PostInt32(key string) int32

	// PostUint32 获取指定参数的值，并将结果转为 uint32
	PostUint32(key string) uint32

	// PostInt64 获取指定参数的值，并将结果转为 int64
	PostInt64(key string) int64

	// PostUint64 获取指定参数的值，并将结果转为 uint64
	PostUint64(key string) uint64

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

// ContextDynamic 动态参数
type ContextDynamic interface {
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

// ContextFile 文件相关
type ContextFile interface {
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

// ContextWrite 响应相关
type ContextWrite interface {
	// Bytes 将数据写入响应
	Bytes(bytes []byte) (int, error)

	// Text 将数据写入响应中
	Text(value string) (int, error)

	// Textf 将数据写入响应中
	Textf(format string, a ...interface{}) error

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

// ContextCookie cookie 相关
type ContextCookie interface {

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

// ContextHook 钩子，一般用于中间件中
type ContextHook interface {
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

// Writer 实现 http.ResponseWriter
type Writer interface {
	http.ResponseWriter

	Writer() http.ResponseWriter

	SetWriter(w http.ResponseWriter)
}

// Server http 服务器
type Server interface {
	http.Handler

	// Start 根据配置调用 ListenAndServe 或者 ListenAndServeTLS 启动 http 服务，接收连接请求
	// addr: host:port，例如: ":8080"，"192.168.1.8:80"
	Start(addr string) error

	// HTTPServer 实际使用的 http 服务器
	HTTPServer() zerograceful.Server

	// SetTLS 指定 tls 证书，密钥路径
	// certFile: 证书路径
	// keyFile: 私钥路径
	SetTLS(certFile, keyFile string) bool

	// SetShutdownTimeout 设置优雅退出超时时间
	// 服务器会每隔500毫秒检查一次连接是否都断开处理完毕
	// 如果超过超时时间，就不再检查，直接退出
	//
	// ms: 单位：毫秒，当 <= 0 时无效，直接退出
	SetShutdownTimeout(ms int)

	// RegisterShutdownHandler 注册关闭函数
	// 按照注册的顺序调用这些函数
	// 所有已经添加的服务器都会响应这个函数
	RegisterShutdownHandler(f func())
}

// Router 路由管理器
type Router interface {

	// Prefix 设置前缀，设置前就已添加的路由不会有该前缀
	// 例如: prefix = "/blog"，则 "/user" -> "/blog/user"
	Prefix(prefix string)

	// Register 注册路由处理函数，以及中间件
	// method: HTTP Method，见 core/const.go Methodxxxx
	// path: 路径，以 "/" 开头，不可以为空
	// handles: 处理函数和路由级别中间件，匹配成功后会调用该函数
	Register(method, path string, handlers ...Handler) bool

	// Build 解析路由，包括动态参数，正则表达式，验证函数
	Build() bool

	// Lookup 查找路由
	Lookup(method, path string) ([]Handler, map[string]string)

	// RegisterRouterValidator 注册路由验证函数
	RegisterRouterValidator(name string, validator RouterValidator)

	// Validator 获取路由验证函数
	Validator(name string) RouterValidator
}

// Group 组路由，相同前缀的一组路由，共享相同的中间件
type Group interface {
	// Use 添加 Group 级别 中间件
	Use(handlers ...Handler) Group

	// Get method = "GET"
	Get(path string, handlers ...Handler) Group

	// Post method = "POST"
	Post(path string, handlers ...Handler) Group

	// Put method = "PUT"
	Put(path string, handlers ...Handler) Group

	// Delete method = "DELETE"
	Delete(path string, handlers ...Handler) Group

	// Head method = "HEAD"
	Head(path string, handlers ...Handler) Group

	// Patch method = "PATCH"
	Patch(path string, handlers ...Handler) Group

	// Options method = "OPTIONS"
	Options(path string, handlers ...Handler) Group
}

// RouteNode 一颗基数树的一个节点
type RouteNode interface {
	// Put 添加路由，路由不可重复
	Put(fullPath string, paths []string, height int, handlers ...Handler)

	// Build 解析路由，包括动态参数，正则表达式，验证函数。路由优化
	Build(router Router) bool

	// Lookup 查找路由
	Lookup(path string, dynamic map[string]string) ([]Handler, map[string]string)

	// Path 获取当前节点路径
	Path() string

	// Child 查找节点信息
	Child(path string) RouteNode

	// Children 获取节点列表
	Children() []RouteNode

	// Handlers 获取路由处理函数和中间件
	Handlers() []Handler

	// Reset 重置，清理所有数据
	Reset()

	RouteNodeInternal
}

// RouteNodeInternal 内部使用的一些接口
type RouteNodeInternal interface {

	// IsStatic 静态路由
	IsStatic() bool

	// IsDynamic 含有动态参数
	IsDynamic() bool

	// IsWildcard 含有通配符
	IsWildcard() bool

	// IsRegexp 含有正则表达式
	IsRegexp() bool

	// IsValidator 含有验证函数
	IsValidator() bool

	// IsHandler 是否有路由处理函数或者中间件
	IsHandler() bool

	// Flag 获取标记
	Flag() int

	// DynamicNum 获取动态节点数量
	DynamicNum() int
}
