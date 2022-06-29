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
	ID     string // TODO change to autoincrement id's
	ctx    context.Context
	MsgC   chan Message
	unsubC chan *Subscriber
}

func NewSubscriber(ctx context.Context, unsubC chan *Subscriber) *Subscriber {
	sub := &Subscriber{
		ID:     uuid.New().String(),
		MsgC:   make(chan Message, defaultMsgBufferSize),
		ctx:    ctx,
		unsubC: unsubC,
	}

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
