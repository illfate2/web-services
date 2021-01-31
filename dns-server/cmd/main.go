package main

import (
	"net"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/sirupsen/logrus"

	"github.com/illfate2/web-services/dns-server/pkg/cache"
	"github.com/illfate2/web-services/dns-server/pkg/config"
	"github.com/illfate2/web-services/dns-server/pkg/dns"
	"github.com/illfate2/web-services/dns-server/pkg/resolver"
)

func mustAddConfigToCache(cache cache.Cache, filePath string) {
	conf := config.GetConfigResources(filePath)
	for k, v := range conf {
		err := cache.Set(k, v.List, time.Second*time.Duration(v.TTL))
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	cfg := config.MustParseCLIConfig()
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: cfg.ServerPort})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clientConn, err := net.Dial("udp", cfg.ForwardAddr)
	if err != nil {
		panic(err)
	}
	defer clientConn.Close()

	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	logger := logrus.New()
	if cfg.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}
	udpResolver := resolver.NewUDPResolver(clientConn)
	cache := cache.NewBadgerCache(db)
	mustAddConfigToCache(cache, cfg.PathToConfigFile)
	dnsServer := dns.NewServer(conn, resolver.NewUDPCacheResolver(udpResolver, cache), logger)
	dnsServer.Handle()
}
