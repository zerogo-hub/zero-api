package main

import (
	zeroapi "github.com/zerogo-hub/zero-api"
	app "github.com/zerogo-hub/zero-api/app"
)

type Student struct {
	ID   int64  `json:'id'`
	Name string `json:'name'`
	Age  int32  `json:'age'`
}

func indexHandle(ctx zeroapi.Context) {
	id := ctx.QueryInt64("id")
	name := ctx.Post("name")
	age := ctx.PostInt32("age")

	ctx.JSON(map[string]interface{}{
		"id":   id,
		"name": name,
		"age":  age,
	})
}

func index2Handle(ctx zeroapi.Context) {
	var student Student
	if err := ctx.Body(&student); err != nil {
		ctx.Text(err.Error())
		return
	}

	ctx.JSON(student)
}

func main() {
	a := app.Default()

	a.Post("/", indexHandle)
	a.Post("/2", index2Handle)

	server := a.Server()

	// 监听信号，比如优雅关闭
	server.HTTPServer().ListenSignal()

	// 在退出前清理
	server.RegisterShutdownHandler(func() {
		a.Logger().Info("done before exit")
	})

	if err := a.Run("127.0.0.1:8877"); err != nil {
		a.Logger().Error(err.Error())
	}
}
