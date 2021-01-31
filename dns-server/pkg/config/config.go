package config

import (
	"io/ioutil"
	"net"

	"github.com/jszwec/csvutil"
	"golang.org/x/net/dns/dnsmessage"
)

func GetConfigResources() map[dnsmessage.Question]*Resources {
	file, err := ioutil.ReadFile("/home/illfate/go/src/github.com/illfate2/web-services/dns-server/config-example/config.csv")
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
			resConfig[c.dnsQuestion()] = &Resources{TTL: c.TTL}
		}
		resConfig[c.dnsQuestion()].List = append(resConfig[c.dnsQuestion()].List, c.dnsResource())
	}
	return resConfig
}

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
	List []dnsmessage.Resource
	TTL  int
}
