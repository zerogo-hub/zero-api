package zeroapi

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
