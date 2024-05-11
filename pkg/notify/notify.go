package notify

import (
	"github.com/hipeday/azir-plugin-golang/pkg/properties"
	"go.uber.org/zap"
)

var (
	Registry NotificationRegistry
)

type Notification interface {
	Type() properties.NotifyType
	Push(message interface{}, target *properties.Notification) error
	GetLogger() *zap.SugaredLogger
	SetLogger(logger *zap.SugaredLogger)
}

type NotificationRegistry struct {
	notifications map[string]Notification
}

func NewNotificationRegistry() *NotificationRegistry {
	return &NotificationRegistry{
		notifications: make(map[string]Notification),
	}
}

func (r *NotificationRegistry) Register(notification Notification) {
	r.notifications[string(notification.Type())] = notification
}

func (r *NotificationRegistry) GetNotification(notificationType properties.NotifyType) Notification {
	return r.notifications[string(notificationType)]
}

type LoggerNotification struct {
	logger *zap.SugaredLogger
}

func (l *LoggerNotification) SetLogger(logger *zap.SugaredLogger) {
	l.logger = logger
}

func (l *LoggerNotification) GetLogger() *zap.SugaredLogger {
	return l.logger
}
