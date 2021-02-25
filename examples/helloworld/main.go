package main

import (
	"os"

	zeroapi "github.com/zerogo-hub/zero-api"
)

func helloworldHandle(ctx zeroapi.Context) {
	pid := os.Getpid()
	ctx.Textf("`ctrl+c` to close, `kill %d` to shutdown, `kill -USR2 %d` to restart", pid, pid)
}

func main() {
	a := zeroapi.Default()

	a.Get("/", helloworldHandle)

	a.Server().HTTPServer().ListenSignal()

	a.Run("127.0.0.1:8877")
}
