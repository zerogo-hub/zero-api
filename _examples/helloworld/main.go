package main

import (
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
	app "github.com/zerogo-hub/zero-api/app"
)

func helloworldHandle(ctx zeroapi.Context) {
	pid := os.Getpid()
	if err := ctx.Textf("`ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", pid, pid); err != nil {
		ctx.App().Logger().Error(err.Error())
	}
}

func main() {
	a := app.Default()

	a.Get("/", helloworldHandle)

	// 监听信号，比如优雅关闭
	a.Server().HTTPServer().ListenSignal()

	if err := a.Run("127.0.0.1:8877"); err != nil {
		a.Logger().Error(err.Error())
	}
}
