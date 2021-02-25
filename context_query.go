package zeroapi

import (
	"net/url"
	"strconv"
)

// Query 包括 GET, POST, PUT
type Query interface {
	// Query 获取指定参数的值
	// 多个参数同名，只获取第一个
	// 例如 /user?id=Yaha&id=Gama
	// Query("id") 的结果为 "Yaha"
	Query(key string) string

	// QueryAll 获取所有参数的值
	// 例如 /user?id=Yaha&id=Gama&age=18
	// QueryAll() 的结果为 {"id": ["Yaha", "Gama"], "age": ["18"]}
	QueryAll() map[string][]string

	// QueryStrings 获取指定参数的值
	// 例如 /user?id=Yaha&id=Gama
	// QueryStrings("id") 的结果为 ["Yaha", "Gama"]
	QueryStrings(key string) []string

	// QueryEscape 获取指定参数的值，并对被转码的结果进行还原
	QueryEscape(key string) string

	// QueryBool 获取指定参数的值，并将结果转为 bool
	QueryBool(key string) bool

	// QueryInt32 获取指定参数的值，并将结果转为 int32
	QueryInt32(key string) int32

	// QueryInt64 获取指定参数的值，并将结果转为 int64
	QueryInt64(key string) int64

	// QueryFloat32 获取指定参数的值，并将结果转为 float32
	QueryFloat32(key string) float32

	// QueryFloat64 获取指定参数的值，并将结果转为 float64
	QueryFloat64(key string) float64

	// QueryDefault 获取指定参数的值，如果不存在，则返回默认值 def
	QueryDefault(key, def string) string

	// QueryBoolDefault 获取指定参数的值(结果转为 bool)，如果不存在，则返回默认值 def
	QueryBoolDefault(key string, def bool) bool

	// QueryInt32Default 获取指定参数的值(结果转为 int32)，如果不存在，则返回默认值 def
	QueryInt32Default(key string, def int32) int32

	// QueryInt64Default 获取指定参数的值(结果转为 int64)，如果不存在，则返回默认值 def
	QueryInt64Default(key string, def int64) int64

	// QueryFloat32Default 获取指定参数的值(结果转为 float32)，如果不存在，则返回默认值 def
	QueryFloat32Default(key string, def float32) float32

	// QueryFloat64Default 获取指定参数的值(结果转为 float64)，如果不存在，则返回默认值 def
	QueryFloat64Default(key string, def float64) float64
}

// queryAll 获取所有的参数值，内部使用
func (ctx *context) queryAll() (map[string][]string, bool) {
	if err := ctx.req.ParseForm(); err != nil {
		return nil, false
	}
	// Form: 需要先调用 ParseForm()
	// including both the URL field's query parameters and the POST or PUT form data
	if form := ctx.req.Form; len(form) > 0 {
		return form, true
	}
	// Form 中没有数据，则使用了 PATCH ?
	//
	// PostForm: 需要先调用 ParseForm()
	// contains the parsed form data from POST, PATCH, or PUT body parameters
	if form := ctx.req.PostForm; len(form) > 0 {
		return form, true
	}

	return nil, false
}

func (ctx *context) Query(key string) string {
	if form, exist := ctx.queryAll(); exist {
		if values := form[key]; len(values) > 0 {
			return values[0]
		}
	}

	return ""
}

func (ctx *context) QueryAll() map[string][]string {
	form, _ := ctx.queryAll()
	return form
}

func (ctx *context) QueryStrings(key string) []string {
	if form, exist := ctx.queryAll(); exist {
		if values := form[key]; len(values) > 0 {
			return values
		}
	}

	return nil
}

// Unescape 解码
// 解码以下编码的结果:
// js: escape, encodeURI, encodeURIComponent
// go: QueryEscape
func unescape(value string) string {
	if value == "" {
		return ""
	}
	result, err := url.QueryUnescape(value)
	if err != nil {
		return value
	}
	return result
}

func (ctx *context) QueryEscape(key string) string {
	return unescape(ctx.Query(key))
}

func (ctx *context) QueryBool(key string) bool {
	v, _ := strconv.ParseBool(ctx.Query(key))
	return v
}

func (ctx *context) QueryInt32(key string) int32 {
	v, _ := strconv.ParseInt(ctx.Query(key), 10, 32)
	return int32(v)
}

func (ctx *context) QueryInt64(key string) int64 {
	v, _ := strconv.ParseInt(ctx.Query(key), 10, 64)
	return v
}

func (ctx *context) QueryFloat32(key string) float32 {
	v, _ := strconv.ParseFloat(ctx.Query(key), 32)
	return float32(v)
}

func (ctx *context) QueryFloat64(key string) float64 {
	v, _ := strconv.ParseFloat(ctx.Query(key), 64)
	return v
}

func (ctx *context) QueryDefault(key, def string) string {
	if value := ctx.Query(key); value != "" {
		return value
	}

	return def
}

func (ctx *context) QueryBoolDefault(key string, def bool) bool {
	if value := ctx.Query(key); value != "" {
		result, err := strconv.ParseBool(value)
		if err != nil {
			return def
		}
		return result
	}

	return def
}

func (ctx *context) QueryInt32Default(key string, def int32) int32 {
	if value := ctx.Query(key); value != "" {
		result, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return def
		}
		return int32(result)
	}

	return def
}

func (ctx *context) QueryInt64Default(key string, def int64) int64 {
	if value := ctx.Query(key); value != "" {
		result, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return def
		}
		return result
	}

	return def
}

func (ctx *context) QueryFloat32Default(key string, def float32) float32 {
	if value := ctx.Query(key); value != "" {
		result, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return def
		}
		return float32(result)
	}

	return def
}

func (ctx *context) QueryFloat64Default(key string, def float64) float64 {
	if value := ctx.Query(key); value != "" {
		result, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return def
		}
		return result
	}

	return def
}
