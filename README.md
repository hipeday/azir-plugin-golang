# Azir Golang Plugin

Azir Workflow scheduler plugin api for golang

## Usage

### Create a new plugin

> Create a new plugin using the following command
> 
> Suppose we now want to create a plugin named `foo`, we can create a new plugin using the following command

- Create a new plugin directory

```shell
mkdir foo-plugin
cd foo-plugin
```

- Create a new plugin(initialize go module replace `${username}` with your GitHub username)

```shell
go mod init github.com/${username}/foo-plugin
```

- Install the plugin dependency

```shell
go get github.com/ideal-rucksack/workflow-golang-plugin
```

- Create a new plugin implementation file

```shell
mkdir plugin
touch plugin/foo.go
vi plugin/foo.go
```

- Add the following code to the `plugin/foo.go` file

```go
package plugin

import (
	"encoding/json"
	"github.com/ideal-rucksack/workflow-golang-plugin/cmd/command"
	"github.com/ideal-rucksack/workflow-golang-plugin/pkg/plugin"
	"github.com/ideal-rucksack/workflow-golang-plugin/pkg/properties"
	"net"
	"path/filepath"
)

var (
	pluginIns FooPlugin
)

// Register the plugin commands
func init() {
	pluginIns = FooPlugin{}
	command.Registry.RegisterCommand("run", plugin.CommandFunctions{Command: pluginIns.Run, Callback: nil})
	command.Registry.RegisterCommand("databases", plugin.CommandFunctions{Command: pluginIns.Databases, Callback: pluginIns.CallbackRender})
}

// Foo this interface is Foo plugin api
type Foo interface {
	plugin.Callback
	// Databases returns the list of databases
	Databases(args []string) (interface{}, error)
}

// FooPlugin this struct is Foo plugin api implementation
type FooPlugin struct {
	plugin.ListenPlugin
	Foo
}

func (m *FooPlugin) Databases(args []string) (interface{}, error) {
	return []string{"db1", "db2"}, nil
}

// CallbackRender this function is Foo plugin api implementation for CallbackRender
// this function invokes the plugin socket server to render the result
// You can render the results where you want. For example, by default, the Run function will listen to a socket service when it is started.
// At this point your data should be rendered to this socket service
func (m *FooPlugin) CallbackRender(result interface{}, args []string) error {
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

	// create a new socket connection
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

	// sent the result to the socket
	_, err = conn.Write(body)

	return err
}
```

- Create a new plugin main file

```shell
touch main.go
vi main.go
```

- Add the following code to the `main.go` file

```go
package main

import (
	"github.com/ideal-rucksack/workflow-golang-plugin/cmd/runner"
	_ "github.com/${username}/foo-plugin/plugin"
)

func main() {
	runner.Run()
}
```

- Build the plugin

```shell
go build -o foo
```

- Run the plugin

> 启动插件的前提我们需要一些配置才能启动插件

```json
{
  "name": "foo-go-100",
  "language": "golang",
  "invoke_id": "1", 
  "home": "/app/workflow/plugins",
  "notification": {
    "type": "UNIX",
    "address": "/app/workflow/plugins/foo-go-100/data/1.sock",
    "enabled": true
  },
  "suffix": "",
  "version": "1.0.0",
  "description": "Foo plugin for golang",
  "logPath": "${name}_${language}_${version}.log",
  "cmd": {
    "linux": ["foo"],
    "darwin": ["foo"],
    "windows": ["foo"]
  },
  "parameter": {
    "webhook": {
      "type": "string",
      "description": "Webhook URL",
      "required": false,
      "default": null
    },
    "action": {
      "type": "string",
      "description": "Action: databases, tables, columns, rows",
      "required": true,
      "default": null
    }
  }
}
```

```shell
./foo run -c '{"name":"foo-go-100","language":"golang","invoke_id":"1","home":"/app/workflow/plugins","notification":{"type":"UNIX","address":"/app/workflow/plugins/foo-go-100/data/1.sock","enabled":true},"suffix":"","version":"1.0.0","description":"Foopluginforgolang","logPath":"${name}_${language}_${version}.log","cmd":{"linux":["foo"],"darwin":["foo"],"windows":["foo"]},"parameter":{"webhook":{"type":"string","description":"WebhookURL","required":false,"default":null},"action":{"type":"string","description":"Action:databases,tables,columns,rows","required":true,"default":null}}}'
```

### Test the plugin

- Get databases command invoke result

```go
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {

	result := make(chan string)
	go getData(result)

	select {
	case res := <-result:
		fmt.Printf("获得结果: %s", res)
	}
}

func getData(result chan<- string) {
	socketPath := "/app/workflow/plugins/foo-go-100/data/1.sock"
	// 确保Socket文件不存在
	os.Remove(socketPath) 
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("Server listen error:", err)
	}

	defer listener.Close()

	fmt.Println("Plugin listening on", socketPath)

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Accept error:", err)
	}

	res, err := handleConnection(conn)
	if err != nil {
		fmt.Println("Handle connection error:", err)
	}
	result <- res
}

func handleConnection(conn net.Conn) (string, error) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf[:])
	if err != nil {
		return "", err
	}

	log.Printf("[收到消息] - [%s] >> %s", time.Now().Format("2006-01-02 15:04:05.000000"), string(buf[:n]))

	return string(buf[:n]), nil
}
```

Execute the above code to listen the result of the `databases` command

- Test the plugin

```shell
./foo databases -c '{"name":"foo-go-100","language":"golang","invoke_id":"1","home":"/app/workflow/plugins","notification":{"type":"UNIX","address":"/app/workflow/plugins/foo-go-100/data/1.sock","enabled":true},"suffix":"","version":"1.0.0","description":"Foopluginforgolang","logPath":"${name}_${language}_${version}.log","cmd":{"linux":["foo"],"darwin":["foo"],"windows":["foo"]},"parameter":{"webhook":{"type":"string","description":"WebhookURL","required":false,"default":null},"action":{"type":"string","description":"Action:databases,tables,columns,rows","required":true,"default":null}}}'
```
