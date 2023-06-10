package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pedroalbanese/readline"
	"github.com/pedroalbanese/color"
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

var (
	clientSKIDs   map[string]bool
	clientSKIDsMu sync.Mutex
)

var clients []Client
var rooms []*Room

// OID for Subject Key Identifier extension
var subjectKeyIdentifierOID = asn1.ObjectIdentifier{2, 5, 29, 14}
var authorityKeyIdentifierOID = []int{2, 5, 29, 35}

var (
	certFile   = flag.String("cert", "", "Certificate file path.")
	crlFile    = flag.String("crl", "", "Certificate revcation list.")
	keyFile    = flag.String("key", "", "Private key file path.")
	mode       = flag.String("mode", "client", "Mode: <server|client>")
	serverAddr = flag.String("ipport", "localhost:8000", "Server address.")
	strict     = flag.Bool("strict", false, "Restrict users.")
)

func init() {
    tls.GOSTInstall()
}

func main() {
	flag.Parse()

	clientSKIDs = make(map[string]bool)

	if *mode == "server" {
		// Load the server certificate and private key
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			log.Fatal(err)
		}

		var certPEM []byte
		file, err := os.Open(*certFile)
		if err != nil {
			log.Fatal(err)
		}
		info, err := file.Stat()
		if err != nil {
			log.Fatal(err)
		}
		buf := make([]byte, info.Size())
		file.Read(buf)
		certPEM = buf
		var certPemBlock, _ = pem.Decode([]byte(certPEM))
		var serverCert, _ = x509.ParseCertificate(certPemBlock.Bytes)

		var crl *pkix.CertificateList
		if *crlFile != "" {
		// Load the CRL from a file
			crlFile, err := os.Open(*crlFile)
			if err != nil {
			    log.Fatal(err)
			}
			defer crlFile.Close()

			// Decode the PEM block of the CRL
			crlBytes, err := ioutil.ReadAll(crlFile)
			if err != nil {
			    log.Fatal(err)
			}
			crl, err = x509.ParseCRL(crlBytes)
			if err != nil {
			    log.Fatal(err)
			}
		}

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

//			go handleClient(conn)
			go handleClient(conn, serverCert, crl)
		}
	} else {
		if *certFile == "" || *keyFile == "" {
			log.Fatal("Both -cert and -key flags must be provided")
		}

		// Load client certificate and key
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			log.Fatal(err)
		}

		// Configure TLS connection
		config := &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		}

		// Connect to the server
		conn, err := tls.Dial("tcp", *serverAddr, config)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		// Get the server's public key
		state := conn.ConnectionState()
		for _, v := range state.PeerCertificates {
			derBytes, err := x509.MarshalPKIXPublicKey(v.PublicKey)
			if err != nil {
				log.Fatal(err)
			}
			pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: derBytes})
			fmt.Printf("%s\n", pubPEM)
		}

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

			if message == "QUIT" {
				break
			}

			rl, err := readline.New("")
			if err != nil {
				fmt.Println("Erro ao criar o leitor de linha:", err)
				os.Exit(1)
			}
			defer rl.Close()
			
			rl.Stdout().Write([]byte("\033[1A\033[K"))
			printMessageln(message)
			_, err = conn.Write([]byte(message + "\n"))
			if err != nil {
				log.Println("Error sending message:", err)
				break
			}
		}

		log.Println("Disconnected from server")
	}
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

//		fmt.Print(message)
		printMessage(message)
	}

	log.Println("Disconnected from server")
}

