package router_test

import (
	"testing"

	app "github.com/zerogo-hub/zero-api/app"
	router "github.com/zerogo-hub/zero-api/router"
)

func TestRouteLookupStatic(t *testing.T) {
	route := router.NewRoute()
	route.Insert("/blog/name", emptyHandle)
	route.Build(nil)

	if handlers, dynamic := route.Lookup("/blog/name"); handlers == nil || dynamic != nil {
		t.Fatal("invalid 1")
	}

	if handlers, dynamic := route.Lookup("/blog/name/add"); handlers != nil || dynamic != nil {
		t.Fatal("invalid 2")
	}

	if handlers, dynamic := route.Lookup("/blog/na"); handlers != nil || dynamic != nil {
		t.Fatal("invalid 3")
	}
}

func TestRouteLookupNotFound(t *testing.T) {
	route := router.NewRoute()
	route.Insert("/blog/name", emptyHandle)
	route.Build(nil)

	if handlers, dynamic := route.Lookup("/blog/10001/name"); handlers != nil || dynamic != nil {
		t.Fatal("invalid 1")
	}
}

func TestRouteLookupDynamic(t *testing.T) {
	a := app.NewApp()
	r := a.Router()

	r.RegisterRouterValidator("isNum", isNum)

	route := router.NewRoute()
	route.Insert("/blog/:id/name", emptyHandle)
	route.Build(nil)

	// 正常解析
	if _, dynamic := route.Lookup("/blog/10001/name"); len(dynamic) == 0 || dynamic["id"] != "10001" {
		t.Fatal("invalid 1")
	}

	if handlers, _ := route.Lookup("/blog/10001/account"); handlers != nil {
		t.Fatal("invalid 2")
	}

	// 动态参数为最后一个
	route.Reset()
	route.Insert("/blog/:id(\\d+)", emptyHandle)
	route.Build(nil)
	if _, dynamic := route.Lookup("/blog/10001"); len(dynamic) == 0 || dynamic["id"] != "10001" {
		t.Fatal("invalid 3")
	}

	// 使用正则表达式判断动态参数值
	route.Reset()
	route.Insert("/blog/:id(\\d+)", emptyHandle)
	route.Build(nil)

	if handlers, dynamic := route.Lookup("/blog/10001"); handlers == nil || len(dynamic) == 0 || dynamic["id"] != "10001" {
		t.Fatal("failed")
	}

	if handlers, dynamic := route.Lookup("/blog/abc"); handlers != nil || len(dynamic) != 0 {
		t.Fatal("failed")
	}

	// 使用验证函数判断动态参数值
	route.Reset()
	route.Insert("/blog/:id|isNum|", emptyHandle)
	route.Build(r)

	if handlers, dynamic := route.Lookup("/blog/10001"); handlers == nil || len(dynamic) == 0 || dynamic["id"] != "10001" {
		t.Fatal("failed")
	}

	if handlers, dynamic := route.Lookup("/blog/abc"); handlers != nil || len(dynamic) != 0 {
		t.Fatal("failed")
	}
}

func TestRouteLookupDynamicWildcard(t *testing.T) {
	route := router.NewRoute()
	route.Insert("/blog/:id/*/name", emptyHandle)
	route.Build(nil)

	handlers, dynamic := route.Lookup("/blog/10001/abc/d/name")
	if handlers == nil || len(dynamic) == 0 || dynamic["id"] != "10001" {
		t.Fatal("invalid 1")
	}
}
