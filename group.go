package zeroapi

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

type group struct {
	app    App
	prefix string

	// middlewares 组路由级别中间件
	middlewares []Handler
}

// NewGroup 创建一个组路由示例
func NewGroup(app App, prefix string) Group {
	if prefix != "" && prefix[0] != '/' {
		prefix = "/" + prefix
	}

	return &group{app: app, prefix: prefix}
}

// Use 添加 Group 级别 中间件
func (g *group) Use(handlers ...Handler) Group {
	g.middlewares = append(g.middlewares, handlers...)
	return g
}

// Get method = "GET"
func (g *group) Get(path string, handlers ...Handler) Group {
	g.app.Router().Register(MethodGet, g.prefix+path, handlers...)
	return g
}

// Post method = "POST"
func (g *group) Post(path string, handlers ...Handler) Group {
	g.app.Router().Register(MethodPost, g.prefix+path, handlers...)
	return g
}

// Put method = "PUT"
func (g *group) Put(path string, handlers ...Handler) Group {
	g.app.Router().Register(MethodPut, g.prefix+path, handlers...)
	return g
}

// Delete method = "DELETE"
func (g *group) Delete(path string, handlers ...Handler) Group {
	g.app.Router().Register(MethodDelete, g.prefix+path, handlers...)
	return g
}

// Head method = "HEAD"
func (g *group) Head(path string, handlers ...Handler) Group {
	g.app.Router().Register(MethodHead, g.prefix+path, handlers...)
	return g
}

// Patch method = "PATCH"
func (g *group) Patch(path string, handlers ...Handler) Group {
	g.app.Router().Register(MethodPatch, g.prefix+path, handlers...)
	return g
}

// Options method = "OPTIONS"
func (g *group) Options(path string, handlers ...Handler) Group {
	g.app.Router().Register(MethodOptions, g.prefix+path, handlers...)
	return g
}
