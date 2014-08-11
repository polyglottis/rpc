// Package rpc provides helper structs for rpc server and client implementations.
package rpc

import (
	"net"
	"net/rpc"
)

type Server struct {
	name      string
	rpcServer interface{}
	Addr      string
	listener  net.Listener
}

// NewServer creates a new server providing with the given rpc name, rpc object, and tcp address.
func NewServer(name string, rpcServer interface{}, addr string) *Server {
	return &Server{
		name:      name,
		rpcServer: rpcServer,
		Addr:      addr,
	}
}

// RegisterAndListen does what is says. The caller is responsible to Close() the server afterwards.
func (s *Server) RegisterAndListen() error {
	err := rpc.RegisterName(s.name, s.rpcServer)
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.listener = l
	return nil
}

func (s *Server) Close() error {
	if s == nil || s.listener == nil {
		return nil
	}
	return s.listener.Close()
}

// Accept starts accepting calls. This method blocks.
func (s *Server) Accept() {
	rpc.Accept(s.listener)
}
