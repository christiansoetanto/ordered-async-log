package oalog

import (
	"context"
	"github.com/tokopedia/tdk/v2/go/log"
	"math/rand"
	"time"
)

func execDebug(ctx context.Context, message string) {
	message = removePII(message)
	metadata := Metadata{
		string(contextKeyProcessNo): GetCtxProcessNo(ctx),
	}
	log.StdDebug(ctx, metadata, nil, message)
}

func removePII(message string) string {
	//to simulate remove PII high process time
	minSleep := 1
	maxSleep := 10
	sleepTime := rand.Intn(maxSleep-minSleep+1) + minSleep
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	return message
}

func Debug(ctx context.Context, message string) {
	q.append(item{
		ctx:     ctx,
		message: message,
	})
}
