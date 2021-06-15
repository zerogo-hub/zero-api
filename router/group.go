package router

import (
	zeroapi "github.com/zerogo-hub/zero-api"
)

type group struct {
	app    zeroapi.App
	prefix string

	// middlewares 组路由级别中间件
	middlewares []zeroapi.Handler
}

// NewGroup 创建一个组路由示例
func NewGroup(app zeroapi.App, prefix string) zeroapi.Group {
	if prefix != "" && prefix[0] != '/' {
		prefix = "/" + prefix
	}

	return &group{app: app, prefix: prefix}
}

// Use 添加 Group 级别 中间件
func (g *group) Use(handlers ...zeroapi.Handler) zeroapi.Group {

	for _, handler := range handlers {
		if handler != nil {
			g.middlewares = append(g.middlewares, handler)
		}
	}

	return g
}

func (g *group) groupHandlers(handlers ...zeroapi.Handler) []zeroapi.Handler {
	lenGroupMiddlewares := len(g.middlewares)
	lenRouteHandlers := len(handlers)

	_handlers := make([]zeroapi.Handler, lenGroupMiddlewares+lenRouteHandlers)

	copy(_handlers, g.middlewares)
	copy(_handlers[lenGroupMiddlewares:], handlers)

	return _handlers
}

// Get method = "GET"
func (g *group) Get(path string, handlers ...zeroapi.Handler) zeroapi.Group {
	g.app.Get(g.prefix+path, g.groupHandlers(handlers...)...)
	return g
}

// Post method = "POST"
func (g *group) Post(path string, handlers ...zeroapi.Handler) zeroapi.Group {
	g.app.Post(g.prefix+path, g.groupHandlers(handlers...)...)
	return g
}

// Put method = "PUT"
func (g *group) Put(path string, handlers ...zeroapi.Handler) zeroapi.Group {
	g.app.Put(g.prefix+path, g.groupHandlers(handlers...)...)
	return g
}

// Delete method = "DELETE"
func (g *group) Delete(path string, handlers ...zeroapi.Handler) zeroapi.Group {
	g.app.Delete(g.prefix+path, g.groupHandlers(handlers...)...)
	return g
}

// Head method = "HEAD"
func (g *group) Head(path string, handlers ...zeroapi.Handler) zeroapi.Group {
	g.app.Head(g.prefix+path, g.groupHandlers(handlers...)...)
	return g
}

// Patch method = "PATCH"
func (g *group) Patch(path string, handlers ...zeroapi.Handler) zeroapi.Group {
	g.app.Patch(g.prefix+path, g.groupHandlers(handlers...)...)
	return g
}

// Options method = "OPTIONS"
func (g *group) Options(path string, handlers ...zeroapi.Handler) zeroapi.Group {
	g.app.Options(g.prefix+path, g.groupHandlers(handlers...)...)
	return g
}
