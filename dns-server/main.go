package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"golang.org/x/net/dns/dnsmessage"
)

type forwardServers struct {
	ip   net.IP
	port int
}

type DNSServer struct {
	forwards []forwardServers
	conn     *net.UDPConn
}

func NewDNSServer(conn *net.UDPConn, forwards []forwardServers) *DNSServer {
	return &DNSServer{
		forwards: forwards,
		conn:     conn,
	}
}

func (s *DNSServer) handle(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		buf := make([]byte, 512)
		_, addr, _ := s.conn.ReadFromUDP(buf)
		var m dnsmessage.Message
		err := m.Unpack(buf)
		if err != nil {
			log.Print(err)
		}
		packed, err := m.Pack()
		if err != nil {
			log.Print(err)
		}
		resolver := net.UDPAddr{IP: s.forwards[0].ip, Port: s.forwards[0].port}
		_, err = s.conn.WriteToUDP(packed, &resolver)
		if err != nil {
			log.Print(err)
		}

		resBuf := make([]byte, 512)
		_, _, err = s.conn.ReadFromUDP(resBuf)
		if err != nil {
			log.Print(err)
		}
		err = m.Unpack(resBuf)
		if err != nil {
			log.Print(err)
		}
		packed, err = m.Pack()
		_, err = s.conn.WriteToUDP(packed, addr)
		if err != nil {
			log.Print(err)
		}
	}
}

func main() {
	forwardIP := net.ParseIP("8.8.8.8")
	forwardPort := 53
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 8090})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	forwards := []forwardServers{
		{
			ip:   forwardIP,
			port: forwardPort,
		},
	}
	dnsServer := NewDNSServer(conn, forwards)
	ctx, cancelF := context.WithCancel(context.Background())
	defer cancelF()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		<-c
		cancelF()
	}()
	dnsServer.handle(ctx)
}
