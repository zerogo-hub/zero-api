package zeroapi

import "strconv"

// Post 包括 POST, PUT, PATCH
type Post interface {
	// Post 获取指定参数的值
	// 多个参数同名，只获取第一个
	// 例如 /user?id=Yaha&id=Gama
	// Post("id") 的结果为 "Yaha"
	Post(key string) string

	// PostStrings 获取指定参数的值
	// 例如 /user?id=Yaha&id=Gama
	// PostStrings("id") 的结果为 ["Yaha", "Gama"]
	PostStrings(key string) []string

	// PostEscape 获取指定参数的值，并对被编码的结果进行解码
	PostEscape(key string) string

	// PostBool 获取指定参数的值，并将结果转为 bool
	PostBool(key string) bool

	// PostInt32 获取指定参数的值，并将结果转为 int32
	PostInt32(key string) int32

	// PostInt64 获取指定参数的值，并将结果转为 int64
	PostInt64(key string) int64

	// PostFloat32 获取指定参数的值，并将结果转为 float32
	PostFloat32(key string) float32

	// PostFloat64 获取指定参数的值，并将结果转为 float64
	PostFloat64(key string) float64

	// PostDefault 获取指定参数的值，如果不存在，则返回默认值 def
	PostDefault(key, def string) string

	// PostBoolDefault 获取指定参数的值(结果转为 bool)，如果不存在，则返回默认值 def
	PostBoolDefault(key string, def bool) bool

	// PostInt32Default 获取指定参数的值(结果转为 int32)，如果不存在，则返回默认值 def
	PostInt32Default(key string, def int32) int32

	// PostInt64Default 获取指定参数的值(结果转为 int64)，如果不存在，则返回默认值 def
	PostInt64Default(key string, def int64) int64

	// PostFloat32Default 获取指定参数的值(结果转为 float32)，如果不存在，则返回默认值 def
	PostFloat32Default(key string, def float32) float32

	// PostFloat64Default 获取指定参数的值(结果转为 float64)，如果不存在，则返回默认值 def
	PostFloat64Default(key string, def float64) float64
}

func (ctx *context) Post(key string) string {
	// PostForm: 需要先调用 ParseForm()
	// contains the parsed form data from POST, PATCH, or PUT body parameters

	if err := ctx.req.ParseForm(); err != nil {
		return ""
	}

	return ctx.req.PostForm.Get(key)
}

func (ctx *context) PostStrings(key string) []string {
	// PostForm: 需要先调用 ParseForm()
	// contains the parsed form data from POST, PATCH, or PUT body parameters

	if err := ctx.req.ParseForm(); err != nil {
		return nil
	}

	if values, exist := ctx.req.PostForm[key]; exist {
		return values
	}

	return nil
}

func (ctx *context) PostEscape(key string) string {
	return unescape(ctx.Post(key))
}

func (ctx *context) PostBool(key string) bool {
	v, _ := strconv.ParseBool(ctx.Post(key))
	return v
}

func (ctx *context) PostInt32(key string) int32 {
	v, _ := strconv.ParseInt(ctx.Post(key), 10, 32)
	return int32(v)
}

func (ctx *context) PostInt64(key string) int64 {
	v, _ := strconv.ParseInt(ctx.Post(key), 10, 64)
	return v
}

func (ctx *context) PostFloat32(key string) float32 {
	v, _ := strconv.ParseFloat(ctx.Post(key), 32)
	return float32(v)
}

func (ctx *context) PostFloat64(key string) float64 {
	v, _ := strconv.ParseFloat(ctx.Post(key), 64)
	return v
}

func (ctx *context) PostDefault(key, def string) string {
	if value := ctx.Post(key); value != "" {
		return value
	}

	return def
}

func (ctx *context) PostBoolDefault(key string, def bool) bool {
	if value := ctx.Post(key); value != "" {
		result, err := strconv.ParseBool(value)
		if err != nil {
			return def
		}
		return result
	}

	return def
}

func (ctx *context) PostInt32Default(key string, def int32) int32 {
	if value := ctx.Post(key); value != "" {
		result, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return def
		}
		return int32(result)
	}

	return def
}

func (ctx *context) PostInt64Default(key string, def int64) int64 {
	if value := ctx.Post(key); value != "" {
		result, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return def
		}
		return result
	}

	return def
}

func (ctx *context) PostFloat32Default(key string, def float32) float32 {
	if value := ctx.Post(key); value != "" {
		result, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return def
		}
		return float32(result)
	}

	return def
}

func (ctx *context) PostFloat64Default(key string, def float64) float64 {
	if value := ctx.Post(key); value != "" {
		result, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return def
		}
		return result
	}

	return def
}
