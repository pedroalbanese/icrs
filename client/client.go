package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	var serverAddr, clientCertFile, clientKeyFile string
	flag.StringVar(&serverAddr, "server", "localhost:8000", "Server address")
	flag.StringVar(&clientCertFile, "cert", "", "Path to client certificate file")
	flag.StringVar(&clientKeyFile, "key", "", "Path to client private key file")
	flag.Parse()

	if clientCertFile == "" || clientKeyFile == "" {
		log.Fatal("Both -cert and -key flags must be provided")
	}

	// Load client certificate and key
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	tls.GOSTInstall()

	// Configure TLS connection
	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	// Connect to the server
	conn, err := tls.Dial("tcp", serverAddr, config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("Connected to server")

	// Read user input from stdin
	reader := bufio.NewReader(os.Stdin)

	// Join the "Home" room
	joinMessage := fmt.Sprintf("JOIN Home")
	_, err = conn.Write([]byte(joinMessage + "\n"))
	if err != nil {
		log.Println("Error sending join message:", err)
		return
	}

	client := &Client{
		conn: conn,
	}

	// Start reading messages from the server in a separate goroutine
	go readMessages(client)

	// Read user input and send messages to the server
	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "quit" {
			break
		}

		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			log.Println("Error sending message:", err)
			break
		}
	}

	log.Println("Disconnected from server")
}

// readMessages reads messages from the server and prints them to the console
func readMessages(client *Client) {
	reader := bufio.NewReader(client.conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading message from server:", err)
			break
		}

		fmt.Print(message)
	}

	log.Println("Disconnected from server")
}

// Client represents a connected client
type Client struct {
	conn net.Conn
}
