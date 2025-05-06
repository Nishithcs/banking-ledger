package pkg

type MessageQueue interface {
	Connect(url string) error
	Publish(queue string, body []byte) error
	Consume(queue string) (<-chan []byte, error)
	Close() error
}


// type MessageQueue interface {
// 	Connect(url string) error
// 	Publish(queue string, body []byte) error
// 	Consume(queue string) (<-chan amqp.Delivery, error)
// 	Close() error
// }