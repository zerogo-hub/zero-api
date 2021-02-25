package zeroapi

// Hook 钩子，一般用于中间件中
type Hook interface {
	// AppendAfter 添加处理函数，这些函数 在中间件和路由函数执行完毕之后执行，如果中间有中断，则不会执行
	// 先入后出顺序执行，越先加入的函数越后执行
	AppendAfter(hook HookHandler)

	// AppendEnd 添加处理函数，这些函数 在处理都完成之后执行，无论中间是否有中断和异常，都会执行
	// 先入后出顺序执行，越先加入的函数越后执行
	AppendEnd(hook HookHandler)

	// RunAfter 执行通过 AppendAfter 加入的处理函数
	RunAfter()

	// RunEnd 执行通过 AppendEnd 加入的处理函数
	RunEnd()
}

func (ctx *context) AppendAfter(hook HookHandler) {
	ctx.afters = append(ctx.afters, hook)
}

func (ctx *context) AppendEnd(hook HookHandler) {
	ctx.ends = append(ctx.ends, hook)
}

func (ctx *context) RunAfter() {
	run(ctx.afters)
}

func (ctx *context) RunEnd() {
	defer ctx.app.ReleaseContext(ctx)
	run(ctx.ends)
}

func run(hooks []HookHandler) {
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
