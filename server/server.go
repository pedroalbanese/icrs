package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Client struct {
	conn        net.Conn
	username    string
	clientCert  *x509.Certificate
	room        *Room
}

type Room struct {
	name    string
	clients []*Client
	mu      sync.Mutex
}

var clients []Client
var rooms []*Room

// OID for Subject Key Identifier extension
var subjectKeyIdentifierOID = asn1.ObjectIdentifier{2, 5, 29, 14}

func main() {
	certFile := flag.String("cert", "server.crt", "Server certificate file path")
	keyFile := flag.String("key", "server.key", "Server private key file path")
	flag.Parse()

	// Load the server certificate and private key
	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}

	tls.GOSTInstall()

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAnyClientCert,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}

	listener, err := tls.Listen("tcp", "localhost:8000", config)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Chat server started. Waiting for TLS connections...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		log.Println("Connection is not a TLS connection.")
		conn.Close()
		return
	}

	// Verify the client certificate
	err := tlsConn.Handshake()
	if err != nil {
		log.Println("Failed to perform TLS handshake:", err)
		conn.Close()
		return
	}

	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		log.Println("Client certificate is missing.")
		conn.Close()
		return
	}

	clientCert := state.PeerCertificates[0]

	// Extract the username from the client certificate
	username := strings.TrimPrefix(clientCert.Subject.CommonName, "CN=")

	client := Client{
		conn:        conn,
		username:    username,
		clientCert:  clientCert,
	}

//	message := fmt.Sprintf("%s joined the chat", client.username)
	message := fmt.Sprintf("%s joined the chat at %s", client.username, time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(message)
	fmt.Println("Client SKID:", getClientSKID(client.clientCert))
	fmt.Println("Client Certificate:")
	printClientCertPEM(client.clientCert)
	broadcastMessage(message)

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		message = strings.TrimSpace(message)

		if client.room == nil {
			if strings.HasPrefix(message, "JOIN ") {
				roomName := strings.TrimPrefix(message, "JOIN ")
				room := findOrCreateRoom(roomName)
				joinRoom(&client, room)
			} else {
				conn.Write([]byte("You are not in a room. Use JOIN <room> command to join a room.\n"))
			}
		} else {
			if strings.HasPrefix(message, "JOIN ") {
				leaveRoom(&client)
				roomName := strings.TrimPrefix(message, "JOIN ")
				room := findOrCreateRoom(roomName)
				joinRoom(&client, room)
			} else if strings.HasPrefix(message, "LEAVE") {
				leaveRoom(&client)
			} else {
				sendMessage(&client, message)
			}
		}
	}

//	message = fmt.Sprintf("%s left the chat", client.username)
	message = fmt.Sprintf("%s left the chat at %s", client.username, time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(message)
	broadcastMessage(message)

	conn.Close()

	removeClient(client)
}

func findOrCreateRoom(roomName string) *Room {
	for _, room := range rooms {
		if room.name == roomName {
			return room
		}
	}

	room := &Room{
		name:    roomName,
		clients: make([]*Client, 0),
	}

	rooms = append(rooms, room)

	return room
}

func joinRoom(client *Client, room *Room) {
	room.mu.Lock()
	defer room.mu.Unlock()

	client.room = room
	client.room.clients = append(client.room.clients, client)

	client.conn.Write([]byte(fmt.Sprintf("Joined room: %s\n", room.name)))

	// Notificar os demais clientes da sala sobre o novo ingresso
	notifyClientJoined(room, client)
}

func notifyClientJoined(room *Room, newClient *Client) {
	for _, client := range room.clients {
		if client != newClient {
			client.conn.Write([]byte(fmt.Sprintf("%s joined the room.\n", newClient.username)))
		}
	}
}

func leaveRoom(client *Client) {
	if client.room == nil {
		return
	}

	room := client.room
	room.mu.Lock()
	defer room.mu.Unlock()

	for i, c := range room.clients {
		if c == client {
			room.clients = append(room.clients[:i], room.clients[i+1:]...)
			client.room = nil
			client.conn.Write([]byte(fmt.Sprintf("Left room: %s\n", room.name)))

			notifyClientLeft(room, client)
			return
		}
	}
}

func notifyClientLeft(room *Room, client *Client) {
	for _, c := range room.clients {
		c.conn.Write([]byte(fmt.Sprintf("%s left the room.\n", client.username)))
	}
}

func sendMessage(client *Client, message string) {
	room := client.room
	room.mu.Lock()
	defer room.mu.Unlock()

	for _, c := range room.clients {
		if c != client {
			c.conn.Write([]byte(fmt.Sprintf("%s: %s\n", client.username, message)))
		}
	}
}

func broadcastMessage(message string) {
	for _, client := range clients {
		fmt.Fprintln(client.conn, message)
	}
}

func removeClient(client Client) {
	for i, c := range clients {
		if c == client {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func getClientSKID(cert *x509.Certificate) string {
	for _, ext := range cert.Extensions {
		if ext.Id.Equal(subjectKeyIdentifierOID) {
			var skid []byte
			if _, err := asn1.Unmarshal(ext.Value, &skid); err == nil {
				return fmt.Sprintf("%X", skid)
			}
		}
	}
	return ""
}

func printClientCertPEM(cert *x509.Certificate) {
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}

	fmt.Println(string(pem.EncodeToMemory(block)))
}
