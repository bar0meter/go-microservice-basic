package db

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"

	"github.com/golang/protobuf/proto"

	"github.com/frost060/go-microservice-basic/basic-messaging-service/configs"
	protos "github.com/frost060/go-microservice-basic/basic-messaging-service/protos/notifications"
)

type Redis struct {
	client *redis.Client
}

func NewRedisClient(serverConfig *configs.ServerConfig) *Redis {
	opts := &redis.Options{
		Addr:     serverConfig.Redis.Addr,
		Password: serverConfig.Redis.Password,
		DB:       serverConfig.Redis.DB,
	}

	client := redis.NewClient(opts)

	return &Redis{client}
}

func (rc *Redis) Push(ctx context.Context, key string, message *protos.MessageRequest) (bool, error) {
	value := proto.MarshalTextString(message)
	result := rc.client.LPush(ctx, key, value)

	if result.Err() != nil {
		return false, result.Err()
	} else if result.Val() == 0 {
		return false, errors.New("invalid key")
	}

	return true, nil
}

func (rc *Redis) Pop(ctx context.Context, key string) (*protos.MessageRequest, error) {
	result := rc.client.RPop(ctx, key)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var message protos.MessageRequest
	err := proto.UnmarshalText(result.Val(), &message)
	if err != nil {
		// Error occurred while marshalling then push ti back to redis
		// TODO: Here check for error type. If message is bad then need to discard the message.
		rc.client.LPush(ctx, key, result.Val())
		return nil, err
	}

	return &message, nil
}
