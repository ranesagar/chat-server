package main

import (
	"net"
	"strings"
	"testing"
	"time"
)

func TestChatServer(t *testing.T) {
	// Start the server
	go main()

	// Wait for the server to start
	time.Sleep(time.Second)

	// Connect two clients
	client1, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Errorf("Error connecting client 1: %s", err)
	}
	defer client1.Close()

	client2, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Errorf("Error connecting client 2: %s", err)
	}
	defer client2.Close()

	// Send messages
	sendMessage(t, client1, "Hello from client 1")
	sendMessage(t, client2, "Hello from client 2")

	// Receive messages
	receivedMessage(t, client1, "< Hello from client 2")
	receivedMessage(t, client2, "< Hello from client 1")
}

func sendMessage(t *testing.T, conn net.Conn, msg string) {
	_, err := conn.Write([]byte(msg + "\n"))
	if err != nil {
		t.Errorf("Error sending message: %s", err)
	}
}

func receivedMessage(t *testing.T, conn net.Conn, expected string) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Errorf("Error receiving message: %s", err)
	}
	receivedMsg := string(buf[:n])
	receivedMsg = strings.TrimSpace(receivedMsg) // Trim any leading/trailing whitespace
	if receivedMsg != expected {
		t.Errorf("Expected message '%s', but received '%s'", expected, receivedMsg)
	}
}
