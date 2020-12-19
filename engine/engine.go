package engine

import (
	"sync"
)

// Command epresents actions that can be performed in a single event loop iteration.
type Command interface {
	Execute(h Handler)
}

// Handler allows to send commands to an event loop it's associated with.
type Handler interface {
	Post(cmd Command)
}

type messageQueue struct {
	sync.Mutex

	data               []Command
	receiveSignal      chan struct{}
	isReceiveRequested bool
}

func (mq *messageQueue) push(command Command) {
	mq.Lock()
	defer mq.Unlock()

	mq.data = append(mq.data, command)
	if mq.isReceiveRequested {
		mq.isReceiveRequested = false
		mq.receiveSignal <- struct{}{}
	}
}

func (mq *messageQueue) isEmpty() bool {
	return len(mq.data) == 0
}

func (mq *messageQueue) pull() Command {
	mq.Lock()
	defer mq.Unlock()

	if mq.isEmpty() {
		mq.isReceiveRequested = true
		mq.Unlock()
		<-mq.receiveSignal
		mq.Lock()
	}

	res := mq.data[0]
	mq.data[0] = nil
	mq.data = mq.data[1:]
	return res
}

// EventLoop pattern
type EventLoop struct {
	mq              *messageQueue
	stopSignal      chan struct{}
	isStopRequested bool
}

// Start initialized message queue
func (el *EventLoop) Start() {
	el.mq = &messageQueue{receiveSignal: make(chan struct{})}
	el.stopSignal = make(chan struct{})

	go func() {
		for !el.isStopRequested || !el.mq.isEmpty() {
			cmd := el.mq.pull()
			cmd.Execute(el)
		}
		el.stopSignal <- struct{}{}
	}()
}

type commangFunc func(h Handler)

func (cf commangFunc) Execute(h Handler) {
	cf(h)
}

// AwaitFinish waits untill commands in queue are finished
func (el *EventLoop) AwaitFinish() {
	el.Post(commangFunc(func(h Handler) {
		el.isStopRequested = true
	}))
	<-el.stopSignal
}

// Post adds commant to queue
func (el *EventLoop) Post(cmd Command) {
	el.mq.push(cmd)
}
