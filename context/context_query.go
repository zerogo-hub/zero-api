package context

import (
	"net/url"
	"strconv"
)

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
