package config

import (
	"flag"
	"io/ioutil"
	"net"

	"github.com/jszwec/csvutil"
	"github.com/octago/sflags/gen/gflag"
	"golang.org/x/net/dns/dnsmessage"
)

func GetConfigResources(path string) map[dnsmessage.Question]*Resources {
	file, err := ioutil.ReadFile(path)
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
	case dnsmessage.TypeCNAME:
		return &dnsmessage.CNAMEResource{
			CNAME: dnsmessage.MustNewName(record),
		}
	case dnsmessage.TypeNS:
		return &dnsmessage.NSResource{
			NS: dnsmessage.MustNewName(record),
		}
	default:
		return nil
	}
}

func getDNSType(t string) dnsmessage.Type {
	switch t {
	case "A":
		return dnsmessage.TypeA
	case "CNAME":
		return dnsmessage.TypeCNAME
	case "NS":
		return dnsmessage.TypeNS
	default:
		return 0
	}
}

type Resources struct {
	List []dnsmessage.Resource
	TTL  int
}

type CLIConfig struct {
	ForwardAddr      string `flag:"forward f" desc:"address of dns to forward"`
	ServerPort       int    `flag:"port p" desc:"server port"`
	PathToConfigFile string `flag:"config c" desc:"path to config file in csv"`
	Debug            bool   `flag:"debug d" desc:"enable debug mode with output logs"`
}

func MustParseCLIConfig() CLIConfig {
	cfg := CLIConfig{
		ForwardAddr: "8.8.8.8:53",
		ServerPort:  8090,
	}
	err := gflag.ParseToDef(&cfg)
	if err != nil {
		panic(err)
	}
	flag.Parse()
	if cfg.PathToConfigFile == "" {
		panic("path to config should not be empty, type --help to see more")
	}
	return cfg
}
