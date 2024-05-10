package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
)

type Client struct {
	conn net.Conn
	name string
}

var (
	clients       []*Client
	clientsMu     sync.Mutex
	clientCounter int32 // Using int32 for atomic operations
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
	defer ln.Close()

	fmt.Println("Chat server started on :8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Add client to the list
	client := &Client{
		conn: conn,
		name: "Client " + strconv.Itoa(int(atomic.AddInt32(&clientCounter, 1))),
	}
	addClient(client)
	defer removeClient(client)

	fmt.Fprintf(conn, "You are %s\n", client.name)

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Error reading from client: %s", err)
			break
		}
		msg := string(buf[:n])
		forwardMessage(client, msg)
	}
}

func addClient(client *Client) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	clients = append(clients, client)
}

func removeClient(client *Client) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for i, c := range clients {
		if c == client {
			clients = append(clients[:i], clients[i+1:]...)
			return
		}
	}
}

func forwardMessage(sender *Client, msg string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for _, client := range clients {
		if client != sender {
			fmt.Fprintf(client.conn, "< %s\n", msg)
		}
	}
}
