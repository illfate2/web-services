package main

import (
	"context"
	"net"
	"os"
	"os/signal"
)



func main() {
	forwardIP := net.ParseIP("8.8.8.8")
	forwardPort := 53
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 8090})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	forward := forwardServers{
		ip:   forwardIP,
		port: forwardPort,
	}
	dnsServer := NewDNSServer(conn, forward)
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
