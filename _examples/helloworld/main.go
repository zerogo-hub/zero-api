package main

import (
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
	app "github.com/zerogo-hub/zero-api/app"
)

func helloWorldHandle(ctx zeroapi.Context) {
	pid := os.Getpid()
	if err := ctx.Textf("`ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", pid, pid); err != nil {
		ctx.App().Logger().Error(err.Error())
	}
}

func main() {
	a := app.Default()

	a.Get("/", helloWorldHandle)

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
