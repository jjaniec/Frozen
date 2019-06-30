package main

import (
	"fmt"
	"net"
)

type server struct {
	prefix   string
	listener net.Listener
	port     string
}

func (s *server) start() {
	srv, err := net.Listen("tcp4", s.port)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.listener = srv
}
