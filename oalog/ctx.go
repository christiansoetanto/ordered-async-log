package oalog

import (
	"context"
	"github.com/rs/xid"
	"github.com/tokopedia/tdk/v2/go/log"
	"strconv"
)

type contextKey string

const (
	contextKeyRequestID = contextKey("tkpd-log-request-id")
	contextKeyProcessNo = contextKey("process-no")
)

func InitLogContext(ctx context.Context) context.Context {
	// Check if request id is in context
	if GetCtxRequestID(ctx) == "" {
		return SetCtxRequestID(ctx, xid.New().String())
	}
	return ctx
}

// SetCtxRequestID set context with request id
func SetCtxRequestID(ctx context.Context, uuid string) context.Context {
	return context.WithValue(ctx, contextKeyRequestID, uuid)
}

// GetCtxRequestID get request id from context
func GetCtxRequestID(ctx context.Context) string {
	id, ok := ctx.Value(contextKeyRequestID).(string)
	if !ok {
		return ""
	}
	return id
}
func SetCtxProcessNo(ctx context.Context, no int) context.Context {
	return context.WithValue(ctx, contextKeyProcessNo, strconv.Itoa(no))
}
func GetCtxProcessNo(ctx context.Context) string {
	id, ok := ctx.Value(contextKeyProcessNo).(string)
	if !ok {
		return ""
	}
	return id
}

type Metadata = log.KV
