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

func TestIntegration_Broker(t *testing.T) {
	b := broker.TestBroker(t)
	msg := testMessage()

	sub1, _ := b.Subscribe(context.Background(), "tt")
	sub2, _ := b.Subscribe(context.Background(), "tt")

	_, _ = b.Publish(context.Background(), "tt", msg)

	respMsg1 := <-sub1
	respMsg2 := <-sub2

	assert.Equal(t, msg, respMsg1)
	assert.Equal(t, msg, respMsg2)
}

func TestUUID(t *testing.T) {
	id := uuid.New().ID()
	t.Log(id)
}
