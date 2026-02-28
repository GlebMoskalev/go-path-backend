package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type StateRepository interface {
	Save(ctx context.Context, state string, ttl time.Duration) error
	Validate(ctx context.Context, state string) (bool, error)
}

type stateRepository struct {
	client *redis.Client
}

func NewStateRepository(client *redis.Client) StateRepository {
	return &stateRepository{client: client}
}

func (r *stateRepository) Save(ctx context.Context, state string, ttl time.Duration) error {
	key := fmt.Sprintf("oauth_state:%s", state)
	return r.client.Set(ctx, key, "1", ttl).Err()
}

func (r *stateRepository) Validate(ctx context.Context, state string) (bool, error) {
	key := fmt.Sprintf("oauth_state:%s", state)

	result, err := r.client.GetDel(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return result == "1", nil
}