//func handleClient(conn net.Conn) {
func handleClient(conn net.Conn, serverCert *x509.Certificate, CRLFile *pkix.CertificateList) {
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

	skid := getClientSKID(clientCert)

	// Check if the SKID is already registered
	clientSKIDsMu.Lock()
	if _, exists := clientSKIDs[skid]; exists {
		log.Println("Client already logged in.")
		message := "You are already logged in from another session."
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			log.Println("Error sending message to client:", err)
		}
		conn.Close()
		clientSKIDsMu.Unlock()
		return
	}
	clientSKIDs[skid] = true
	clientSKIDsMu.Unlock()

	if *strict {
		clientSKIDsMu.Lock()
		if !bytes.Equal(clientCert.AuthorityKeyId, serverCert.AuthorityKeyId) {
			message := "Invalid client certificate."
			_, err := conn.Write([]byte(message + "\n"))
			if err != nil {
				log.Println("Error sending message to client:", err)
			}
			conn.Close()
			delete(clientSKIDs, skid)
			clientSKIDsMu.Unlock()
			return
		}
		clientSKIDs[skid] = true
		clientSKIDsMu.Unlock()
	}

	if *crlFile != "" {
		revoked, revocationTime := isCertificateRevoked(clientCert, CRLFile)
		if revoked {
			message := "Your certificate has been revoked. Please contact the certificate authority.\nRevocation Time: " + revocationTime.String()
			_, err := conn.Write([]byte(message + "\n"))
			if err != nil {
				log.Println("Error sending message to client:", err)
			}
			conn.Close()
			delete(clientSKIDs, skid)
			return
		}
	}

	if isCertificateValid(clientCert) == false {
		message := "Your certificate has been expired."
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			log.Println("Error sending message to client:", err)
		}
		conn.Close()
		delete(clientSKIDs, skid)
		return
	}

	// Extract the username from the client certificate
//	username := strings.TrimPrefix(clientCert.Subject.CommonName, "CN=")
	username := "@" + strings.TrimPrefix(clientCert.Subject.CommonName, "CN=")

	client := Client{
		conn:        conn,
		username:    username,
		clientCert:  clientCert,
	}

//	message := fmt.Sprintf("%s joined the chat", client.username)
	message := fmt.Sprintf("%s joined the chat at %s", client.username, time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(message)
	fmt.Println("SKID:", getClientSKID(client.clientCert))
	fmt.Println("AKID:", getClientAKID(client.clientCert))
	fmt.Println("IP Address:", conn.RemoteAddr())
	fmt.Println("Certificate:")
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
			} else if strings.HasPrefix(message, "QUIT") {
				conn.Close()
				removeClient(&client)
				break 
			} else if strings.TrimSpace(message) == "LIST" {
				response := listUsers(client.room)
				_, err := client.conn.Write([]byte(response))
				if err != nil {
					log.Println("Error sending user list:", err)
				}
				continue
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
	delete(clientSKIDs, skid)
	removeClient(&client)
}



func printMessage(message string) {
//	currentTime := time.Now().Format("15:04:05")
//	fmt.Printf("[%s] %s", currentTime, message)
//	fmt.Print(message)
	if strings.HasPrefix(message, "Users in the chat") || strings.HasPrefix(message, "-") {
		fmt.Print(message)
	} else {
		currentTime := time.Now().Format("15:04:05")
		gray := color.New(color.FgHiBlack)
		gray.Printf("[%s] ", currentTime)
		if strings.HasPrefix(message, "@") && strings.Contains(message, "#") {
			split := strings.SplitN(message, "#", 2)
			if split[0] != "Joined room" && split[0] != "Left room" {
				red := color.New(color.FgHiWhite)
				red.Print(split[0])
				fmt.Print(":")
				fmt.Print(split[1])
				return
			} else {
				fmt.Print(split[0])
				fmt.Print(":")
				fmt.Print(split[1])
				return
			}
		} else {
			fmt.Print(message)
		}
	}
}

