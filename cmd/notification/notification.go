package notification

import (
	"github.com/ideal-rucksack/workflow-glolang-plugin/pkg/notify"
)

func init() {
	notify.Registry = *notify.NewNotificationRegistry()
	notify.Registry.Register(&notify.UnixNotification{})
}
