// Package zeroapi 简单的 api 框架
package zeroapi

const (
	// VERSION 框架版本号
	VERSION = "0.1.0"
)

// New 生成一个应用实例
func New() App {
	a := NewApp()
	return a
}

// Default 生成默认的应用实例
func Default() App {
	a := NewApp()

	return a
}
