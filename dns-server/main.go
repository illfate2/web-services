package main

import (
	"context"
	"net"
	"os"
	"os/signal"
)

func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 8090})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	clientConn, err := net.Dial("udp", "8.8.8.8"+":"+"53")
	if err != nil {
		panic(err)
	}
	defer clientConn.Close()
	resolver := NewUDPResolver(clientConn)
	dnsServer := NewDNSServer(conn, resolver)
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
