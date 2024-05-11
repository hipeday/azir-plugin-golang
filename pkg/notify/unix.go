package notify

import (
	"encoding/json"
	"github.com/hipeday/azir-plugin-golang/pkg/properties"
	"net"
)

type UnixNotification struct {
	LoggerNotification
}

func (n *UnixNotification) Type() properties.NotifyType {
	return properties.UNIX
}

func (n *UnixNotification) Push(message interface{}, target *properties.Notification) error {
	var (
		logger = n.GetLogger()
		err    error
	)

	conn, err := net.Dial("unix", target.Address)

	if err != nil {
		logger.Fatalf("Error dialing socket: %v", err)
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Fatalf("Error closing socket: %v", err)
		}
	}(conn)

	body, err := json.Marshal(message)

	if err != nil {
		return err
	}

	_, err = conn.Write(body)

	return err

}
