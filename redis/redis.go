package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/sheenacodes/sharedutils/logger"

	"github.com/redis/go-redis/v9"
)

const (
	maxRetries     = 5                // Maximum number of retries before giving up
	initialBackoff = 2 * time.Second  // Initial delay before retrying
	maxBackoff     = 30 * time.Second // Maximum delay between retries
)

// RedisClient is a wrapper around the redis.Client to hold the instance
type RedisClient struct {
	Client *redis.Client
}

// Function to connect to Redis with retry and exponential backoff
func GetRedisClient(addr string, pword string, database int) (*RedisClient, error) {
	var client *redis.Client
	var err error

	for retries := 0; retries < maxRetries; retries++ {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pword,
			DB:       database,
		})
		logger.Log.Info().Msg(" connecting to Redis")
		ctx := context.Background()
		_, err = client.Ping(ctx).Result()
		if err == nil {
			// Successful connection
			logger.Log.Info().Msg("Successfully connected to Redis")
			return &RedisClient{Client: client}, nil
		} else {
			backoff := time.Duration((1 << retries) * int(initialBackoff))
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			logger.Log.Warn().Err(err).Msgf("Failed to connect to Redis, retrying in %v...", backoff)
			time.Sleep(backoff)
		}
	}

	//if err != nil {
	logger.Log.Fatal().Err(err).Msg("Failed to connect to Redis after multiple attempts")

	return nil, fmt.Errorf("failed to connect to Redis after %d attempts: %v", maxRetries, err)
	//}

}

func (r *RedisClient) IsSetNotEmpty(redisSetName string) (bool, error) {
	// Use the SCARD command to get the number of members in the set
	ctx := context.Background()
	card, err := r.Client.SCard(ctx, redisSetName).Result()
	if err != nil {
		logger.Log.Error().Err(err).Msg("error checking set size")
		return false, fmt.Errorf("error checking set size: %v", err)
	}

	logger.Log.Debug().Msgf("%d items available in Redis Set", card)
	// Return true if the number of members is greater than 0
	return card > 0, nil
}

func (r *RedisClient) AddItemToSet(item string, redisSetName string) error {
	ctx := context.Background()
	_, err := r.Client.SAdd(ctx, redisSetName, item).Result()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to make entry to Redis Set")
		return err
	}
	logger.Log.Debug().Msgf(" %s added to Redis Set %s", item, redisSetName)
	return nil
}

func (r *RedisClient) RemoveItemFromSet(item string, redisSetName string) error {
	ctx := context.Background()
	_, err := r.Client.SRem(ctx, redisSetName, item).Result()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to remove entry from Redis Set")
		return err
	}
	logger.Log.Debug().Msgf("%s removed from Redis Set %s", item, redisSetName)
	return nil
}

func (r *RedisClient) GetRandomItemFromSet(redisSetName string) (string, error) {
	// Use SRANDMEMBER to get a random member from the set
	ctx := context.Background()
	item, err := r.Client.SRandMember(ctx, redisSetName).Result()
	if err != nil {
		return "", fmt.Errorf("could not get random member from set: %w", err)
	}
	return item, nil
}
