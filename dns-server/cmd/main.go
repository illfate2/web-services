package main

import (
	"net"
	"time"

	"github.com/illfate2/web-services/dns-server/pkg/cache"
	"github.com/illfate2/web-services/dns-server/pkg/config"
	"github.com/illfate2/web-services/dns-server/pkg/dns"
	resolver2 "github.com/illfate2/web-services/dns-server/pkg/resolver"
)

func mustAddConfigToCache(cache cache.Cache){
	conf := config.GetConfigResources()
	for k, v := range conf {
		err:=cache.Set(k, v.List, time.Second*time.Duration(v.TTL))
		if err != nil {
			panic(err)
		}
	}
}

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

	resolver := resolver2.NewUDPResolver(clientConn)
	cache := cache.NewInMemoryCache()
	mustAddConfigToCache(cache)
	dnsServer := dns.NewDNSServer(conn, resolver2.NewUDPCacheResolver(resolver, cache))
	dnsServer.Handle()
}

