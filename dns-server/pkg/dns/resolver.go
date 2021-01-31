package dns

import "golang.org/x/net/dns/dnsmessage"

type Resolver interface {
	ResolveDNS(message dnsmessage.Message) (dnsmessage.Message, error)
}
