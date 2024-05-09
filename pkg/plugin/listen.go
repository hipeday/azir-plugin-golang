package plugin

import (
	"fmt"
	"github.com/ideal-rucksack/workflow-glolang-plugin/pkg/properties"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
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
		err error
	)

	_, err = l.ParseConfig(args)

	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	l.listenSocket()

	return "", nil
}

func (l *ListenPlugin) listenSocket() {
	var (
		property   = l.GetConfig().(properties.DefaultProperty)
		socketHome = filepath.Join(property.Home, property.Name, "socks")
		socketPath = filepath.Join(socketHome, "plugin.sock")
		err        error
		logger     = property.Logger.CreateLogger(property.Name, property.InvokeId)
	)

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
		l.handleConnection(conn)
	}
}

func (l *ListenPlugin) handleConnection(conn net.Conn) {
	var (
		property = l.GetConfig().(properties.DefaultProperty)
		logger   = property.Logger.CreateLogger(property.Name, property.InvokeId)
	)

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

	logger.Infof("[收到消息] - [%s] >> %s", time.Now().Format("2006-01-02 15:04:05.000000"), string(buf[:n]))

	_, err = conn.Write([]byte("这是插件执行结果🌹"))
	if err != nil {
		logger.Fatalf("Failed to send message: %v", err)
	}
	logger.Infof("[回复消息] - [%s] >> %s", time.Now().Format("2006-01-02 15:04:05.000000"), "这是插件执行结果🌹")
}
