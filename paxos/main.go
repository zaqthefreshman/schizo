package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	protocal = "tcp"
)

var (
	host  string
	hosts []string
	conns []*net.TCPAddr
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing argument PORT")
		os.Exit(1)
	}
	port := os.Args[1]
	host := "localhost:" + port

	fmt.Println("Started")

	file, err := os.Open("../hosts")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if h := scanner.Text(); h != host {
			hosts = append(hosts, h)
			fmt.Println("Using host:", h)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	lostConnections := make(chan string, 0)
	newConnections := make(chan Server, 0)

	go findServers(lostConnections, newConnections)
	go healthcheck(newConnections, lostConnections)

	for _, v := range hosts {
		lostConnections <- v
	}

	l, err := net.Listen(protocal, host)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		l.Accept()
	}
}

func healthcheck(source chan Server, sink chan string) {
	//read conns from sink, store in array
	//loop and healthcheck
	//output failed conns connection info
	for {
		server := <-source
		log.Println("Will healthcheck:", server.hostname)
		//TODO: healthcheck :)
	}
}

func findServers(source chan string, sink chan Server) {
	//read strings output conns
	log.Println("Finding servers")
	for {
		time.Sleep(2 * time.Second)
		host := <-source
		log.Println("Recieved host:", host)
		conn, err := net.Dial(protocal, host)
		if err != nil {
			log.Println(err)
			go func() {
				source <- host
			}()
			continue
		}
		sink <- Server{conn, host}
	}
}

type Server struct {
	conn     net.Conn
	hostname string
}
