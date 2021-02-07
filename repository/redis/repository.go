package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/paoloposso/poc/hex-ms/shortener"
	"github.com/pkg/errors"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisUrl string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	_, err = client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewRedisRepository(redisUrl string) (shortener.RedirectRepository, error) {
	repo := &redisRepository{}
	client, err := newRedisClient(redisUrl)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}
	repo.client = client
	return repo, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepository) Find(code string) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	redirect.Code = data["code"]
	redirect.Url = data["url"]
	redirect.CreatedAt = createdAt
	return redirect, nil
}

func (r *redisRepository) Store(redirect *shortener.Redirect) error {
	data := map[string]interface{}{
		"code": redirect.Code,
		"url": redirect.Url,
		"created_at": redirect.CreatedAt,
	}
	_, err := r.client.HSet(context.TODO(), r.generateKey(redirect.Code), data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Find")
	}
	return nil
}