package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	connections sync.Map
	port        int
	listener    net.Listener
}

func NewServer(port int) *Server {
	return &Server{port: port}
}

func (s Server) Run() {
	portString := fmt.Sprintf(":%d", s.port)
	ln, _ := net.Listen("tcp", portString)
	s.listener = ln
	fmt.Println("Listening on port", s.port)
	s.loop()
}

func (s Server) loop() {
	for {
		connection, _ := s.listener.Accept()
		s.openConnection(connection)
		go s.handleUserConnetion(connection)
	}
}

func (s *Server) openConnection(connection net.Conn) {
	s.connections.Store(connection, true)
	fmt.Println("A new connection is now open")
}

func (s *Server) closeConnection(connection net.Conn) {
	s.connections.Swap(connection, false)
	connection.Close()
	fmt.Println("A connection is closed")
}

func (s *Server) handleUserConnetion(connection net.Conn) {
	defer s.closeConnection(connection)
	buffer := make([]byte, BufferSize)
	for {
		bytesRead, error := connection.Read(buffer)
		if error != nil {
			break
		}
		fmt.Print("Message Received: ", string(buffer[:bytesRead]))
	}
}
