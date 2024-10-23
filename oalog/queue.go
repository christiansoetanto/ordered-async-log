package oalog

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var q *queue

type queue struct {
	head *logData
	tail *logData
	s    int
	mtx  sync.Mutex
}
type logData struct {
	i    *item
	next *logData
}
type item struct {
	ctx     context.Context
	message string
}

func (q *queue) pop() (*item, error) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	if q == nil || q.head == nil {
		return nil, errors.New("nil queue")
	}
	i := q.head.i
	q.head = q.head.next
	q.s--
	return i, nil
}

func (q *queue) append(i item) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	n := logData{
		i: &i,
	}
	if q == nil {
		panic("must init first")
	}
	if q.head == nil {
		q.head = &n
		q.tail = &n
	} else {
		q.tail.next = &n
		q.tail = &n
	}
	q.s++
}
func (q *queue) size() int {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	return q.s
}

func processQueue() {
	go func() {
		for {
			i, err := q.pop()
			if err == nil {
				execDebug(i.ctx, i.message)
			}
		}
	}()
}

func WaitForEmptyQueue(ch chan struct{}) {
	fmt.Println("waiting for q to empty.")
	for {
		if q.size() == 0 {
			ch <- struct{}{}
			return
		}
	}
}
