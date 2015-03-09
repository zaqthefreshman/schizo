package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_TYPE = "tcp"
)

var (
	HOSTS = []string{"8086", "8087", "8088"}
	CONNS = make([]net.Conn, 0, len(HOSTS))
	count = 0
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing argument PORT")
		os.Exit(1)
	}
	CONN_PORT := os.Args[1]

	processPort := findStringInSlice(CONN_PORT, HOSTS)
	if processPort == -1 {
		fmt.Println("Host not in pool, please use a different port")
		os.Exit(1)
	} else {
		HOSTS = append(HOSTS[:processPort], HOSTS[processPort+1:]...)
	}

	go findFriends()

	// Listen for incoming connections
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	go countStuff()

	for {
		// Listen for an incoming connection.
		fmt.Println("waiting on " + l.Addr().String())
		conn, err := l.Accept()
		fmt.Println("recieved")
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	for {
		// Make a buffer to hold incoming data.
		buf := make([]byte, 1024)
		// Read the incoming connection into the buffer.
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		fmt.Println(fmt.Sprintf("LocalCount:    %d\nReceivedCount: %s", count, buf))
	}
}

func countStuff() {
	for {
		time.Sleep(1000 * time.Millisecond)
		count++
		for _, v := range CONNS {
			fmt.Println("Sending count to " + v.RemoteAddr().String())
			v.Write([]byte(strconv.Itoa(count)))
		}
	}
}

func findStringInSlice(s string, l []string) int {
	for i, v := range l {
		if s == v {
			return i
		}
	}
	return -1
}

func findFriends() {
	for {
		fmt.Println("Looking for friends")
		for _, v := range HOSTS {
			tcpAddr, _ := net.ResolveTCPAddr(CONN_TYPE, CONN_HOST+":"+v)
			conn, err := net.DialTCP(CONN_TYPE, nil, tcpAddr)
			if err != nil {
				fmt.Println("No available host at: " + v)
			} else {
				conn.SetKeepAlive(true)
				fmt.Println("Added Host :" + v)
				CONNS = append(CONNS, conn)
			}
		}
		//FIX logic is satified with finding only one friend needs to find all of them
		for _, v := range CONNS {
			if findStringInSlice(strings.Split(v.RemoteAddr().String(), ":")[1], HOSTS) == -1 {
				continue
			}
			fmt.Println("Done looking for friends")
			return
		}
		time.Sleep(10000 * time.Millisecond)
	}
}
