package broker

import (
	"context"
	"github.com/google/uuid"
)

const (
	// This constant represents how many messages can be in the message queue of subscriber
	defaultMsgBufferSize = 50
)

// Subscriber represents a client who subscribed to a particular topic
type Subscriber struct {
	ID     string
	ctx    context.Context
	MsgC   chan Message
	unsubC chan *Subscriber
}

// NewSubscriber return a new subscriber instance
func NewSubscriber(ctx context.Context, unsubC chan *Subscriber) *Subscriber {
	sub := &Subscriber{
		ID:     uuid.New().String(), // TODO change to autoincrement id's
		MsgC:   make(chan Message, defaultMsgBufferSize),
		ctx:    ctx,
		unsubC: unsubC,
	}

	go sub.eventListener()
	return sub
}

func (s *Subscriber) eventListener() {
	for {
		select {
		case <-s.ctx.Done():
			s.unsubC <- s
			return
		}
	}
}
