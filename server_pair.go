package rpc

import (
	"log"
)

type ServerPair struct {
	name string
	main *Server
	op   *Server
}

// NewServerPair creates a new rpc server pair with the given name.
func NewServerPair(name string, main, operations *Server) *ServerPair {
	return &ServerPair{
		name: name,
		main: main,
		op:   operations,
	}
}

func (p *ServerPair) RegisterAndListen() error {
	err := p.main.RegisterAndListen()
	if err != nil {
		return err
	}
	log.Printf("Main %s listening on %v", p.name, p.main.Addr)

	err = p.op.RegisterAndListen()
	if err != nil {
		return err
	}
	log.Printf("Operations %s listening on %v", p.name, p.op.Addr)
	return nil
}

func (p *ServerPair) Close() error {
	err := p.main.Close()
	if err != nil {
		return err
	}
	return p.op.Close()
}

// Accept starts accepting calls. This method blocks.
func (p *ServerPair) Accept() {
	log.Printf("Launching %s...", p.name)
	go p.op.Accept()
	p.main.Accept()
}
