package broker

import (
	"context"
	"fmt"
)

var (
	ErrTopicDoesntExist error = fmt.Errorf("topic doesn't exist")
)

type (

	// Service is a broker service interface
	Service interface {
		Subscribe(ctx context.Context, topicName string) (<-chan Message, error)
		Publish(ctx context.Context, topicName string, message Message) (string, error)
	}

	ActualBroker struct {
		topics *TopicStorage
	}

	Message struct {
		ID   string
		Body string
	}
)

// NewBroker return ActualBroker
func NewBroker() *ActualBroker {
	return &ActualBroker{
		topics: NewTopicStorage(),
	}
}

func (b *ActualBroker) Subscribe(ctx context.Context, topicName string) (<-chan Message, error) {
	topic, exist := b.topics.Topic(topicName)
	if !exist {
		topic = b.topics.CreateTopic(topicName)
	}

	msgChannel := topic.RegisterNewSubscriber(ctx)
	return msgChannel, nil
}

func (b *ActualBroker) Publish(ctx context.Context, topicName string, message Message) (string, error) {
	topic, exist := b.topics.Topic(topicName)
	if !exist {
		topic = b.topics.CreateTopic(topicName)
	}

	id, err := topic.Publish(ctx, topicName, message)
	if err != nil {
		return "", fmt.Errorf("failed to publish message")
	}

	return id, nil
}
