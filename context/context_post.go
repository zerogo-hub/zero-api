package context

import (
	"strconv"
)

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

func (ctx *context) PostInt8(key string) int8 {
	v, _ := strconv.ParseInt(ctx.Post(key), 10, 32)
	return int8(v)
}

func (ctx *context) PostUint8(key string) uint8 {
	v, _ := strconv.ParseInt(ctx.Post(key), 10, 32)
	return uint8(v)
}

func (ctx *context) PostInt16(key string) int16 {
	v, _ := strconv.ParseInt(ctx.Post(key), 10, 32)
	return int16(v)
}

func (ctx *context) PostUint16(key string) uint16 {
	v, _ := strconv.ParseInt(ctx.Post(key), 10, 32)
	return uint16(v)
}

func (ctx *context) PostInt32(key string) int32 {
	v, _ := strconv.ParseInt(ctx.Post(key), 10, 32)
	return int32(v)
}

func (ctx *context) PostUint32(key string) uint32 {
	v, _ := strconv.ParseInt(ctx.Post(key), 10, 32)
	return uint32(v)
}

func (ctx *context) PostInt64(key string) int64 {
	v, _ := strconv.ParseInt(ctx.Post(key), 10, 64)
	return v
}

func (ctx *context) PostUint64(key string) uint64 {
	v, _ := strconv.ParseInt(ctx.Post(key), 10, 64)
	return uint64(v)
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
