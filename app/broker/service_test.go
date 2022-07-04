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
	return broker.Message{Body: []byte(body)}
}

func TestShouldUnsubscribeAndNotReceiveMessage(t *testing.T) {
	b := broker.TestBroker(t)

	ctxSub, cancel := context.WithCancel(context.Background())
	sub1Chan, _ := b.Subscribe(ctxSub, "tt")

	cancel()
	_, _ = b.Publish(context.Background(), "tt", testMessage())

	select {
	case _ = <-sub1Chan:
		t.Fail()
	default:
	}
}

func TestShouldReceiveOneMessageBeforeUnsubscribing(t *testing.T) {
	b := broker.TestBroker(t)

	ctxSub, cancel := context.WithCancel(context.Background())
	sub1Chan, _ := b.Subscribe(ctxSub, "tt")

	msg := testMessage()
	id, _ := b.Publish(context.Background(), "tt", msg)
	msg.ID = id

	cancel()

	in := <-sub1Chan
	assert.Equal(t, msg, in)
}

func TestShouldPublishMessage(t *testing.T) {
	b := broker.TestBroker(t)
	msg := testMessage()

	id, err := b.Publish(context.Background(), "tt", msg)
	assert.NoError(t, err)
	assert.NotNil(t, id)
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
