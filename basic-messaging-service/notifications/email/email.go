package email

import (
	"github.com/frost060/go-microservice-basic/basic-messaging-service/configs"
	notifications "github.com/frost060/go-microservice-basic/basic-messaging-service/notifications"
)

// Email Service Providers
const (
	SendGrid = 0
)

// Dispatcher => Dispatcher Factory For all Email dispatcher
func Dispatcher(sender int, to, subject, msg string, config *configs.ServerConfig) notifications.Dispatcher {
	switch sender {
	case SendGrid:
		return NewSendGridDispatcher(config.SendGrid.APIKey, to, subject, msg)
	default:
		return nil
	}
}
