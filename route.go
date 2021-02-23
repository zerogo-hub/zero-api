package zeroapi

import (
	"regexp"
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

	// Path 获取当前节点路径
	Path() string

	// Child 查找子节点信息
	Child(path string) RouteNode

	// Children 获取子节点列表
	Children() []RouteNode

	// Reset 重置，清理所有数据
	Reset()
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

	// Child 查找子节点信息
	Child(path string) RouteNode

	// Children 获取子节点列表
	Children() []RouteNode

	// Flag 获取标记
	Flag() int

	// Handlers 获取路由处理函数和中间件
	Handlers() []Handler

	// DynamicNum 获取动态节点数量
	DynamicNum() int

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

	// Reset 重置，清理所有数据
	Reset()
}

const (
	// DynamicCharacter 动态路由符号，比如 /user/:name
	DynamicCharacter = ':'

	// WildcardCharacter 通配符，比如 /blog/hi/*
	WildcardCharacter = '*'
)

const (
	// STATIC 静态路由
	STATIC = 0

	// DYNAMIC 含有动态参数标记 :
	DYNAMIC = 2

	// WILDCARD 含有通配符 *
	WILDCARD = 2 << 1

	// REGEXP 含有正则表达式
	REGEXP = 2 << 2

	// VALIDATOR 是否含有验证函数
	VALIDATOR = 2 << 3
)

// route 实现一颗基数树
type route struct {
	// root 基数树根节点
	root RouteNode
}

// routeNode 一颗基数树的一个节点
type routeNode struct {

	// fullPath 路由全路径
	fullPath string

	// path 路由部分路径，由该节点持有
	path string

	// handlers 路由处理函数 + 路由级别中间件
	handlers []Handler

	// validators 参数校验
	validators []RouterValidator

	// flag 用于标记节点是否含有 动态参数(:param),通配符(*),正则表达式(regexp),验证函数(validator)
	flag int

	// dynamicName 动态参数名称，假设动态参数 :id，则 dynamicName = id
	dynamicName string

	// dynamicNum 本节点 + 子节点动态参数个数
	dynamicNum int

	// pattern 编译好的正则表达式
	pattern *regexp.Regexp

	// children 子节点
	children []RouteNode
}

// NewRoute ..
func NewRoute() Route {
	return &route{root: new(routeNode)}
}

