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

// acquireWriter 从池中获取 Writer
func acquireWriter() zeroapi.Writer {
	return writerPool.Get().(*writer)
}

// releaseWriter 将 writer 放入池中
func releaseWriter(w zeroapi.Writer) {
	writerPool.Put(w)
}

func init() {
	writerPool = &sync.Pool{}
	writerPool.New = func() interface{} {
		return &writer{}
	}
}
