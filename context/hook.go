package context

import (
	zeroapi "github.com/zerogo-hub/zero-api"
)

func (ctx *context) AppendAfter(hook zeroapi.HookHandler) {
	ctx.afters = append(ctx.afters, hook)
}

func (ctx *context) AppendEnd(hook zeroapi.HookHandler) {
	ctx.ends = append(ctx.ends, hook)
}

func (ctx *context) RunAfter() {
	run(ctx.afters)
}

func (ctx *context) RunEnd() {
	defer ctx.app.ReleaseContext(ctx)
	run(ctx.ends)
}

func run(hooks []zeroapi.HookHandler) {
	if len(hooks) == 0 {
		return
	}

	// 从尾巴开始执行
	for i := len(hooks) - 1; i >= 0; i-- {
		if err := hooks[i](); err != nil {
			return
		}
	}
}
