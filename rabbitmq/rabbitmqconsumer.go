package rabbitmq

type Client interface {
	ProcessMessage(msg []byte) error
}

// ConsumeQueue consumes messages from the specified RabbitMQ queue and uses the provided handler.
func (client *RabbitMQClient) ConsumeQueue(queueName string, handler Client) error {
	msgs, err := client.Channel.Consume(
		queueName, // Queue
		"",        // Consumer
		true,      // Auto-ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			//messageContent := string(msg.Body)
			//TODO catch error here
			if err := handler.ProcessMessage(msg.Body); err != nil {
				// Handle error
			}
		}
	}()

	return nil
}
