package router_test

import (
	"testing"

	zeroapp "github.com/zerogo-hub/zero-api/app"
	zerorouter "github.com/zerogo-hub/zero-api/router"
)

func TestGroup(t *testing.T) {
	a := zeroapp.NewApp()

	// 路由前缀自动添加 "/"
	g := zerorouter.NewGroup(a, "blog")

	// 使用中间件
	g.Use(emptyHandle)

	// 获取列表
	g.Get("/", emptyHandle)

	// 新增
	g.Post("/", emptyHandle)

	// 修改
	g.Put("/", emptyHandle)
	g.Patch("/", emptyHandle)

	// 删除
	g.Delete("/", emptyHandle)

	g.Head("/", emptyHandle)
	g.Options("/", emptyHandle)

	if !a.Router().Build() {
		t.Fatal("build failed")
	}
}
