package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/frost060/go-microservice-basic/basic-messaging-service/db"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/frost060/go-microservice-basic/basic-messaging-service/configs"
	log "github.com/frost060/go-microservice-basic/basic-messaging-service/logging"
	"github.com/frost060/go-microservice-basic/basic-messaging-service/notifications"
	"github.com/frost060/go-microservice-basic/basic-messaging-service/notifications/email"
	protos "github.com/frost060/go-microservice-basic/basic-messaging-service/protos/notifications"
)

// MessageService => Sends Notificaiotns
type MessageService struct {
	config *configs.ServerConfig
	Redis  *db.Redis
}

// NewMessageService => returns a new message service
func NewMessageService(config *configs.ServerConfig, redis *db.Redis) *MessageService {
	return &MessageService{config, redis}
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
	log.Info(fmt.Sprintf("Success: %v, Error: %v", success, err))

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
