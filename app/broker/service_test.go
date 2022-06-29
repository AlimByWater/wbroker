package broker_test

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	grpcOrig "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"wbroker/app/broker"
	wb "wbroker/external/grpc/wbroker"
)

func testMessage() broker.Message {
	body := fmt.Sprintf("test_%s", uuid.New().String())
	return broker.Message{Body: body}
}

func testGRPCClient(t *testing.T) wb.WBrokerClient {
	conn, err := grpcOrig.Dial("localhost:8018", grpcOrig.WithTransportCredentials(insecure.NewCredentials()), grpcOrig.WithBlock())
	if err != nil {
		t.Fatalf("couldnt connect: %s", err)
	}
	defer conn.Close()

	return wb.NewWBrokerClient(conn)
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
