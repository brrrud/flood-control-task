package floodControl

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type RedisImpl struct {
	client   *redis.Client
	K        int
	N        time.Duration
	redisKey string
}

func NewFloodControlRedisImpl(client *redis.Client, limit int, duration time.Duration, redisKey string) *RedisImpl {
	return &RedisImpl{client: client, K: limit, N: duration, redisKey: redisKey}
}

func (fc *RedisImpl) Check(ctx context.Context, userID int64) (bool, error) {
	key := fc.redisKey + ":" + strconv.FormatInt(userID, 10)

	count, err := fc.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		if err := fc.client.Expire(ctx, key, fc.N).Err(); err != nil {
			return false, err
		}
	} else {
		if err := fc.client.Expire(ctx, key, fc.N).Err(); err != nil {
			return false, err
		}
	}

	if count > int64(fc.K) {
		if err := fc.client.Del(ctx, key).Err(); err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}
