package redis

import (
	"context"
	"fmt"
	"github.com/SuperMatch/zapLogger"
	"go.uber.org/zap"
	"strconv"
	"time"

	Redis "github.com/redis/go-redis/v9"
)

const TTL int = 2 * 24 * 60 * 60

//go:generate mockgen -package mocks -destination mocks/cache_mock.go github.com/SuperMatch/pkg/redis LikeDislikeCacheInterface

type LikeDislikeCacheInterface interface {
	GetLikeDislike(key string) (*string, error)
	PutLikeDislike(key string, value string) error
	GetMatchList(key string) ([]int, error)
	AddToUserMatchList(key string, value string) error
	RemoveFromUserMatchList(key string) error
	GetUserLikes(key string) ([]int, error)
	PutLiker(key, value string) error
	RemoveLikerFromLikeeList(key, value string) error
}

type LikeDislikeCache struct {
	redisClient *Redis.Client
}

func LikeDislikeCacheConstructor() *LikeDislikeCache {
	return &LikeDislikeCache{
		redisClient: RedisClient,
	}
}

func (l *LikeDislikeCache) GetLikeDislike(key string) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := l.redisClient.Get(ctx, key).Result()

	if err == Redis.Nil {
		return nil, nil
	}

	if err != nil {
		zapLogger.Logger.Error("Error while getting like/dislike from redis cache:", zap.Error(err))
		return nil, err
	}

	return &val, nil
}

func (l *LikeDislikeCache) PutLikeDislike(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := l.redisClient.Set(ctx, key, value, 48*time.Hour).Err()

	zapLogger.Logger.Error("error in setting PutLikeDislike:", zap.Error(err))

	if err != nil {
		return err
	}
	return nil
}

func (l *LikeDislikeCache) AddToUserMatchList(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := l.redisClient.RPush(ctx, key, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (l *LikeDislikeCache) RemoveFromUserMatchList(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := l.redisClient.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (l *LikeDislikeCache) GetMatchList(key string) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	val, err := l.redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	matches := make([]int, 0)
	for _, val := range val {
		c, err := strconv.Atoi(val)
		if err != nil {
			zapLogger.Logger.Error(fmt.Sprintf("error while converting string: %s to int: %d ", val, err))
		}
		matches = append(matches, c)
	}
	return matches, nil
}

func (l *LikeDislikeCache) GetUserLikes(key string) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	val, err := l.redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	likes := make([]int, 0)
	for _, val := range val {
		l, err := strconv.Atoi(val)
		if err != nil {
			zapLogger.Logger.Error(fmt.Sprintf("error while converting string: %s to int: %d ", val, err))
		}
		likes = append(likes, l)
	}
	return likes, nil

}

func (l *LikeDislikeCache) PutLiker(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := l.redisClient.LPush(ctx, key, value).Err()

	zapLogger.Logger.Error("error in setting Liker to likee's likers list", zap.Error(err))

	if err != nil {
		return err
	}
	return nil
}

func (l *LikeDislikeCache) RemoveLikerFromLikeeList(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := l.redisClient.LRem(ctx, key, 0, value).Err()

	zapLogger.Logger.Error("error in removing Liker from likee's likers list", zap.Error(err))

	if err != nil {
		return err
	}
	return nil
}
