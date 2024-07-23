package context

import (
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"

	zerobytes "github.com/zerogo-hub/zero-helper/bytes"
	zerojson "github.com/zerogo-hub/zero-helper/json"
)

func (ctx *context) Bytes(bytes []byte) (int, error) {
	var size int
	var err error

	size, err = ctx.res.Writer().Write(bytes)

	if err != nil {
		return 0, err
	}

	ctx.responseSize += int64(size)

	return size, nil
}

func (ctx *context) Text(value string) (int, error) {
	var size int
	var err error

	size, err = ctx.res.Writer().Write(zerobytes.StringToBytes(value))

	if err != nil {
		return 0, err
	}

	ctx.responseSize += int64(size)
	ctx.SetHeader("Content-Type", "text/plain;charset=utf-8")

	return size, nil
}

func (ctx *context) Textf(format string, a ...interface{}) error {
	_, err := ctx.Text(fmt.Sprintf(format, a...))
	return err
}

func (ctx *context) Map(obj interface{}) (int, error) {
	bytes, err := zerojson.Marshal(obj)
	if err != nil {
		return 0, err
	}

	return ctx.Bytes(bytes)
}

func (ctx *context) JSON(obj interface{}) (int, error) {
	bytes, err := zerojson.Marshal(obj)
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

	return ctx.Bytes(zerobytes.StringToBytes(template.HTMLEscapeString(html)))
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
	http.Redirect(ctx.res.Writer(), ctx.req, url, httpCode)

	return nil
}

func (ctx *context) Flush() {
	if flusher, ok := ctx.res.Writer().(http.Flusher); ok {
		flusher.Flush()
	}
}

func (ctx *context) Push(value string, opts *http.PushOptions) error {
	if push, ok := ctx.res.Writer().(http.Pusher); ok {
		if err := push.Push(value, opts); err != nil {
			return err
		}
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
