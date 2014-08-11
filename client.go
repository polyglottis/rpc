package rpc

import (
	"fmt"
	"log"
	"net/rpc"
	"strings"
	"time"
)

type Client struct {
	Name    string
	Network string
	Address string
	c       *rpc.Client
}

func NewClient(name, network, address string) (*Client, error) {
	c := &Client{
		Name:    name,
		Network: network,
		Address: address,
	}
	err := c.redial()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) redial() error {
	if c.c != nil {
		err := c.c.Close()
		c.c = nil
		if err != nil && err != rpc.ErrShutdown {
			return err
		}
	}
	return c.tryDialing()
}

func (c *Client) tryDialing() error {
	var err error
	c.c, err = rpc.Dial(c.Network, c.Address)
	if err != nil {
		if strings.HasSuffix(err.Error(), "connection refused") {
			log.Printf("Unable to reach %s rpc server. Trying again later.", c.Name)
			return nil
		}
		return err
	}
	return nil
}

var waitIntervals = []time.Duration{100 * time.Millisecond, 1 * time.Second, 3 * time.Second, 5 * time.Second}

// Call performs the rpc call, trying to redial if the connection was shut down.
func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	if c.c != nil {
		err := c.c.Call(serviceMethod, args, reply)
		if err == nil || err != rpc.ErrShutdown {
			return err
		}
	}
	var err error
	for _, duration := range waitIntervals {
		log.Printf("%s rpc client recovering...", c.Name)
		time.Sleep(duration)
		err = c.redial()
		if err != nil {
			return err
		}
		if c.c != nil {
			err = c.c.Call(serviceMethod, args, reply)
			if err == nil || err != rpc.ErrShutdown {
				log.Printf("%s rpc client recovered!", c.Name)
				return err
			}
		}

	}
	return fmt.Errorf("Unable to execute %s rpc call", c.Name)
}
