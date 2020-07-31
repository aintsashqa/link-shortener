package redis

import (
	"fmt"
	"strconv"

	"github.com/aintsashqa/link-shortener/src/shortener"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(rURL string) (*redis.Client, error) {
	options, err := redis.ParseURL(rURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(options)
	_, err = client.Ping().Result()
	return client, err
}

func NewRedisRepository(rURL string) (shortener.RedirectRepositoryInterface, error) {
	repository := &redisRepository{}
	client, err := newRedisClient(rURL)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}
	repository.client = client
	return repository, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepository) Find(code string) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(key).Result()
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
	redirect.Link = data["link"]
	redirect.CreatedAt = createdAt
	return redirect, nil
}

func (r *redisRepository) Add(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"link":       redirect.Link,
		"created_at": redirect.CreatedAt,
	}
	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Add")
	}
	return nil
}
