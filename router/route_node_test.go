package router_test

import (
	"regexp"
	"testing"

	zeroapi "github.com/zerogo-hub/zero-api"
	app "github.com/zerogo-hub/zero-api/app"
	router "github.com/zerogo-hub/zero-api/router"
)

func emptyHandle(zeroapi.Context) {}

func TestRouteNodeStatic(t *testing.T) {
	route := router.NewRoute()

	route.Insert("/blog", emptyHandle)
	route.Insert("/blog/a/b/c/d", emptyHandle)

	route.Build(nil)

	root := route.Child("/blog")
	if root.Path() != "/blog" {
		t.Fatal("invalid path")
	}

	children := route.Children()
	if len(children) != 1 {
		t.Fatal("children' size must be 1")
	}

	if child1 := route.Child("/a/b/c/d"); child1 == nil {
		t.Fatal("child1 not found")
	}

	if child2 := route.Child("/fake"); child2 != nil {
		t.Fatal("child2 exist")
	}
}

func TestRouteNodeRoot(t *testing.T) {
	route := router.NewRoute()
	route.Insert("/", emptyHandle)
	route.Build(nil)

	root := route.Child("/")
	if root.Path() != "/" {
		t.Fatal("invalid root path")
	}
}

func TestRouteNodeDynamic(t *testing.T) {
	route := router.NewRoute()

	route.Insert("/blog/:id/borrow", emptyHandle)
	route.Insert("/blog/:id/name", emptyHandle)

	route.Build(nil)

	root := route.Child("/blog")
	if root.Path() != "/blog" {
		t.Fatal("invalid path")
	}

	children := route.Children()
	if len(children) != 1 {
		t.Fatal("invalid children")
	}

	child := route.Child("/:id")
	if child == nil {
		t.Fatalf("invalid child: %s", "/:id")
	}

	if len(child.Children()) != 2 {
		t.Fatal("invalid children")
	}

	if fakeChild := child.Child("/fake"); fakeChild != nil {
		t.Fatal("fake child exist")
	}
}

func TestRouteNodeMultiDynamic(t *testing.T) {
	route := router.NewRoute()

	route.Insert("/blog/:id/borrow", emptyHandle)
	route.Insert("/blog/:id/:account/:app/name", emptyHandle)

	route.Build(nil)

	root := route.Child("/blog")
	if root.Path() != "/blog" {
		t.Fatal("invalid path")
	}

	children := route.Children()
	if len(children) != 1 {
		t.Fatal("invalid children")
	}

	child := route.Child("/:id")
	if child == nil {
		t.Fatalf("invalid child: %s", "/:id")
	}

	if len(child.Children()) != 2 {
		t.Fatal("child: invalid children")
	}

	child2 := child.Child("/:account")
	if child2 == nil || child2.IsHandler() {
		t.Fatalf("invalid child: %s", "/:account")
	}

	if len(child2.Children()) != 1 {
		t.Fatal("child2: invalid children")
	}

	child3 := child2.Child("/:app")
	if child3 == nil || child3.IsHandler() {
		t.Fatalf("invalid child: %s", "/:app")
	}

	if len(child3.Children()) != 1 {
		t.Fatal("child3: invalid children")
	}

	child4 := child3.Child("/name")
	if child4 == nil || !child4.IsHandler() {
		t.Fatalf("invalid child: %s", "/name")
	}

	if len(child4.Children()) != 0 {
		t.Fatal("child4: invalid children")
	}
}

func TestRouteNodeDynamicNum(t *testing.T) {
	route := router.NewRoute()

	route.Insert("/blog-1/:id/borrow/:account", emptyHandle)
	route.Insert("/blog-2/:id/:account/:app/name", emptyHandle)
	route.Insert("/blog-2/add/:id/:account/:app/:name/:version", emptyHandle)

	if len(route.Children()) != 2 {
		t.Fatal("route: invalid children")
	}

	child2 := route.Child("/blog-2")
	if child2 == nil {
		t.Fatal("invalid child2")
	}

	if len(child2.Children()) != 2 {
		t.Fatal("child2: invalid children")
	}

	route.Build(nil)
}

func TestRouteNodeStaticDynamic(t *testing.T) {
	route := router.NewRoute()

	route.Insert("/blog/user/add", emptyHandle)
	route.Insert("/blog/user/:id/del", emptyHandle)
	route.Insert("/blog/invalidHandler")

	route.Build(nil)

	root := route.Child("/blog/user")
	if root.Path() != "/blog/user" {
		t.Fatal("invalid route path")
	}

	child1 := route.Child("/add")
	if child1 == nil || !child1.IsHandler() {
		t.Fatal("invalid child1")
	}

	child2 := route.Child("/:id")
	if child2 == nil || !child2.IsDynamic() || child2.IsHandler() {
		t.Fatal("invalid child2")
	}

	if len(child2.Children()) != 1 {
		t.Fatal("child2: invalid child")
	}

	child3 := child2.Child("/del")
	if child3 == nil || !child3.IsHandler() || !child3.IsStatic() {
		t.Fatalf("child2: invalid child: %s", "/del")
	}
}

func TestRouteNodeWildcard(t *testing.T) {
	route := router.NewRoute()

	route.Insert("/blog/notfound/*/abc", emptyHandle)

	route.Build(nil)

	root := route.Child("/blog/notfound")
	if root.Path() != "/blog/notfound" {
		t.Fatal("invalid path")
	}

	child := route.Child("/*")
	if child == nil || !child.IsWildcard() || !child.IsHandler() {
		t.Fatal("invalid child")
	}

	if len(child.Children()) > 0 {
		t.Fatal("child: invalid children")
	}
}

func TestRouteDynamicParseRegexp(t *testing.T) {
	route := router.NewRoute()

	// 缺失右括号
	route.Insert("/blog/list/:id(^\\d+$", emptyHandle)
	if route.Build(nil) {
		t.Fatal("miss )")
	}

	// 左右括号对调
	route.Reset()
	route.Insert("/blog/list/:id)^\\d+$(", emptyHandle)
	if route.Build(nil) {
		t.Fatal("invalid )(")
	}

	// 正常
	route.Reset()
	route.Insert("/blog/list/:id(^\\d+$)", emptyHandle)
	if !route.Build(nil) {
		t.Fatal("invalid regexp")
	}
}

func isNum(s string) bool {
	result, _ := regexp.MatchString("\\d+", s)
	return result
}

func less4(s string) bool {
	return len(s) < 4
}

func TestRouteDynamicParseValidator(t *testing.T) {
	a := app.NewApp()
	r := a.Router()

	r.RegisterRouterValidator("isNum", isNum)

	route := router.NewRoute()

	// 没有验证函数
	route.Insert("/blog/list/:id", emptyHandle)
	if !route.Build(r) {
		t.Fatal("no validator")
	}

	// 缺少 | 将验证函数包裹
	route.Reset()
	route.Insert("/blog/list/:id|isNum", emptyHandle)
	if route.Build(r) {
		t.Fatal("miss \"|\"")
	}

	// 缺少验证函数
	route.Reset()
	route.Insert("/blog/list/:id||", emptyHandle)
	if route.Build(r) {
		t.Fatal("miss validator")
	}

	// 不存在的验证函数
	route.Reset()
	route.Insert("/blog/list/:id|isNum|less4|", emptyHandle)
	if route.Build(r) {
		t.Fatal("validator not found")
	}

	// 正常路由
	route.Reset()
	route.Insert("/blog/list/:id|isNum|", emptyHandle)
	if !route.Build(r) {
		t.Fatal("failed")
	}
}
