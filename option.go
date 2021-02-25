package zeroapi

import (
	"github.com/zerogo-hub/zero-helper/logger"
)

var (
	// defaultFileMaxMemory 用于限制使用内存大小 multipart/form-data，比如文件上传
	defaultFileMaxMemory = int64(32 * 1024 * 1024) // 32M
)

// config app 配置
type config struct {
	// version 框架版本号
	version string

	// fileMaxMemory 文件系统使用的最大内存
	fileMaxMemory int64

	// logger 日志管理器
	logger logger.Logger

	// cookieEncode 对 cookie 键值编码函数
	cookieEncode CookieEncodeHandler

	// cookieDecode 对 cookie 键值解码函数
	cookieDecode CookieDecodeHandler
}

func defaultConfig() *config {
	return &config{
		version:       VERSION,
		fileMaxMemory: defaultFileMaxMemory,
		logger:        logger.NewSampleLogger(),
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

// WithFileMaxMemory 设置文件系统使用的最大内存
func WithFileMaxMemory(fileMaxMemory int64) Option {
	return func(config *config) {
		config.fileMaxMemory = fileMaxMemory
	}
}

// WithLogger 设置日志
func WithLogger(logger logger.Logger) Option {
	return func(config *config) {
		config.logger = logger
	}
}

// WithCookieHandler 设置 cookie 编码与解码函数
func WithCookieHandler(encoder CookieEncodeHandler, decoder CookieDecodeHandler) Option {
	return func(config *config) {
		config.cookieEncode = encoder
		config.cookieDecode = decoder
	}
}
