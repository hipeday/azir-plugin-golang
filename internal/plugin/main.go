package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Job struct {
	Nickname string `json:"nickname"`
}

type Plugin interface {
	Invoke(Job) string
}

type MyPlugin struct {
}

func (p *MyPlugin) Invoke(job Job) string {
	return "你好👋Unix Domain Socket！"
}

func main() {
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	executeCmd := flag.NewFlagSet("exec", flag.ExitOnError)
	payloadArg := executeCmd.String("payload", "", "JSON payload")

	if len(os.Args) < 2 {
		fmt.Println("expected 'run' or 'exec' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		runCmd.Parse(os.Args[2:])
		run()
	case "exec":
		executeCmd.Parse(os.Args[2:])
		execute(*payloadArg)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func run() {
	socketPath := "/Users/coloxan/GolandProjects/webhook-golang/sock/plugin.sock"
	os.Remove(socketPath) // 确保Socket文件不存在

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("Server listen error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Plugin listening on", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			return
		}
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf[:])
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	log.Printf("[收到消息] - [%s] >> %s", time.Now().Format("2006-01-02 15:04:05.000000"), string(buf[:n]))

	var job Job
	err = json.Unmarshal(buf[:n], &job)
	if err != nil {
		fmt.Println("Failed to parse job:", err)
		return
	}

	plugin := &MyPlugin{} // Create your plugin here
	result := plugin.Invoke(job)

	_, err = conn.Write([]byte(result))
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return
	}
	log.Printf("[回复消息] - [%s] >> %s", time.Now().Format("2006-01-02 15:04:05.000000"), result)
}

func execute(payload string) {
	conn, err := net.Dial("unix", "/Users/coloxan/GolandProjects/webhook-golang/sock/plugin.sock")
	if err != nil {
		fmt.Println("Failed to connect to plugin:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(payload))
	if err != nil {
		fmt.Println("Failed to send job:", err)
		return
	}
	log.Printf("[发送消息] - [%s] >> %s", time.Now().Format("2006-01-02 15:04:05.000000"), payload)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf[:])
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return
	}
	log.Printf("[收到回复] - [%s] >> %s", time.Now().Format("2006-01-02 15:04:05.000000"), string(buf[:n]))
	dial, err := net.Dial("unix", "/Users/coloxan/GolandProjects/webhook-golang/sock/foo/data.sock")
	if err != nil {
		fmt.Println("Failed to connect to plugin:", err)
		return
	}
	defer dial.Close()
	dial.Write(buf[:n])
}
