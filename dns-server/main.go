package main

import (
	"io/ioutil"
	"net"
	"time"

	"github.com/jszwec/csvutil"
	"golang.org/x/net/dns/dnsmessage"
)

type Config struct {
	Name   string `csv:"name"`
	Type   string `csv:"type"`
	TTL    int    `csv:"ttl"`
	Record string `csv:"record"`
}

func (c Config) dnsQuestion() dnsmessage.Question {
	return dnsmessage.Question{
		Name:  dnsmessage.MustNewName(c.Name),
		Type:  getDNSType(c.Type),
		Class: dnsmessage.ClassINET,
	}
}

func (c Config) dnsResource() dnsmessage.Resource {
	return dnsmessage.Resource{
		Header: dnsmessage.ResourceHeader{
			Name:  dnsmessage.MustNewName("domain.con."),
			Class: dnsmessage.ClassINET,
			TTL:   uint32(c.TTL),
		},
		Body: getDNSResourceBody(getDNSType(c.Type), c.Record),
	}
}

func getDNSResourceBody(t dnsmessage.Type, record string) dnsmessage.ResourceBody {
	switch t {
	case dnsmessage.TypeA:
		ip := net.ParseIP(record)
		return &dnsmessage.AResource{
			A: [4]byte{ip[12], ip[13], ip[14], ip[15]},
		}
	default:
		return nil
	}
}

func getDNSType(t string) dnsmessage.Type {
	switch t {
	case "A":
		return dnsmessage.TypeA
	default:
		return 0
	}
}

type Resources struct {
	list []dnsmessage.Resource
	ttl  int
}

func getConfigResources() map[dnsmessage.Question]*Resources {
	file, err := ioutil.ReadFile("config.csv")
	if err != nil {
		panic(err)
	}
	var cc []Config
	err = csvutil.Unmarshal(file, &cc)
	if err != nil {
		panic(err)
	}
	resConfig := make(map[dnsmessage.Question]*Resources)
	for _, c := range cc {
		if resConfig[c.dnsQuestion()] == nil {
			resConfig[c.dnsQuestion()] = &Resources{ttl: c.TTL}
		}
		resConfig[c.dnsQuestion()].list = append(resConfig[c.dnsQuestion()].list, c.dnsResource())
	}
	return resConfig
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
	resolver := NewUDPResolver(clientConn)
	cache := newInMemoryCache()
	conf := getConfigResources()
	for k, v := range conf {
		cache.Set(k, v.list, time.Second*time.Duration(v.ttl))
	}
	dnsServer := NewDNSServer(conn, NewUDPCacheResolver(resolver, cache))
	dnsServer.handle()
}
