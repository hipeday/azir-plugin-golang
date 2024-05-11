package plugin

import (
	"encoding/json"
	"github.com/ideal-rucksack/workflow-glolang-plugin/cmd/command"
	"github.com/ideal-rucksack/workflow-glolang-plugin/pkg/plugin"
	"github.com/ideal-rucksack/workflow-glolang-plugin/pkg/properties"
	"net"
	"path/filepath"
)

var (
	pluginIns MySQLPlugin
)

func init() {
	pluginIns = MySQLPlugin{}
	command.Registry.RegisterCommand("run", plugin.CommandFunctions{Command: pluginIns.Run, Callback: nil})
	command.Registry.RegisterCommand("databases", plugin.CommandFunctions{Command: pluginIns.Databases, Callback: pluginIns.CallbackRender})
}

type MySQL interface {
	plugin.Callback
	Databases(args []string) (interface{}, error)
}

type MySQLPlugin struct {
	plugin.ListenPlugin
	MySQL
}

func (m *MySQLPlugin) Run(args []string) (interface{}, error) {
	return m.ListenPlugin.Run(args)
}

func (m *MySQLPlugin) Databases(args []string) (interface{}, error) {
	return []string{"db1", "db2"}, nil
}

func (m *MySQLPlugin) CallbackRender(result interface{}, args []string) error {
	var err error
	_, err = m.ParseConfig(args)
	if err != nil {
		return err
	}
	var (
		property   = m.GetConfig().(properties.DefaultProperty)
		logger     = m.GetLogger()
		socketHome = filepath.Join(property.Home, property.Name, "socks")
		socketPath = filepath.Join(socketHome, "plugin.sock")
	)

	conn, err := net.Dial("unix", socketPath)

	if err != nil {
		logger.Fatalf("Error dialing socket: %v", err)
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Fatalf("Error closing connection: %v", err)
		}
	}(conn)

	body, err := json.Marshal(result)

	if err != nil {
		return err
	}

	_, err = conn.Write(body)

	return err
}
