package context

import (
	"net/http"
	"sync"

	zeroapi "github.com/zerogo-hub/zero-api"
)

type writer struct {
	http.ResponseWriter
}

func (w *writer) Writer() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *writer) SetWriter(sw http.ResponseWriter) {
	w.ResponseWriter = sw
}

var writerPool *sync.Pool

// AcquireWriter 从池中获取 Writer
func AcquireWriter() zeroapi.Writer {
	return writerPool.Get().(*writer)
}

// ReleaseWriter 将 writer 放入池中
func ReleaseWriter(w zeroapi.Writer) {
	writerPool.Put(w)
}

func init() {
	writerPool = &sync.Pool{}
	writerPool.New = func() interface{} {
		return &writer{}
	}
}
