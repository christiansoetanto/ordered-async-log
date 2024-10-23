package oalog

import (
	"context"
)

var q *queue

type queue struct {
	msgChan chan item
}
type item struct {
	ctx     context.Context
	message string
}

func processQueue() {
	go func() {
		for {
			i, ok := <-q.msgChan
			if !ok {
				break
			}
			execDebug(i.ctx, i.message)
		}
	}()
}

func Close() {
	close(q.msgChan)
}
