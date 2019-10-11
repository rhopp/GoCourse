package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConnection(c net.Conn) {
	defer c.Close()
	netData, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("Received: " + netData)
	c.Write([]byte("You said: " + netData))
}

func handleConnectionRoutine(c chan net.Conn) {
	for connection := range c {
		handleConnection(connection)
	}
}

func main() {

	//create channel
	channel := make(chan net.Conn)

	// close channel on exit
	defer func(c chan net.Conn) {
		close(c)
	}(channel)

	// create 4 goroutines
	for i := 1; i < 5; i++ {
		go handleConnectionRoutine(channel)
	}

	l, err := net.Listen("tcp4", ":5002")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		// send accepted connection to channel
		channel <- c
	}
}
