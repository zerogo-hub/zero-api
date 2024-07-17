package context

func (ctx *context) Header(key string) string {
	return ctx.req.Header.Get(key)
}

func (ctx *context) AddHeader(key, value string) {
	if key != "" && value != "" {
		ctx.res.Writer().Header().Add(key, value)
	}
}

func (ctx *context) SetHeader(key, value string) {
	if key != "" && value != "" {
		ctx.res.Writer().Header().Set(key, value)
	}
}

func (ctx *context) DelHeader(key string) {
	if key != "" {
		ctx.res.Writer().Header().Del(key)
	}
}

func (ctx *context) ContentType() string {
	return ctx.Header("Content-Type")
}
