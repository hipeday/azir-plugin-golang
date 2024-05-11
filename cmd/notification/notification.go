package notification

import (
	"github.com/hipeday/azir-plugin-golang/pkg/notify"
)

func init() {
	notify.Registry = *notify.NewNotificationRegistry()
	notify.Registry.Register(&notify.UnixNotification{})
}
