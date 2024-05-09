package plugin

import (
	"errors"
	"fmt"
	"github.com/ideal-rucksack/workflow-glolang-plugin/cmd/command"
	"github.com/ideal-rucksack/workflow-glolang-plugin/pkg/plugin"
	"github.com/ideal-rucksack/workflow-glolang-plugin/pkg/properties"
	"net"
	"os"
	"path/filepath"
	"time"
)

var (
	pluginIns MySQLPlugin
)

func init() {
	pluginIns = MySQLPlugin{}
	command.Registry.RegisterCommand("run", pluginIns.Run)
	command.Registry.RegisterCommand("databases", pluginIns.Databases)
}

type MySQL interface {
	Databases(args []string) (interface{}, error)
}

type MySQLPlugin struct {
	plugin.ListenPlugin
	MySQL
}

func (m *MySQLPlugin) Databases(args []string) (interface{}, error) {
	var err error
	_, err = m.ParseConfig(args)
	if err != nil {
		return nil, err
	}

	var (
		result = make(chan string)
	)

	go m.getData(result)

	select {
	case res := <-result:
		return res, err
	case <-time.After(6 * time.Second):
		return nil, errors.New("timeout")
	}
}

func (m *MySQLPlugin) getData(result chan<- string) {

	var (
		property   = m.GetConfig().(properties.DefaultProperty)
		logger     = property.Logger.CreateLogger(property.Name, property.InvokeId)
		socketHome = filepath.Join(property.Home, property.Name, "data")
		socketPath = filepath.Join(socketHome, property.InvokeId+".sock")
		err        error
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

	conn, err := listener.Accept()
	if err != nil {
		logger.Fatalf("Accept error: %v", err)
	}

	res, err := m.handleConnection(conn)
	if err != nil {
		fmt.Println("Handle connection error:", err)
	}
	result <- res
}

func (m *MySQLPlugin) handleConnection(conn net.Conn) (string, error) {
	var (
		property = m.GetConfig().(properties.DefaultProperty)
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
		return "", err
	}

	logger.Infof("[收到消息] - [%s] >> %s", time.Now().Format("2006-01-02 15:04:05.000000"), string(buf[:n]))

	return string(buf[:n]), nil
}
