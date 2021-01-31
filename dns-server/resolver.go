package main

import (
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

type Resolver interface {
	ResolveDNS(message dnsmessage.Message) (dnsmessage.Message, error)
}

type UDPResolver struct {
	conn           net.Conn
	defaultBufSize int
}

func NewUDPResolver(conn net.Conn) *UDPResolver {
	return &UDPResolver{
		conn:           conn,
		defaultBufSize: 512,
	}
}

func (r *UDPResolver) ResolveDNS(msg dnsmessage.Message) (dnsmessage.Message, error) {
	packedMsg, err := msg.Pack()
	if err != nil {
		return dnsmessage.Message{}, err
	}
	_, err = r.conn.Write(packedMsg)
	if err != nil {
		return dnsmessage.Message{}, err
	}
	resBuf := make([]byte, r.defaultBufSize)
	_, err = r.conn.Read(resBuf)
	if err != nil {
		return dnsmessage.Message{}, err
	}

	var resMsg dnsmessage.Message
	err = resMsg.Unpack(resBuf)
	if err != nil {
		return dnsmessage.Message{}, err
	}
	return resMsg, nil
}
