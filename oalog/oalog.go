package oalog

import (
	"context"
	"fmt"
	"github.com/tokopedia/tdk/v2/go/log"
	"math/rand"
	"time"
)

func execDebug(ctx context.Context, message string) {

	message = removePII(message)
	metadata := Metadata{
		string(contextKeyProcessNo): GetCtxProcessNo(ctx),
	}
	fmt.Println("exec debug ", message)
	log.StdDebug(ctx, metadata, nil, message)
}

func removePII(message string) string {
	//to simulate remove PII high process time
	minSleep := 1
	maxSleep := 10
	if message == "999" {
		minSleep = 5000
		maxSleep = 9000
	}
	sleepTime := rand.Intn(maxSleep-minSleep+1) + minSleep
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	return message
}

func Debug(ctx context.Context, message string) {
	q.msgChan <- item{
		ctx:     ctx,
		message: message,
	}
}
