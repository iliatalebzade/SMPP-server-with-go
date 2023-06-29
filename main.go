package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func main() {
	// Start the TCP server
	listener, err := net.Listen("tcp", "0.0.0.0:2775")
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started. Listening on localhost:9000")

	for {
		// Accept incoming client connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle client request in a goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	fmt.Println("New client connected:", conn.RemoteAddr())

	// Read client request from the connection
	request := make([]byte, 1024)
	n, err := conn.Read(request)
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	fmt.Printf("Received request from client %s: %s\n", conn.RemoteAddr(), string(request[:n]))

	// Validate the request
	if isValidBindTransmitterRequest(request[:n]) {
		// Send a successful response
		response := []byte{0x00, 0x00, 0x00, 0x10, 0x80, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
		_, err = conn.Write(response)
		if err != nil {
			fmt.Println("Error sending response:", err)
		} else {
			fmt.Printf("Sent response to client %s\n", conn.RemoteAddr())
		}
	} else {
		// Send an error response
		response := []byte{0x00, 0x00, 0x00, 0x10, 0x80, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01}
		_, err = conn.Write(response)
		if err != nil {
			fmt.Println("Error sending response:", err)
		} else {
			fmt.Printf("Sent response to client %s\n", conn.RemoteAddr())
		}
	}
}

func isValidBindTransmitterRequest(request []byte) bool {
	// Validate the request based on the SMPP protocol specifications
	// Check the command length
	if len(request) < 16 {
		log.Println("Invalid PDU: Length was too short")
		return false
	}

	// Check the command ID
	commandID := binary.BigEndian.Uint32(request[4:8])
	if commandID != 0x00000002 {
		log.Println("Invalid PDU: the command id was not valid for a bind_transmitter request")
		return false
	}

	// Additional validation logic can be added here if required

	log.Println("Valid PDU!")
	return true
}
