package rabbitmq

import (
	"encoding/json"
	"time"

	"github.com/sheenacodes/sharedutils/logger"

	"github.com/rabbitmq/amqp091-go"
)

const (
	maxRetries     = 5                // Maximum number of retries before giving up
	initialBackoff = 2 * time.Second  // Initial delay before retrying
	maxBackoff     = 30 * time.Second // Maximum delay between retries
)

// RabbitMQClient holds the RabbitMQ connection and channel
type RabbitMQClient struct {
	Connection *amqp091.Connection
	Channel    *amqp091.Channel
}

// NewRabbitMQClient creates a new RabbitMQ client with retry logic
func GetRabbitMQClient(url string) (*RabbitMQClient, error) {
	var conn *amqp091.Connection
	var ch *amqp091.Channel
	var err error

	for retries := 0; retries < maxRetries; retries++ {
		conn, err = amqp091.Dial(url)
		if err != nil {
			backoff := time.Duration((1 << retries) * int(initialBackoff))
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			logger.Log.Warn().Err(err).Msgf("Failed to connect to RabbitMQ, retrying in %v...", backoff)
			time.Sleep(backoff)
		} else {
			logger.Log.Info().Msg("Successfully connected to RabbitMQ")
			break
		}
	}

	if err != nil {
		return nil, err
	}

	// Create a new channel
	ch, err = conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQClient{
		Connection: conn,
		Channel:    ch,
	}, nil
}

// Close cleans up the RabbitMQ connection and channel
func (r *RabbitMQClient) Close() {
	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to close RabbitMQ channel")
		}
	}
	if r.Connection != nil {
		if err := r.Connection.Close(); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to close RabbitMQ connection")
		}
	}
}
func PublishEvent(conn *amqp091.Connection, queueName string, eventPayload any) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	body, err := json.Marshal(eventPayload)
	if err != nil {
		logger.Log.Fatal().Err(err).Msgf("JSON conversion error in %v...", eventPayload)
		return err
	}

	err = channel.Publish(
		"",        // exchange
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	logger.Log.Info().Msgf("Published event: %s", body)
	return nil
}
