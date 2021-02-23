package zeroapi

import (
	"bytes"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/zerogo-hub/zero-helper/crypto"
	"github.com/zerogo-hub/zero-helper/time"
)

func (ctx *context) Cookie(name string, opts ...CookieOption) (string, error) {
	oname := name

	if ctx.app.IsCookieEncode() {
		handler := ctx.app.CookieEncodeHandler()
		name = handler(name)
	}

	cookie, err := ctx.req.Cookie(name)
	if err != nil {
		return "", err
	}

	cookie.Name = oname

	if ctx.app.IsCookieEncode() {
		handler := ctx.app.CookieDecodeHandler()
		value, err := handler(cookie.Value)
		if err != nil {
			return "", err
		}

		cookie.Value = value
	}

	for _, opt := range opts {
		opt(cookie)
	}

	val, err := url.QueryUnescape(cookie.Value)
	return val, err
}

func (ctx *context) SetCookie(name, value string, opts ...CookieOption) {
	cookie := &http.Cookie{Name: name, Value: url.QueryEscape(value)}

	for _, opt := range opts {
		opt(cookie)
	}

	// 默认存在 1 小时
	if cookie.MaxAge == 0 {
		cookie.MaxAge = 3600
	}

	if ctx.app.IsCookieEncode() {
		handler := ctx.app.CookieEncodeHandler()
		cookie.Name = handler(cookie.Name)
		cookie.Value = handler(cookie.Value)
	}

	http.SetCookie(ctx.res, cookie)
}

func (ctx *context) RemoveCookie(name string, opts ...CookieOption) {
	ctx.SetCookie(name, "", CookieMaxAge(-1))
}

func (ctx *context) SetHTTPCookie(cookie *http.Cookie) {
	if cookie == nil {
		panic("Cookie cannot be empty")
	}

	http.SetCookie(ctx.res, cookie)
}

func (ctx *context) HTTPCookies() []*http.Cookie {
	return ctx.req.Cookies()
}

// CookieMaxAge ..
// maxAge: 见 https://tools.ietf.org/html/rfc6265#section-4.1.2.2
// 		 = 0: 表示不指定存活时间
// 		 < 0: 表示立即删除
// 		 > 0: cookie 生存时间，单位秒
func CookieMaxAge(maxAge int) CookieOption {
	return func(cookie *http.Cookie) {
		cookie.MaxAge = maxAge
	}
}

// CookiePath path: https://tools.ietf.org/html/rfc6265#section-4.1.2.4
func CookiePath(path string) CookieOption {
	return func(cookie *http.Cookie) {
		if path != "" {
			cookie.Path = path
		}
	}
}

// CookieDomain domain: https://tools.ietf.org/html/rfc6265#section-4.1.2.3
func CookieDomain(domain string) CookieOption {
	return func(cookie *http.Cookie) {
		if domain != "" {
			cookie.Domain = domain
		}
	}
}

// CookieSecure secure: https://tools.ietf.org/html/rfc6265#section-4.1.2.5
func CookieSecure(secure bool) CookieOption {
	return func(cookie *http.Cookie) {
		cookie.Secure = secure
	}
}

// CookieHTTPOnly secure: https://tools.ietf.org/html/rfc6265#section-4.1.2.6
func CookieHTTPOnly(httpOnly bool) CookieOption {
	return func(cookie *http.Cookie) {
		cookie.HttpOnly = httpOnly
	}
}

// CookieSign 对 cookie 进行签名
func CookieSign(signKey string) CookieOption {
	return func(cookie *http.Cookie) {
		if cookie.Name == "" {
			return
		}

		timestamp := strconv.Itoa(int(time.Now()))

		buf := buffer()
		defer releaseBuffer(buf)

		buf.WriteString(cookie.Name)
		buf.WriteString(cookie.Value)
		buf.WriteString(timestamp)

		sign := crypto.HmacMd5(buf.String(), signKey)

		buf.Reset()
		buf.WriteString(cookie.Value)
		buf.WriteString("|")
		buf.WriteString(timestamp)
		buf.WriteString("|")
		buf.WriteString(sign)

		cookie.Value = buf.String()
	}
}

// CookieVerify 对有签名的 cookie 进行验证
func CookieVerify(signKey string) CookieOption {
	return func(cookie *http.Cookie) {
		if cookie.Value == "" {
			return
		}

		l := strings.Split(cookie.Value, "|")
		if len(l) != 3 {
			// cookie 值被篡改
			cookie.Value = ""
			return
		}

		value := l[0]
		timestamp := l[1]
		sign := l[2]

		buf := buffer()
		defer releaseBuffer(buf)

		buf.WriteString(cookie.Name)
		buf.WriteString(value)
		buf.WriteString(timestamp)
		calcSign := crypto.HmacMd5(buf.String(), signKey)

		if calcSign != sign {
			// cookie 值被篡改
			cookie.Value = ""
		} else {
			cookie.Value = value
		}
	}
}

var bufferPool *sync.Pool

// buffer 从池中获取 buffer
func buffer() *bytes.Buffer {
	buff := bufferPool.Get().(*bytes.Buffer)
	buff.Reset()
	return buff
}

// releaseBuffer 将 buff 放入池中
func releaseBuffer(buff *bytes.Buffer) {
	bufferPool.Put(buff)
}

func init() {
	bufferPool = &sync.Pool{}
	bufferPool.New = func() interface{} {
		return &bytes.Buffer{}
	}
}
