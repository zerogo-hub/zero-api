package router_test

import (
	"testing"

	zeroapi "github.com/zerogo-hub/zero-api"
	zeroapp "github.com/zerogo-hub/zero-api/app"
)

func TestRouterRegister(t *testing.T) {
	a := zeroapp.NewApp()
	r := a.Router()

	// 错误的路由前缀
	r.Prefix("")

	// 路由前缀自动添加 "/"
	r.Prefix("blog")

	// 注册错误的路由
	if r.Register(zeroapi.MethodGet, "") {
		t.Fatal("invalid path")
	}

	// 注册错误的路由
	if r.Register(zeroapi.MethodGet, "/list") {
		t.Fatal("miss handlers")
	}

	// 注册正确的路由
	if !r.Register(zeroapi.MethodGet, "/list", emptyHandle) {
		t.Fatal("failed")
	}
}

func TestRouterBuildFailed(t *testing.T) {
	a := zeroapp.NewApp()
	r := a.Router()

	r.Register(zeroapi.MethodGet, "/list/:id(\\d+", emptyHandle)

	if r.Build() {
		t.Fatal("invalid regexp")
	}
}

func TestRouterBuildSuccess(t *testing.T) {
	a := zeroapp.NewApp()
	r := a.Router()
	r.RegisterRouterValidator("isNum", isNum)

	// 重复添加会被忽略
	r.RegisterRouterValidator("isNum", isNum)

	r.Register(zeroapi.MethodGet, "/list/:id(\\d+)|isNum|", emptyHandle)

	if !r.Build() {
		t.Fatal("failed")
	}
}

func TestRouterLookup(t *testing.T) {
	a := zeroapp.NewApp()
	r := a.Router()
	r.RegisterRouterValidator("less4", less4)

	r.Register(zeroapi.MethodGet, "/list/:id(\\d+)|less4|", emptyHandle)

	if !r.Build() {
		t.Fatal("build failed")
	}

	// 正确的找到
	if handlers, dynamic := r.Lookup(zeroapi.MethodGet, "/list/101"); handlers == nil || len(dynamic) == 0 || dynamic["id"] != "101" {
		t.Fatal("lookup failed")
	}

	// 不匹配的路由，不能通过正则表达式
	if handlers, dynamic := r.Lookup(zeroapi.MethodGet, "/list/abcd"); handlers != nil || len(dynamic) > 0 {
		t.Fatal("lookup failed")
	}

	// 不匹配的路由，不能通过验证函数检查
	if handlers, dynamic := r.Lookup(zeroapi.MethodGet, "/list/1001"); handlers != nil || len(dynamic) > 0 {
		t.Fatal("lookup failed")
	}

	// 不匹配的路由，错误的 Method
	if handlers, dynamic := r.Lookup(zeroapi.MethodPost, "/list/10001"); handlers != nil || len(dynamic) > 0 {
		t.Fatal("lookup failed")
	}
}

func TestRouterLookup2(t *testing.T) {
	a := zeroapp.NewApp()
	r := a.Router()

	g := a.Group("/account")
	g.Post("/v1/signup", emptyHandle)
	g.Post("/v1/signin", emptyHandle)
	g.Post("/v1/signout", emptyHandle)
	g.Put("/v1/password", emptyHandle)

	if !r.Build() {
		t.Fatal("build failed")
	}

	if handlers, _ := r.Lookup(zeroapi.MethodPost, "/account/v1/signin"); handlers == nil {
		t.Fatal("lookup /account/v1/signin failed")
	}

	if handlers, _ := r.Lookup(zeroapi.MethodPut, "/account/v1/password"); handlers == nil {
		t.Fatal("lookup /account/v1/password failed")
	}
}

func TestRouterLookup3(t *testing.T) {
	a := zeroapp.NewApp()
	r := a.Router()

	a.Get("/app", emptyHandle)
	a.Get("/app/category", emptyHandle)
	a.Get("/app/category/v1", emptyHandle)
	if !r.Build() {
		t.Fatal("build failed")
	}

	// 查找不存在的路由
	if handlers, _ := r.Lookup(zeroapi.MethodGet, "/app/recharge/v1"); handlers != nil {
		t.Fatal("lookup /app/recharge/v1 failed")
	}
	if handlers, _ := r.Lookup(zeroapi.MethodGet, "/app/category2/v1"); handlers != nil {
		t.Fatal("lookup /app/recharge/v1 failed")
	}
}
