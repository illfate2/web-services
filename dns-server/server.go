package main

import (
	"context"
	"log"
	"net"
	"strconv"

	"golang.org/x/net/dns/dnsmessage"
)

type forwardServers struct {
	ip   net.IP
	port int
}

type DNSServer struct {
	forward        forwardServers
	conn           *net.UDPConn
	defaultBufSize int
}

func NewDNSServer(conn *net.UDPConn, forward forwardServers) *DNSServer {
	return &DNSServer{
		forward:        forward,
		conn:           conn,
		defaultBufSize: 512,
	}
}

func (s *DNSServer) resolveDNS(dnsMsg dnsmessage.Message) (dnsmessage.Message, error) {
	var m dnsmessage.Message
	conn, err := net.Dial("udp", s.forward.ip.String()+":"+strconv.Itoa(s.forward.port))
	if err != nil {
		return dnsmessage.Message{}, err
	}
	packedDnsMsg, err := dnsMsg.Pack()
	if err != nil {
		return dnsmessage.Message{}, err
	}
	_, err = conn.Write(packedDnsMsg)
	if err != nil {
		return dnsmessage.Message{}, err
	}
	resBuf := make([]byte, s.defaultBufSize)
	_, err = conn.Read(resBuf)
	if err != nil {
		return dnsmessage.Message{}, err
	}
	err = m.Unpack(resBuf)
	if err != nil {
		return dnsmessage.Message{}, err
	}
	return m, nil
}

func (s *DNSServer) readDNSMsg() (dnsmessage.Message, *net.UDPAddr, error) {
	buf := make([]byte, s.defaultBufSize)
	_, addr, err := s.conn.ReadFromUDP(buf)
	if err != nil {
		return dnsmessage.Message{}, nil, err
	}
	var msg dnsmessage.Message
	err = msg.Unpack(buf)
	if err != nil {
		return dnsmessage.Message{}, nil, err
	}
	return msg, addr, nil
}

func (s *DNSServer) handle(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		msg, clientAddr, err := s.readDNSMsg()
		if err != nil {
			log.Print(err)
			continue
		}
		resolvedMsg, err := s.resolveDNS(msg)
		if err != nil {
			log.Print(err)
			continue
		}
		err = s.sendDNSMsg(clientAddr, resolvedMsg)
		if err != nil {
			log.Print(err)
			continue
		}
	}
}

func (s *DNSServer) sendDNSMsg(addr *net.UDPAddr, message dnsmessage.Message) error {
	packed, err := message.Pack()
	if err != nil {
		return err
	}
	_, err = s.conn.WriteToUDP(packed, addr)
	return err
}
