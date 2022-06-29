package broker

import (
	"context"
	"sync"
)

// Topic contains all its subscribers
type Topic struct {
	Name            string
	Subscribers     map[string]*Subscriber
	MessagesStorage *MessageStorage
	msgPubC         chan Message
	deleteSubC      chan *Subscriber
	mu              sync.RWMutex
}

type TopicStorage struct {
	topics map[string]*Topic
	mu     sync.RWMutex
}

func NewTopicStorage() *TopicStorage {
	return &TopicStorage{
		topics: make(map[string]*Topic),
	}
}

func NewTopic(name string) *Topic {
	topic := &Topic{
		Name:            name,
		Subscribers:     make(map[string]*Subscriber),
		MessagesStorage: NewMessageStorage(),
		msgPubC:         make(chan Message),
		deleteSubC:      make(chan *Subscriber),
	}

	go topic.eventListener()
	return topic
}

func (t *Topic) RegisterNewSubscriber(ctx context.Context) <-chan Message {
	sub := NewSubscriber(ctx, t.deleteSubC)

	t.mu.Lock()
	defer t.mu.Unlock()
	t.Subscribers[sub.ID] = sub

	return sub.MsgC
}

func (t *Topic) Publish(ctx context.Context, topicName string, message Message) (string, error) {
	message = t.MessagesStorage.RegisterNewMessage(message)

	t.msgPubC <- message
	return message.ID, nil
}

func (t *Topic) eventListener() {
	for {
		select {
		case msg := <-t.msgPubC:
			var wg sync.WaitGroup
			wg.Add(len(t.Subscribers))

			for _, sub := range t.Subscribers {
				go func(s *Subscriber) {
					sub.MsgC <- msg
					wg.Done()
				}(sub)
			}
		case sub := <-t.deleteSubC:
			t.mu.Lock()
			delete(t.Subscribers, sub.ID)
			t.mu.Unlock()
		}

	}
}

func (ts *TopicStorage) Topic(name string) (*Topic, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	topic, exist := ts.topics[name]
	if !exist {
		return nil, false
	}

	return topic, true
}

func (ts *TopicStorage) CreateTopic(name string) *Topic {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	topic := NewTopic(name)
	ts.topics[topic.Name] = topic
	return topic
}
