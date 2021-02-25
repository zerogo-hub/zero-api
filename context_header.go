package zeroapi

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
