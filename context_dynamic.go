package zeroapi

import "errors"

// Dynamic 动态参数
type Dynamic interface {
	// Dynamic 获取动态参数的值
	// 示例:
	// 定义路由: /blog/:id
	// 调用路由: /blog/1001
	// params = {id: "1001"}
	// Dynamic("id") -> 1001
	Dynamic(key string) string

	// SetDynamic 设置动态参数，key 的格式为 "param" 或者 ":param"
	SetDynamic(key string, value string) error

	// SetDynamics 替换动态参数
	SetDynamics(dynamics map[string]string)
}

func (ctx *context) Dynamic(key string) string {
	if len(key) == 0 {
		return ""
	}

	// 兼容判断
	if key[0] == ':' {
		key = key[1:]
	}

	if ctx.dynamics != nil {
		if value, exist := ctx.dynamics[key]; exist {
			return value
		}
	}

	return ""
}

func (ctx *context) SetDynamic(key string, value string) error {
	if len(key) == 0 {
		return errors.New("Parameter key cannot be empty")
	}

	// 兼容判断
	if key[0] == ':' {
		key = key[1:]
	}

	if ctx.dynamics == nil {
		ctx.dynamics = make(map[string]string)
	}

	ctx.dynamics[key] = value

	return nil
}

func (ctx *context) SetDynamics(dynamics map[string]string) {
	ctx.dynamics = dynamics
}
