package main

import (
	"context"
	wb "github.com/AlimByWater/wbroker/external/grpc/wbroker"
	"github.com/sirupsen/logrus"
)

func main() {
	client, err := wb.ProvideClients()
	if err != nil {
		logrus.Fatalf("failed to initialize client: %s", err.Error())
	}

	ctx := context.Background()
	stream, err := client.WBrokerClient.Subscribe(ctx, &wb.SubscribeRequest{Topic: "new_order"})
	if err != nil {
		logrus.Fatalf("failed to subscribe: %s", err.Error())
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			logrus.Errorf("stream: %s", err.Error())
		}

		logrus.Infof("message: %s", string(msg.Body))
	}
}
