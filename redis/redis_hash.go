package redis

import (
	"context"

	"github.com/sheenacodes/sharedutils/logger"
)

func (r *RedisClient) AddFieldValueToHash(hashKey string, fieldName string, fieldValue string) error {
	ctx := context.Background()
	err := r.Client.HSet(ctx, hashKey, fieldName, fieldValue).Err()
	if err != nil {
		logger.Log.Error().Err(err).Msgf("Error Setting Hash Field %s for key %s", fieldName, hashKey)
		return err
	}
	logger.Log.Debug().Msgf(" %s:%s added to Redis Hash for key %s", fieldName, fieldValue, hashKey)
	return nil
}
