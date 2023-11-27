package main

import (
	"fmt"
	"net"
	"strconv"
)

const SERVER_HOST string = "localhost"
const SERVER_PORT uint16 = 9999

type ServerInstance struct {
	host     string
	port     uint16
	listener net.Listener
	live     bool
}

func beginServerInstance(host string, port uint16) ServerInstance {
	listenString := host + ":" + strconv.Itoa(int(port))
	fmt.Printf("Beginning master on %v.\n", listenString)
	listener, err := net.Listen("tcp", listenString)
	if err != nil {
		fmt.Printf("Master could not listen on %v. Panicking!\n", listenString)
	} else {
		fmt.Printf("Master successfully listening on %v.\n", listenString)
	}

	return ServerInstance{host, port, listener, true}
}

func (instance *ServerInstance) endServerInstance() {
	err := instance.listener.Close()
	if err != nil {
		fmt.Printf("Failed to close listener: \n%v\n", err)
	} else {
		fmt.Printf("Successfuly closed listener.\n")
	}
}

func (instance *ServerInstance) beginAcceptLoop() {
	fmt.Println("Beginning accept loop.")
	for {
		connection, err := instance.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting. Panicking!")
			panic(err)
		}
		fmt.Println("Accepted connection")
		go processClient(&connection)
	}
}

func processClient(connection *net.Conn) {
	defer func(connection *net.Conn) {
		_ = (*connection).Close()
	}(connection)

	buffer := make([]byte, 128)
	bLen, err := (*connection).Read(buffer)
	if err != nil {
		fmt.Println("Error reading. Panicking!")
		panic(err)
	}

	fmt.Printf("Received: %v\n\n", string(buffer[:bLen]))
	_, _ = (*connection).Write([]byte("Thanks!"))
}

func main() {
	instance := beginServerInstance(SERVER_HOST, SERVER_PORT)
	defer instance.endServerInstance()
	instance.beginAcceptLoop()
}