func buildPath(path string) []string {
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

// Path 获取当前节点路径
func (re *route) Path() string {
	return re.root.Path()
}

// Child 查找子节点信息
func (re *route) Child(path string) RouteNode {
	for _, child := range re.root.Children() {
		if child.Path() == path {
			return child
		}
	}

	return nil
}

// Children 获取子节点列表
func (re *route) Children() []RouteNode {
	return re.root.Children()
}

// Reset 重置，清理所有数据
func (re *route) Reset() {
	re.root.Reset()
}

// IsStatic 静态路由
func (rn *routeNode) IsStatic() bool {
	return rn.flag == 0
}

// IsDynamic 含有动态参数
func (rn *routeNode) IsDynamic() bool {
	return rn.flag&DYNAMIC != 0
}

// IsWildcard 含有通配符
func (rn *routeNode) IsWildcard() bool {
	return rn.flag&WILDCARD != 0
}

// IsRegexp 含有正则表达式
func (rn *routeNode) IsRegexp() bool {
	return rn.flag&REGEXP != 0
}

// IsValidator 含有验证函数
func (rn *routeNode) IsValidator() bool {
	return rn.flag&VALIDATOR != 0
}

// IsHandler 是否有路由处理函数或者中间件
func (rn *routeNode) IsHandler() bool {
	return rn.handlers != nil && len(rn.handlers) > 0
}

// Child 查找子节点信息
func (rn *routeNode) Child(path string) RouteNode {
	for _, child := range rn.children {
		if child.Path() == path {
			return child
		}
	}

	return nil
}

// Children 获取子节点列表
func (rn *routeNode) Children() []RouteNode {
	return rn.children
}

// Flag 获取标记
func (rn *routeNode) Flag() int {
	return rn.flag
}

// Handlers 获取路由处理函数和中间件
func (rn *routeNode) Handlers() []Handler {
	return rn.handlers
}

// DynamicNum 获取动态节点数量
func (rn *routeNode) DynamicNum() int {
	return rn.dynamicNum
}

// Path 获取当前节点路径
func (rn *routeNode) Path() string {
	return rn.path
}

// put 添加路由
//
// fullPath 完整路径，例如 /blog/:id/borrow
func (rn *routeNode) Put(fullPath string, paths []string, height int, handlers ...Handler) {
	if len(handlers) == 0 {
		return
	}

	if len(paths) == height || rn.IsWildcard() {
		// 本次路由的最终节点
		rn.fullPath = fullPath
		rn.handlers = handlers
		return
	}

	path := paths[height]

	// 是否已存在对应的子节点
	child := rn.child(path)
	if child == nil {
		child = newChild(path)
		rn.children = append(rn.children, child)
	}

	child.Put(fullPath, paths, height+1, handlers...)
}

func newChild(path string) RouteNode {
	flag := STATIC

	if len(path) > 1 {
		if path[1] == DynamicCharacter {
			flag |= DYNAMIC
		} else if path[1] == WildcardCharacter {
			flag |= WILDCARD
		}
	}

	child := &routeNode{path: path, flag: flag}

	return child
}

// child 在子节点中查找已存在的节点
func (rn *routeNode) child(path string) RouteNode {
	for _, child := range rn.children {
		if child.Path() == path || child.IsWildcard() {
			return child
		}
	}

	return nil
}

func (rn *routeNode) Build(router Router) bool {
	if rn.IsWildcard() {
		return true
	}

	if rn.IsDynamic() {
		if !(rn.parseRegexp() && rn.parseValidator(router) && rn.parseDynamic()) {
			return false
		}
	}

	rn.merge()

	// 解析子节点
	for _, child := range rn.children {
		if !child.Build(router) {
			return false
		}
	}

	rn.countDynamicNum()

	return true
}

// parseRegexp 解析当前节点 path 上的正则表达式
//
// 一个节点只包含一个正则表达式
func (rn *routeNode) parseRegexp() bool {
	// 示例: /blog/list/:id(^\d+$)
	pos := strings.Index(rn.path, "(")

	if pos == -1 {
		return true
	}

	posEnd := strings.Index(rn.path, ")")
	if posEnd == -1 {
		// 缺失右括号
		return false
	}
	if pos+1 >= posEnd {
		// )(
		return false
	}

	regexpExpress := rn.path[pos+1 : posEnd]
	rn.pattern = regexp.MustCompile(regexpExpress)
	rn.flag |= REGEXP

	return true
}

// parseValidator 解析当前节点 path 上的验证函数
//
// 验证函数必须现在 Router 中注册
func (rn *routeNode) parseValidator(router Router) bool {
	if router == nil {
		return true
	}

	// 示例: /blog/list/:id|isNum|less4|
	pos := strings.Index(rn.path, "|")

	if pos == -1 {
		return true
	}

	posEnd := strings.LastIndex(rn.path, "|")
	if posEnd == -1 || pos == posEnd {
		// 必须包含在 |...| 中间
		return false
	}

	handlerNames := strings.Split(rn.path[pos+1:posEnd], "|")
	if len(handlerNames) == 0 || handlerNames[0] == "" {
		return false
	}

	rn.validators = make([]RouterValidator, 0, len(handlerNames))

	for _, handlerName := range handlerNames {
		handler := router.Validator(handlerName)
		if handler == nil {
			return false
		}

		rn.validators = append(rn.validators, handler)
	}

	rn.flag |= VALIDATOR

	return true
}

// parseDynamic 解析当前节点 path 上的动态参数
func (rn *routeNode) parseDynamic() bool {
	// 示例: /blog/article/:id(^\d+$)|less4|/del

	// 当前节点的 path = /:id(^\d+$)|less4|
	// rn.dynamicName = id

	// 开头两个符号为 /:，所以从 2 开始
	i := 2
	for ; i < len(rn.path); i++ {
		c := rn.path[i]

		if c == '|' || c == '(' {
			break
		}
	}

	rn.dynamicName = rn.path[2:i]

	return true
}

// merge 合并，如果只有一个子节点，且子节点是 STATIC 的，则合并
func (rn *routeNode) merge() {
	if len(rn.children) != 1 || !rn.IsStatic() || rn.IsHandler() {
		return
	}

	child := rn.children[0]

	if !child.IsStatic() {
		return
	}

	// 拼接 path
	if rn.path == "" {
		rn.path = child.Path()
	} else {
		rn.path += child.Path()
	}

	rn.flag |= child.Flag()
	rn.children = child.Children()
	rn.handlers = child.Handlers()

	rn.merge()
}

func (rn *routeNode) countDynamicNum() {

	dynamicNum := 0

	for _, child := range rn.children {
		// 只保留 dynamicNum 最多的那个
		if child.DynamicNum() > dynamicNum {
			dynamicNum = child.DynamicNum()
		}
	}

	if rn.IsDynamic() {
		dynamicNum++
	}

	rn.dynamicNum = dynamicNum
}

func (rn *routeNode) Lookup(path string, dynamic map[string]string) ([]Handler, map[string]string) {

	if rn.IsWildcard() {
		return rn.handlers, dynamic
	}

	if rn.IsDynamic() {
		return rn.lookupByDynamic(path, dynamic)
	}

	return rn.lookupByStatic(path, dynamic)
}

func (rn *routeNode) lookupByStatic(path string, dynamic map[string]string) ([]Handler, map[string]string) {
	if rn.path == path {
		return rn.handlers, dynamic
	}

	// rn.path = /users，path = /user
	// 当前节点 rn 不匹配 path
	if len(rn.path) > len(path) {
		return nil, nil
	}

	// rn.path = /user，path = /user/add
	// 从子节点中匹配，childPath = /add
	childPath := path[len(rn.path):]
	if childPath[0] != '/' {
		return nil, nil
	}

	for _, child := range rn.children {
		if handlers, dynamic := child.Lookup(childPath, dynamic); handlers != nil {
			return handlers, dynamic
		}
	}

	return nil, nil
}

func (rn *routeNode) lookupByDynamic(path string, dynamic map[string]string) ([]Handler, map[string]string) {

	// rn.path = /:id，path = /1001/add
	if dynamic == nil {
		dynamic = make(map[string]string, rn.dynamicNum)
	}

	// 获取 id 值，id = 1001
	pos := strings.Index(path[1:], "/")
	dynamicValueEnd := pos
	if pos < 0 {
		// rn.path = /:id, path = /1001
		dynamicValueEnd = len(path) - 1
	}
	dynamicValue := path[1 : dynamicValueEnd+1]

	if !rn.checkDynamicValueValid(dynamicValue) {
		return nil, nil
	}

	// rn.dynamicName = id
	dynamic[rn.dynamicName] = dynamicValue

	// 如果 path[1:] 没有 '/' 或者 '/' 在最后一个，表示该节点是最后一个节点了
	if pos == -1 || pos == len(path)-1 {
		return rn.handlers, dynamic
	}

	// 向子节点查找
	childPath := path[pos+1:]

	for _, child := range rn.children {
		if handlers, dynamic := child.Lookup(childPath, dynamic); handlers != nil {
			return handlers, dynamic
		}
	}

	return nil, nil
}

func (rn *routeNode) checkDynamicValueValid(dynamicValue string) bool {

	if rn.IsRegexp() && !rn.checkRegexp(dynamicValue) {
		return false
	}

	if rn.IsValidator() && !rn.checkValidator(dynamicValue) {
		return false
	}

	return true
}

func (rn *routeNode) checkRegexp(dynamicValue string) bool {
	return rn.pattern.MatchString(dynamicValue)
}

func (rn *routeNode) checkValidator(dynamicValue string) bool {
	for _, validator := range rn.validators {
		if !validator(dynamicValue) {
			return false
		}
	}

	return true
}

// Reset 重置，清理所有数据
func (rn *routeNode) Reset() {
	rn.fullPath = ""
	rn.path = ""
	rn.handlers = nil
	rn.validators = nil
	rn.flag = STATIC
	rn.dynamicName = ""
	rn.dynamicNum = 0
	rn.pattern = nil
	rn.children = nil
}
