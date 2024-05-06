package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

func main() {

	result := make(chan string)
	go func() {
		cmd := exec.Command("/Users/coloxan/GolandProjects/webhook-golang/plugin/foo", "exec", "-payload", `{"nickname": "ggboy"}`)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()
	go getData(result)

	select {
	case res := <-result:
		fmt.Printf("获得结果: %s", res)
	case <-time.After(6 * time.Second):
		fmt.Println("timeout")
	}
}

func getData(result chan<- string) {
	socketPath := "/Users/coloxan/GolandProjects/webhook-golang/sock/foo/data.sock"
	os.Remove(socketPath) // 确保Socket文件不存在
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
