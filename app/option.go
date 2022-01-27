package app

import (
	zeroapi "github.com/zerogo-hub/zero-api"

	zerologger "github.com/zerogo-hub/zero-helper/logger"
)

var (
	// defaultMaxMemory 用于限制使用内存大小 multipart/form-data，比如文件上传
	// 也会被 http.MaxBytesReader 调用
	defaultMaxMemory = int64(32 * 1024 * 1024) // 32M
)

// config app 配置
type config struct {
	// version 框架版本号
	version string

	// maxMemory 最大内存
	maxMemory int64

	// logger 日志管理器
	logger zerologger.Logger

	// cookieEncode 对 cookie 键值编码函数
	cookieEncode zeroapi.CookieEncodeHandler

	// cookieDecode 对 cookie 键值解码函数
	cookieDecode zeroapi.CookieDecodeHandler
}

func defaultConfig() *config {
	return &config{
		version:   zeroapi.VERSION,
		maxMemory: defaultMaxMemory,
		logger:    zerologger.NewSampleLogger(),
	}
}

// Option app 配置选项
type Option func(config *config)

// WithVersion 设置框架版本号
func WithVersion(version string) Option {
	return func(config *config) {
		config.version = version
	}
}

// WithMaxMemory 设置使用的最大内存
func WithMaxMemory(maxMemory int64) Option {
	return func(config *config) {
		config.maxMemory = maxMemory
	}
}

// WithLogger 设置日志
func WithLogger(logger zerologger.Logger) Option {
	return func(config *config) {
		config.logger = logger
	}
}

// WithCookieHandler 设置 cookie 编码与解码函数
func WithCookieHandler(encoder zeroapi.CookieEncodeHandler, decoder zeroapi.CookieDecodeHandler) Option {
	return func(config *config) {
		config.cookieEncode = encoder
		config.cookieDecode = decoder
	}
}
