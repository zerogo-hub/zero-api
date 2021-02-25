package zeroapi

import (
	"strings"
)

// Route 路由，每一个 Route 表示一颗基数树，每种 HTTP Method 一个实例
type Route interface {
	// Insert 添加路由，路由不可重复
	Insert(path string, handlers ...Handler)

	// Build 解析路由，包括动态参数，正则表达式，验证函数。路由优化
	Build(router Router) bool

	// Lookup 查找路由
	Lookup(path string) ([]Handler, map[string]string)

	// Child 查找节点信息
	Child(path string) RouteNode

	// Children 获取节点列表
	Children() []RouteNode

	// Reset 重置，清理所有数据
	Reset()
}

// route 实现一颗基数树
type route struct {
	// root 基数树根节点
	root RouteNode
}

// NewRoute ..
func NewRoute() Route {
	return &route{root: new(routeNode)}
}

// Insert 添加路由，路由不可重复
func (re *route) Insert(path string, handlers ...Handler) {
	paths := buildPath(path)
	re.root.Put(path, paths, 0, handlers...)
}

// Build 解析路由，包括动态参数，正则表达式，验证函数
func (re *route) Build(router Router) bool {
	return re.root.Build(router)
}

// Lookup 查找路由
func (re *route) Lookup(path string) ([]Handler, map[string]string) {
	return re.root.Lookup(path, nil)
}

// Child 查找节点信息
func (re *route) Child(path string) RouteNode {
	for _, child := range re.root.Children() {
		if child.Path() == path {
			return child
		}
	}

	if re.root.Path() == path {
		return re.root
	}

	return nil
}

// Children 获取节点列表
func (re *route) Children() []RouteNode {
	return re.root.Children()
}

// Reset 重置，清理所有数据
func (re *route) Reset() {
	re.root.Reset()
}

func buildPath(path string) []string {
	if path == "/" {
		return []string{"/"}
	}

	paths := strings.Split(path, "/")

	out := make([]string, 0, len(paths)-1)

	for _, p := range paths {
		if p == "" {
			continue
		}

		out = append(out, "/"+p)
		if p[0] == WildcardCharacter {
			break
		}
	}

	return out
}
