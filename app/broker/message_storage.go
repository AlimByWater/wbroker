package broker

import (
	"github.com/google/uuid"
	"sync"
)

// TODO make message storage as eventlog file

// TODO add time-to-live parameter

// MessageStorage is in-memory storage for messages
type MessageStorage struct {
	messages map[string]Message
	mu       sync.RWMutex
}

// NewMessageStorage return new MessageStorage
func NewMessageStorage() *MessageStorage {
	return &MessageStorage{
		messages: make(map[string]Message),
	}
}

// RegisterNewMessage adds message to a message storage
func (ms *MessageStorage) RegisterNewMessage(message Message) Message {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	message.ID = uuid.New().String()
	ms.messages[message.ID] = message

	return message
}
