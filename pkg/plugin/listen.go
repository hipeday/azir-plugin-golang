package plugin

import (
	"encoding/json"
	"fmt"
	"github.com/hipeday/azir-plugin-golang/pkg/notify"
	"github.com/hipeday/azir-plugin-golang/pkg/properties"
	"go.uber.org/zap"
	"log"
	"net"
	"os"
	"path/filepath"
)

type Listen interface {
	Plugin
}

type ListenPlugin struct {
	Listen
	ConfigPlugin
}

func (l *ListenPlugin) Run(args []string) (interface{}, error) {
	var (
		err     error
		runConf string
	)

	var commands = []CommandParam{
		{
			Name:    "c",
			Default: "",
			Usage:   "Start plug in configuration",
			Value:   runConf,
		},
	}

	err = l.ParseConfig("run", commands, args)

	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	l.listenSocket(runConf, args)

	return "", nil
}

func (l *ListenPlugin) listenSocket(runConf string, args []string) {
	var (
		property   properties.DefaultProperty
		socketHome string
		socketPath string
		err        error
		logger     *zap.SugaredLogger
	)

	err = json.Unmarshal([]byte(runConf), &property)
	if err != nil {
		log.Fatalf("Error parsing run config: %v", err)
	} else {
		logger = l.GetLogger(args)
		socketHome = filepath.Join(property.Home, property.Name, "socks")
		socketPath = filepath.Join(socketHome, "plugin.sock")
	}

	err = os.MkdirAll(socketHome, os.ModePerm)
	if err != nil {
		logger.Fatalf("Error creating socket home: %v", err)
	}

	_, err = os.Stat(socketPath)
	if err == nil {
		err = os.Remove(socketPath)
		if err != nil {
			logger.Fatalf("Error removing socket: %v", err)
		}
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		logger.Fatalf("Server listen error: %v", err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			logger.Fatalf("Error closing listener: %v", err)
		}
	}(listener)

	logger.Infof("Plugin listening on %s", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Fatalf("Accept error: %v", err)
		}
		l.handleConnection(conn, logger, property.Notification)
	}
}

func (l *ListenPlugin) handleConnection(conn net.Conn, logger *zap.SugaredLogger, notification *properties.Notification) {

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Fatalf("Error closing connection: %v", err)
		}
	}(conn)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf[:])
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	logger.Infof("Received: %s", string(buf[:n]))

	if notification != nil && notification.Enabled {
		notifier := notify.Registry.GetNotification(notification.Type)
		notifier.SetLogger(logger)
		err = notifier.Push(string(buf[:n]), notification)
		if err != nil {
			logger.Fatalf("Error sending notification: %v", err)
		} else {
			logger.Infof("Notification sent")
		}
	} else {
		logger.Infof("Notification not enabled")
	}
}
