package main

import (
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"net"
	"os"
	"strconv"
)

const MASTER_HOST string = "localhost"
const MASTER_PORT uint16 = 9999

type MasterLink struct {
	host       string
	port       uint16
	name       string
	connection net.Conn
}

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

func beginConnection(host string, port uint16, name string) MasterLink {
	connString := host + ":" + strconv.Itoa(int(port))
	connection, err := net.Dial("tcp", connString)

	if err != nil {
		fmt.Printf("Connection to master could not be established at %v. Panicking!\n", connString)
		panic(err)
	}

	link := MasterLink{host, port, name, connection}

	_, err = link.connection.Write([]byte(fmt.Sprintf("Hello from bot %v!", name)))
	if err != nil {
		fmt.Printf("Error writing: %v\n", err)
	}
	buffer := make([]byte, 128)
	bLen, err := link.connection.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading: %v\n", err)
	}
	fmt.Println("Received: ", string(buffer[:bLen]))

	return link
}

func (link *MasterLink) endConnection() {
	fmt.Println(link.connection)
	_ = link.connection.Close()
}

func main() {
	botName, err := generateBotName()
	if err != nil {
		fmt.Printf("Failed to generate unique bot name, %v. Using \"%v\"\n", err, botName)
	} else {
		fmt.Printf("Generated unique bot name: \"%v\"\n", botName)
	}

	link := beginConnection(MASTER_HOST, MASTER_PORT, botName)
	defer link.endConnection()
}
