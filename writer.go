package zeroapi

import (
	"net/http"
	"sync"
)

// Writer 实现 http.ResponseWriter
type Writer interface {
	http.ResponseWriter

	Writer() http.ResponseWriter

	SetWriter(w http.ResponseWriter)
}

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
func acquireWriter() Writer {
	return writerPool.Get().(*writer)
}

// releaseWriter 将 writer 放入池中
func releaseWriter(w Writer) {
	writerPool.Put(w)
}

func init() {
	writerPool = &sync.Pool{}
	writerPool.New = func() interface{} {
		return &writer{}
	}
}
