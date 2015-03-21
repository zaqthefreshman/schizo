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
	conn_host = "localhost"
	conn_type = "tcp"
)

var (
	hosts       = []string{"8086", "8087", "8088"}
	connections = make([]net.Conn, 0, len(hosts))
	count       = 0
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing argument PORT")
		os.Exit(1)
	}
	conn_port := os.Args[1]

	processPort := findStringInSlice(conn_port, hosts)
	if processPort == -1 {
		fmt.Println("Host not in pool, please use a different port")
		os.Exit(1)
	}

	hosts = append(hosts[:processPort], hosts[processPort+1:]...)

	go findFriends()

	// Listen for incoming connections
	l, err := net.Listen(conn_type, conn_host+":"+conn_port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + conn_host + ":" + conn_port)

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
		fmt.Println(fmt.Sprintf("LocalCount:    %d", count))
		fmt.Println(fmt.Sprintf("ReceivedCount: %s", buf))
	}
}

func countStuff() {
	for {
		time.Sleep(1000 * time.Millisecond)
		count++
		for _, v := range connections {
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
	friendsFound := []String{}
	for {
		if len(friendsFound) >= len(hosts) {
			fmt.Println("Looking for friends")
			for _, v := range hosts {
				tcpAddr, _ := net.ResolveTCPAddr(conn_type, conn_host+":"+v)
				conn, err := net.DialTCP(conn_type, nil, tcpAddr)
				if err != nil {
					fmt.Println("No available host at: " + v)
				} else {
					conn.SetKeepAlive(true)
					fmt.Println("Added Host :" + v)
					connections = append(connections, conn)
				}
			}
			//FIX logic is satified with finding only one friend needs to find all of them
			for _, v := range connections {
				friend := strings.Split(v.RemoteAddr().String(), ":")[1]
				if findStringInSlice(friend, hosts) == -1 {
					continue
				}
				friendsFound = append(friendsFound, friend)
			}
		}
		time.Sleep(10000 * time.Millisecond)
	}
}

func ArrayCompare(hosts, friendsFound) bool {
	for _, v := range hosts {
		for _, v2 := range friendsFound {
			return true
		}
	}
	return false
}
