package main

import (
	"encoding/binary"
	"log"
	"net"
)

const (
	bindTransmitterCommandID = 0x00000002
)

func main() {
	// Start the TCP server
	listener, err := net.Listen("tcp", "0.0.0.0:2775")
	if err != nil {
		log.Println("Failed to start server:", err)
		return
	}
	defer listener.Close()

	log.Println("Server started. Listening on 0.0.0.0:2775")

	for {
		// Accept incoming client connections
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		// Handle client request in a goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	log.Println("New client connected:", conn.RemoteAddr())

	// Read client request from the connection
	request := make([]byte, 1024)
	n, err := conn.Read(request)
	if err != nil {
		log.Println("Error reading request:", err)
		return
	}

	log.Printf("Received request from client %s: %s\n", conn.RemoteAddr(), string(request[:n]))

	// Validate the request
	if isValidRequest(request[:n], bindTransmitterCommandID) {
		// Send a successful response
		response := []byte{0x00, 0x00, 0x00, 0x10, 0x80, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
		sendResponse(conn, response)
		log.Printf("Sent response to client %s\n", conn.RemoteAddr())
	} else {
		// Send an error response
		response := []byte{0x00, 0x00, 0x00, 0x10, 0x80, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01}
		sendResponse(conn, response)
		log.Printf("Sent response to client %s\n", conn.RemoteAddr())
	}
}

func isValidRequest(request []byte, expectedCommandID uint32) bool {
	// Validate the request based on the SMPP protocol specifications
	// Check the command length
	if len(request) < 16 {
		log.Println("Invalid PDU: Length was too short")
		return false
	}

	// Check the command ID
	commandID := binary.BigEndian.Uint32(request[4:8])
	if commandID != expectedCommandID {
		log.Println("Invalid PDU: the command ID was not valid for the expected command")
		return false
	}

	// Additional validation logic can be added here if required

	log.Println("Valid PDU!")
	return true
}

func sendResponse(conn net.Conn, response []byte) {
	_, err := conn.Write(response)
	if err != nil {
		log.Println("Error sending response:", err)
	}
}
