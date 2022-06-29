package broker

import "testing"

func TestBroker(t *testing.T) *ActualBroker {
	return NewBroker()
}
