package broker

import (
	"context"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
)

var (
	ErrFailedToPublishMessage error = fmt.Errorf("failed to publish message")
)

type (

	// Service is a broker service interface
	Service interface {
		Subscribe(ctx context.Context, topicName string) (chan Message, error)
		Publish(ctx context.Context, topicName string, message Message) (string, error)
	}

	ActualBroker struct {
		topics *TopicStorage
	}

	Message struct {
		ID   string
		Body []byte
	}
)

// NewBroker return ActualBroker
func NewBroker() *ActualBroker {
	return &ActualBroker{
		topics: NewTopicStorage(),
	}
}

// Subscribe subscribes a client to a particular topic
func (b *ActualBroker) Subscribe(ctx context.Context, topicName string) (chan Message, error) {
	topic, exist := b.topics.Topic(topicName)
	if !exist {
		topic = b.topics.CreateTopic(topicName)
	}

	msgChannel := topic.RegisterNewSubscriber(ctx)
	return msgChannel, nil
}

// Publish publishes a message to a particular topic
func (b *ActualBroker) Publish(ctx context.Context, topicName string, message Message) (string, error) {
	topic, exist := b.topics.Topic(topicName)
	if !exist {
		topic = b.topics.CreateTopic(topicName)
	}

	id, err := topic.Publish(message)
	if err != nil {
		return "", ErrFailedToPublishMessage
	}

	return id, nil
}

// Validate itself and returns error
func (m *Message) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Body, validation.Required),
	)
}
