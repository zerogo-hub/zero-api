package context

import (
	"strconv"
)

func (ctx *context) Get(key string) string {
	return ctx.req.URL.Query().Get(key)
}

func (ctx *context) Gets(key string) []string {
	return ctx.req.URL.Query()[key]
}

func (ctx *context) GetEscape(key string) string {
	return unescape(ctx.Get(key))
}

func (ctx *context) GetBool(key string) bool {
	v, _ := strconv.ParseBool(ctx.Get(key))
	return v
}

func (ctx *context) GetInt8(key string) int8 {
	v, _ := strconv.ParseInt(ctx.Get(key), 10, 32)
	return int8(v)
}

func (ctx *context) GetUint8(key string) uint8 {
	v, _ := strconv.ParseInt(ctx.Get(key), 10, 32)
	return uint8(v)
}

func (ctx *context) GetInt16(key string) int16 {
	v, _ := strconv.ParseInt(ctx.Get(key), 10, 32)
	return int16(v)
}

func (ctx *context) GetUint16(key string) uint16 {
	v, _ := strconv.ParseInt(ctx.Get(key), 10, 32)
	return uint16(v)
}

func (ctx *context) GetInt32(key string) int32 {
	v, _ := strconv.ParseInt(ctx.Get(key), 10, 32)
	return int32(v)
}

func (ctx *context) GetUint32(key string) uint32 {
	v, _ := strconv.ParseInt(ctx.Get(key), 10, 32)
	return uint32(v)
}

func (ctx *context) GetInt64(key string) int64 {
	v, _ := strconv.ParseInt(ctx.Get(key), 10, 64)
	return v
}

func (ctx *context) GetUint64(key string) uint64 {
	v, _ := strconv.ParseInt(ctx.Get(key), 10, 64)
	return uint64(v)
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
