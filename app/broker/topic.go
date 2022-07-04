package broker

import (
	"context"
	"sync"
)

// Topic contains all its subscribers and messages
// Topic also manages publication of messages to its respective subscribers
type Topic struct {
	Name            string
	Subscribers     map[string]*Subscriber
	MessagesStorage *MessageStorage
	msgPubC         chan Message
	deleteSubC      chan *Subscriber
	mu              sync.RWMutex
}

// TopicStorage is a storage of all existing topics
type TopicStorage struct {
	topics map[string]*Topic
	mu     sync.RWMutex
}

// NewTopicStorage creates new TopicStorage
func NewTopicStorage() *TopicStorage {
	return &TopicStorage{
		topics: make(map[string]*Topic),
	}
}

// NewTopic creates new Topic
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

// RegisterNewSubscriber registers new subscriber, assigns it to the topic and returns message channel
func (t *Topic) RegisterNewSubscriber(ctx context.Context) chan Message {
	sub := NewSubscriber(ctx, t.deleteSubC)

	t.Subscribers[sub.ID] = sub

	return sub.MsgC
}

// Publish registers new message and publishes it to all subscribers
func (t *Topic) Publish(message Message) (string, error) {
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
					defer wg.Done()
					sub.MsgC <- msg
				}(sub)
			}
			wg.Wait()
		case sub := <-t.deleteSubC:
			t.mu.Lock()
			close(sub.MsgC) // we should close msg channel when user unsubscribes
			delete(t.Subscribers, sub.ID)
			t.mu.Unlock()
		}

	}
}

// Topic return topic from TopicStorage by name
func (ts *TopicStorage) Topic(name string) (*Topic, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	topic, exist := ts.topics[name]
	if !exist {
		return nil, false
	}

	return topic, true
}

// CreateTopic creates new topics
func (ts *TopicStorage) CreateTopic(name string) *Topic {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	topic := NewTopic(name)
	ts.topics[topic.Name] = topic
	return topic
}
