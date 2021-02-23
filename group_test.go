package zeroapi_test

import (
	"testing"

	zeroapi "github.com/zerogo-hub/zero-api"
)

func TestGroup(t *testing.T) {
	a := zeroapi.NewApp()

	// 路由前缀自动添加 "/"
	g := zeroapi.NewGroup(a, "blog")

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
