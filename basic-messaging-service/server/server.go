package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/frost060/go-microservice-basic/basic-messaging-service/db"
	"github.com/frost060/go-microservice-basic/basic-messaging-service/logging"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/frost060/go-microservice-basic/basic-messaging-service/configs"
	"github.com/frost060/go-microservice-basic/basic-messaging-service/notifications"
	"github.com/frost060/go-microservice-basic/basic-messaging-service/notifications/email"
	protos "github.com/frost060/go-microservice-basic/basic-messaging-service/protos/notifications"
)

// MessageService => Sends Notificaiotns
type MessageService struct {
	config *configs.ServerConfig
	Redis  *db.Redis
	log    *logging.LogWrapper
}

// NewMessageService => returns a new message service
func NewMessageService(config *configs.ServerConfig, redis *db.Redis, l *logging.LogWrapper) *MessageService {
	return &MessageService{config, redis, l}
}

// SendNotification => Sends a notification without processing (dont add to queue)
// Used for forgot password, verify account, login OTP, etc.
func (ms *MessageService) SendNotification(
	ctx context.Context, req *protos.MessageRequest) (*protos.MessageResponse, error) {
	messageType := req.GetType()

	var dispatcher notifications.Dispatcher
	msg := req.GetMsg()
	to := req.GetTo()
	subject := req.GetSubject()

	switch messageType {
	case protos.NotificationType_EMAIL:
		dispatcher = email.Dispatcher(
			email.GetProvider(ms.config.Providers.Email), to, subject, msg, ms.config)
	default:
		dispatcher = nil
	}

	if dispatcher == nil {
		return &protos.MessageResponse{
			Success: false,
		}, errors.New("invalid message type")
	}

	success, err := dispatcher.Dispatch()
	ms.log.Info(fmt.Sprintf("Success: %v, Error: %v", success, err))

	return &protos.MessageResponse{
		Success: success,
	}, err
}

func (ms *MessageService) AddToQueue(
	ctx context.Context, req *protos.MessageRequest) (*protos.MessageResponse, error) {

	ok, err := ms.Redis.Push(ctx, "default", req)

	return &protos.MessageResponse{
		Success: ok,
	}, err
}

func (ms *MessageService) RemoveFromQueue(ctx context.Context, _ *empty.Empty) (*protos.MessageRequest, error) {
	return ms.Redis.Pop(ctx, "default")
}

type Worker struct {
	ID int
}

type WorkerPool struct {
	Pool chan Worker
}

func (ms *MessageService) newWorkerPool(noOfRoutines int) *WorkerPool {
	workerPool := make(chan Worker, noOfRoutines)
	for i := 0; i < noOfRoutines; i++ {
		worker := &Worker{
			ID: i,
		}
		ms.log.Info("Created worker with id %d", i)
		workerPool <- *worker
	}

	return &WorkerPool{
		Pool: workerPool,
	}
}

// https://play.golang.org/p/HovNRgp6FxH
func (ms *MessageService) StartDispatchRedis(noOfRoutines int, redis *db.Redis) {
	workerPool := ms.newWorkerPool(noOfRoutines)
	for {
		worker := <-workerPool.Pool
		go func() {
			ctx := context.Background()
			message, err := redis.Pop(ctx, "default")
			if err != nil {
				time.Sleep(1 * time.Minute)
				workerPool.Pool <- worker
				return
			}

			resp, err := ms.SendNotification(ctx, message)
			if err != nil || !resp.Success {
				ms.log.Error("Error occurred while dispatching message, pushing back to redis")
				_, _ = redis.Push(ctx, "default", message)
			} else {
				ms.log.Info("Successfully sent message, by worker: %d", worker.ID)
			}

			workerPool.Pool <- worker
		}()
	}
}
