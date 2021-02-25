package zeroapi

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/zerogo-hub/zero-helper/bytes"
	"google.golang.org/protobuf/proto"
)

// Write 响应相关
type Write interface {
	// Bytes 将数据写入响应
	Bytes(bytes []byte) (int, error)

	// Text 将数据写入响应中
	Text(value string) (int, error)

	// Textf 将数据写入响应中
	Textf(format string, a ...interface{})

	// Map map 转 text
	Map(obj interface{}) (int, error)

	// JSON 将数据转为 JSON 格式写入响应
	JSON(obj interface{}) (int, error)

	// XML 将数据转为 XML 格式写入响应
	XML(obj interface{}) (int, error)

	// HTML 发送 html 响应
	HTML(html string) (int, error)

	// HTMLf 发送 html 响应
	HTMLf(format string, a ...interface{}) (int, error)

	// Protobuf 将数据装为 google protobuf 格式，写入响应
	Protobuf(obj interface{}) (int, error)

	// Size 响应的数据大小
	Size() int64

	// Redirect 重定向
	// httpCode: HTTP Code， 需要在 3xx 范围内, 比如 301, 302, 303 ... 308
	// url: 重定向后的地址
	Redirect(httpCode int, url string) error

	// Flush 将数据推向客户端
	Flush()

	// Push HTTP/2 服务器推送
	Push(value string, opts *http.PushOptions) error

	// AutoContentType 根据给定的文件类型，自动设置 Content-Type
	// .json -> app/json
	// fileExt: 文件后缀名，例如 .json
	AutoContentType(fileExt string)

	// Message 传递 {"code": xx, "message": xxx}
	Message(code int, message ...string) (int, error)
}

func (ctx *context) Bytes(bytes []byte) (int, error) {
	var size int
	var err error

	size, err = ctx.res.Write(bytes)

	if err != nil {
		return 0, err
	}

	ctx.responseSize += int64(size)

	return size, nil
}

func (ctx *context) Text(value string) (int, error) {
	var size int
	var err error

	size, err = ctx.res.Write(bytes.StringToBytes(value))

	if err != nil {
		return 0, err
	}

	ctx.responseSize += int64(size)
	ctx.SetHeader("Content-Type", "text/plain;charset=utf-8")

	return size, nil
}

func (ctx *context) Textf(format string, a ...interface{}) {
	ctx.Text(fmt.Sprintf(format, a...))
}

func (ctx *context) Map(obj interface{}) (int, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return 0, err
	}

	return ctx.Bytes(bytes)
}

func (ctx *context) JSON(obj interface{}) (int, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return 0, err
	}

	ctx.SetHeader("Content-Type", "application/json;charset=utf-8")

	return ctx.Bytes(bytes)
}

func (ctx *context) XML(obj interface{}) (int, error) {
	bytes, err := xml.Marshal(obj)
	if err != nil {
		return 0, err
	}

	ctx.SetHeader("Content-Type", "application/xml;charset=utf-8")

	return ctx.Bytes(bytes)
}

func (ctx *context) HTML(html string) (int, error) {
	ctx.SetHeader("Content-Type", "text/html;charset=utf-8")

	return ctx.Bytes(bytes.StringToBytes(template.HTMLEscapeString(html)))
}

func (ctx *context) HTMLf(format string, a ...interface{}) (int, error) {
	return ctx.HTML(fmt.Sprintf(format, a...))
}

func (ctx *context) Protobuf(obj interface{}) (int, error) {
	bytes, err := proto.Marshal(obj.(proto.Message))
	if err != nil {
		return 0, err
	}

	ctx.SetHeader("Content-Type", "application/x-protobuf;charset=utf-8")

	return ctx.Bytes(bytes)
}

func (ctx *context) Size() int64 {
	return ctx.responseSize
}

func (ctx *context) Redirect(httpCode int, url string) error {
	if httpCode < http.StatusMultipleChoices || httpCode > http.StatusPermanentRedirect {
		return errors.New("httpCode should be in the 3xx, like 301, 302 etc")
	}

	ctx.Stopped()
	http.Redirect(ctx.res, ctx.req, url, httpCode)

	return nil
}

func (ctx *context) Flush() {
	if flusher, ok := ctx.res.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (ctx *context) Push(value string, opts *http.PushOptions) error {
	if push, ok := ctx.res.(http.Pusher); ok {
		push.Push(value, opts)
	}

	return http.ErrNotSupported
}

func (ctx *context) AutoContentType(fileExt string) {
	if !strings.HasPrefix(fileExt, ".") {
		fileExt = "." + fileExt
	}

	if ct := mime.TypeByExtension(fileExt); ct != "" {
		ctx.AddHeader("Content-Type", ct)
	}
}

func (ctx *context) Message(code int, message ...string) (int, error) {
	result := make(map[string]string)

	result["code"] = strconv.Itoa(code)
	if len(message) > 0 {
		result["message"] = message[0]
	}

	return ctx.Map(result)
}
