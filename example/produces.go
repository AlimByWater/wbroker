package main

import (
	"context"
	"flag"
	wb "github.com/AlimByWater/wbroker/external/grpc/wbroker"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	flagTopic = flag.String("topic", "new_order", "")
	flagMsg   = flag.String("msg", randomMessage(), "")
)

func main() {
	flag.Parse()

	client, err := wb.ProvideClients()
	if err != nil {
		logrus.Fatalf("failed to initialize client: %s", err.Error())
	}

	ctx := context.Background()
	resp, err := client.WBrokerClient.Publish(ctx, &wb.PublishRequest{
		Topic: *flagTopic,
		Body:  []byte(*flagMsg),
	})

	if err != nil {
		logrus.Fatalf("failed to send message: %s", err.Error())
		return
	}

	logrus.Infof("message id: %s\t message body: %s", resp.Id, *flagMsg)
}

func randomMessage() string {
	return uuid.New().String()
}
