package handlers_test

import (
	"fmt"
)

// MockMessageQueue is a mock implementation of the pkg.MessageQueue interface
type MockMessageQueue struct {
	PublishedMessages [][]byte
	ShouldFail        bool
}

func (m *MockMessageQueue) Connect(url string) error {
	return nil
}

func (m *MockMessageQueue) Publish(queue string, body []byte) error {
	if m.ShouldFail {
		return fmt.Errorf("mock publish error")
	}
	m.PublishedMessages = append(m.PublishedMessages, body)
	return nil
}

func (m *MockMessageQueue) Consume(queue string) (<-chan []byte, error) {
	ch := make(chan []byte)
	close(ch)
	return ch, nil
}

func (m *MockMessageQueue) Close() error {
	return nil
}