package context

import (
	"errors"
	"io"

	"google.golang.org/protobuf/proto"

	zerojson "github.com/zerogo-hub/zero-helper/json"
)

func (ctx *context) Body(in interface{}) error {
	b, err := ctx.readBody()
	if err != nil {
		return err
	}

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

func (ctx *context) readBody() ([]byte, error) {
	return io.ReadAll(ctx.req.Body)
}
