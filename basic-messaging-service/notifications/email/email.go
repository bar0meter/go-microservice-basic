package email

import (
	"github.com/frost060/go-microservice-basic/basic-messaging-service/configs"
	"github.com/frost060/go-microservice-basic/basic-messaging-service/notifications"
)

// Email Service Providers Ordinal
const (
	SendGrid = 0
)

// Email Service Providers
const (
	SENDGRID = "sendgrid"
)

// Dispatcher => Dispatcher Factory For all Email dispatcher
func Dispatcher(sender int, to, subject, msg string, config *configs.ServerConfig) notifications.Dispatcher {
	switch sender {
	case SendGrid:
		return NewSendGridDispatcher(to, subject, msg, config.SendGrid.APIKey)
	default:
		return nil
	}
}

// GetProvider => Returns provider ordinal
func GetProvider(sender string) int {
	switch sender {
	case SENDGRID:
		return SendGrid
	default:
		return -1
	}
}
