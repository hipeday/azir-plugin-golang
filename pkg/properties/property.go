package properties

import (
	"github.com/hipeday/azir-plugin-golang/pkg/logging"
)

type NotifyType string

const (
	UNIX NotifyType = "UNIX"
)

type Property struct {
	InvokeId     string         `json:"invoke_id"`
	Logger       logging.Config `json:"logger"`
	Name         string         `json:"name"`
	Version      string         `json:"version"`
	Home         string         `json:"home"`
	Notification *Notification  `json:"notification"`
	Description  string         `json:"description"`
	Logfile      string         `json:"logfile"`
	Cmd          CmdProperty    `json:"cmd"`
}

type CmdProperty struct {
	Linux   []string `json:"linux"`
	Windows []string `json:"windows"`
	Darwin  []string `json:"darwin"`
}

type Notification struct {
	// Type is the type of notification default is "UNIX"
	Type NotifyType `json:"type"`
	// Address is the address of the notification
	Address string `json:"address"`
	// Enabled is the flag to enable or disable the notification
	Enabled bool `json:"enabled"`
}
