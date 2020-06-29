package server

import (
	"context"
	"errors"

	"github.com/frost060/go-microservice-basic/basic-messaging-service/configs"
	"github.com/frost060/go-microservice-basic/basic-messaging-service/notifications"
	"github.com/frost060/go-microservice-basic/basic-messaging-service/notifications/email"
	protos "github.com/frost060/go-microservice-basic/basic-messaging-service/protos/notifications"
	"github.com/hashicorp/go-hclog"
)

// MessageService => Sends Notificaiotns
type MessageService struct {
	log    hclog.Logger
	config *configs.ServerConfig
}

// NewMessageService => returns a new message service
func NewMessageService(l hclog.Logger, config *configs.ServerConfig) *MessageService {
	return &MessageService{l, config}
}

// SendNotification => Sends a notification without processing (dont add to queue)
// Used for forgot password, verify account, login OTP, etc.
func (ms *MessageService) SendNotification(ctx context.Context, req *protos.MessageRequest) (*protos.MessageResponse, error) {
	messageType := req.GetType()

	var dispatcher notifications.Dispatcher
	msg := req.GetMsg()
	to := req.GetTo()
	subject := req.GetSubject()

	switch messageType {
	case protos.NotificationType_EMAIL:
		dispatcher = email.Dispatcher(email.SendGrid, to, subject, msg, ms.config)
	default:
		dispatcher = nil
	}

	if dispatcher == nil {
		return &protos.MessageResponse{
			Success: false,
		}, errors.New("invalid message type")
	}

	success, err := dispatcher.Dispatch()
	return &protos.MessageResponse{
		Success: success,
	}, err
}
