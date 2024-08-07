package context

import (
	"bytes"
	"errors"
	"io"

	"google.golang.org/protobuf/proto"

	zerojson "github.com/zerogo-hub/zero-helper/json"
)

func (ctx *context) Body(in interface{}) error {
	b, release, err := ctx.ReadBody(true)
	if err != nil {
		return err
	}
	release()

	switch ctx.ContentType() {
	case "application/json":
		return zerojson.Unmarshal(b, in)
	case "application/x-protobuf":
		msg, ok := in.(proto.Message)
		if !ok {
			return errors.New("not a protobuf message")
		}
		return proto.Unmarshal(b, msg)
	}

	return nil
}

var emptyFunc = func() {}

// ReadBody reads the request body
func (ctx *context) ReadBody(isMultiTimes bool) ([]byte, func(), error) {
	data, err := io.ReadAll(ctx.req.Body)
	if err != nil {
		return nil, nil, err
	}

	if !isMultiTimes {
		return data, emptyFunc, nil
	}

	return data, func() {
		ctx.req.Body = io.NopCloser(bytes.NewBuffer(data))
	}, nil
}
