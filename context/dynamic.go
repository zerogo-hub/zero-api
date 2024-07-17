package context

import (
	"errors"
)

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
		return errors.New("parameter key cannot be empty")
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