func printMessageln(message string) {
//	currentTime := time.Now().Format("15:04:05")
//	fmt.Printf("[%s] %s\n", currentTime, message)
//	fmt.Println(message)
	if strings.HasPrefix(message, "Users in the chat") || strings.HasPrefix(message, "-") {
		fmt.Println(message)
	} else {
		currentTime := time.Now().Format("15:04:05")
		gray := color.New(color.FgHiBlack)
		gray.Printf("[%s] ", currentTime)
		if strings.HasPrefix(message, "@") && strings.Contains(message, "#") {
			split := strings.SplitN(message, "#", 2)
			if split[0] != "Joined room" && split[0] != "Left room" {
				red := color.New(color.FgHiWhite)
				red.Print(split[0])
				fmt.Print(":")
				fmt.Println(split[1])
				return
			} else {
				fmt.Print(split[0])
				fmt.Print(":")
				fmt.Println(split[1])
				return
			}
		} else {
			fmt.Println(message)
		}
	}
}

func isCertificateRevoked(cert *x509.Certificate, crl *pkix.CertificateList) (bool, time.Time) {
	for _, revokedCert := range crl.TBSCertList.RevokedCertificates {
		if revokedCert.SerialNumber.Cmp(cert.SerialNumber) == 0 {
			return true, revokedCert.RevocationTime
		}
	}
	return false, time.Time{}
}

func isCertificateValid(cert *x509.Certificate) bool {
	currentTime := time.Now()
	if currentTime.Before(cert.NotBefore) || currentTime.After(cert.NotAfter) {
		return false
	}
	return true
}

func listUsers(room *Room) string {
	userList := "Users in the chat:\n"
	for _, client := range room.clients {
		userList += "- " + client.username + "\n"
	}
	return userList
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
	// Notify all clients in the room that a client has left
	for _, c := range room.clients {
		c.conn.Write([]byte(fmt.Sprintf("%s left the room.\n", client.username)))
	}
}

func sendMessage(client *Client, message string) {
	room := client.room
	room.mu.Lock()
	defer room.mu.Unlock()

	// Send the message to all clients in the same room except the sender
	for _, c := range room.clients {
		if c != client {
//			c.conn.Write([]byte(fmt.Sprintf("%s: %s\n", client.username, message)))
			c.conn.Write([]byte(fmt.Sprintf("%s# %s\n", client.username, message)))
		}
	}
}

func broadcastMessage(message string) {
	// Send the message to all connected clients
	for _, client := range clients {
		fmt.Fprintln(client.conn, message)
	}
}

func removeClient(client *Client) {
	// Check if the client is associated with a room
	if client.room == nil {
		return
	}

	// Lock the room's mutex to ensure exclusive access to the room's data
	client.room.mu.Lock()
	defer client.room.mu.Unlock()

	// Find the client in the room's client list and remove it
	for i, c := range client.room.clients {
		if c == client {
			// Create a new slice that excludes the client to be removed
			client.room.clients = append(client.room.clients[:i], client.room.clients[i+1:]...)
			break
		}
	}

	// Set the client's room reference to nil
	client.room = nil
}

func getClientSKID(cert *x509.Certificate) string {
	// Get the Subject Key Identifier (SKID) from the client certificate
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
type authorityKeyIdentifier struct {
	Raw       asn1.RawContent
	Authority []byte `asn1:"optional,tag:0"`
}

func getClientAKID(cert *x509.Certificate) string {
	// Get the Authority Key Identifier (AKID) from the client certificate
	for _, ext := range cert.Extensions {
		if ext.Id.Equal(authorityKeyIdentifierOID) {
			var akid authorityKeyIdentifier
			if _, err := asn1.Unmarshal(ext.Value, &akid); err == nil {
				if len(akid.Authority) > 0 {
					return fmt.Sprintf("%X", akid.Authority)
				}
			}
		}
	}
	return ""
}

func printClientCertPEM(cert *x509.Certificate) {
	// Print the client certificate in PEM format
	block := &pem.Block{
		Type:  "CERTIFICATE",
			Bytes: cert.Raw,
	}

	fmt.Println(string(pem.EncodeToMemory(block)))
}
