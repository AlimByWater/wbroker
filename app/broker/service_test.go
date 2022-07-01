package broker_test

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"wbroker/app/broker"
)

func testMessage() broker.Message {
	body := fmt.Sprintf("test_%s", uuid.New().String())
	return broker.Message{Body: body}
}

func TestTwoSubsShouldRecieveMessage(t *testing.T) {
	b := broker.TestBroker(t)
	msg := testMessage()

	sub1Chan, _ := b.Subscribe(context.Background(), "tt")
	sub2Chan, _ := b.Subscribe(context.Background(), "tt")

	id, _ := b.Publish(context.Background(), "tt", msg)
	msg.ID = id

	for i := 0; i < 2; i++ {
		select {
		case respMsg1 := <-sub1Chan:
			assert.Equal(t, msg, respMsg1)
		case respMsg2 := <-sub2Chan:
			assert.Equal(t, msg, respMsg2)
		}
	}
}
