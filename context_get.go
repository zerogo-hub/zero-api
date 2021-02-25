package zeroapi

import "strconv"

// Get 只包括 GET
type Get interface {
	// Get 获取指定参数的值
	// 多个参数同名，只获取第一个
	// 例如 /user?id=Yaha&id=Gama
	// Get("id") 的结果为 "Yaha"
	Get(key string) string

	// GetEscape 获取指定参数的值，并对被编码的结果进行解码
	GetEscape(key string) string

	// GetBool 获取指定参数的值，并将结果转为 bool
	GetBool(key string) bool

	// GetInt32 获取指定参数的值，并将结果转为 int32
	GetInt32(key string) int32

	// GetInt64 获取指定参数的值，并将结果转为 int64
	GetInt64(key string) int64

	// GetFloat32 获取指定参数的值，并将结果转为 float32
	GetFloat32(key string) float32

	// GetFloat64 获取指定参数的值，并将结果转为 float64
	GetFloat64(key string) float64

	// GetDefault 获取指定参数的值，如果不存在，则返回默认值 def
	GetDefault(key, def string) string

	// GetBoolDefault 获取指定参数的值(结果转为 bool)，如果不存在，则返回默认值 def
	GetBoolDefault(key string, def bool) bool

	// GetInt32Default 获取指定参数的值(结果转为 int32)，如果不存在，则返回默认值 def
	GetInt32Default(key string, def int32) int32

	// GetInt64Default 获取指定参数的值(结果转为 int64)，如果不存在，则返回默认值 def
	GetInt64Default(key string, def int64) int64

	// GetFloat32Default 获取指定参数的值(结果转为 float32)，如果不存在，则返回默认值 def
	GetFloat32Default(key string, def float32) float32

	// GetFloat64Default 获取指定参数的值(结果转为 float64)，如果不存在，则返回默认值 def
	GetFloat64Default(key string, def float64) float64
}

func (ctx *context) Get(key string) string {
	return ctx.req.URL.Query().Get(key)
}

func (ctx *context) GetEscape(key string) string {
	return unescape(ctx.Get(key))
}

func (ctx *context) GetBool(key string) bool {
	v, _ := strconv.ParseBool(ctx.Get(key))
	return v
}

func (ctx *context) GetInt32(key string) int32 {
	v, _ := strconv.ParseInt(ctx.Get(key), 10, 32)
	return int32(v)
}

func (ctx *context) GetInt64(key string) int64 {
	v, _ := strconv.ParseInt(ctx.Get(key), 10, 64)
	return v
}

func (ctx *context) GetFloat32(key string) float32 {
	v, _ := strconv.ParseFloat(ctx.Get(key), 32)
	return float32(v)
}

func (ctx *context) GetFloat64(key string) float64 {
	v, _ := strconv.ParseFloat(ctx.Get(key), 64)
	return v
}

func (ctx *context) GetDefault(key, def string) string {
	if value := ctx.Get(key); value != "" {
		return value
	}

	return def
}

func (ctx *context) GetBoolDefault(key string, def bool) bool {
	if value := ctx.Get(key); value != "" {
		result, err := strconv.ParseBool(value)
		if err != nil {
			return def
		}
		return result
	}

	return def
}

func (ctx *context) GetInt32Default(key string, def int32) int32 {
	if value := ctx.Get(key); value != "" {
		result, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return def
		}
		return int32(result)
	}

	return def
}

func (ctx *context) GetInt64Default(key string, def int64) int64 {
	if value := ctx.Get(key); value != "" {
		result, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return def
		}
		return result
	}

	return def
}

func (ctx *context) GetFloat32Default(key string, def float32) float32 {
	if value := ctx.Get(key); value != "" {
		result, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return def
		}
		return float32(result)
	}

	return def
}

func (ctx *context) GetFloat64Default(key string, def float64) float64 {
	if value := ctx.Get(key); value != "" {
		result, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return def
		}
		return result
	}

	return def
}
