package resolver

import (
	"log"
	"net"
	"time"

	"golang.org/x/net/dns/dnsmessage"

	"github.com/illfate2/web-services/dns-server/pkg/cache"
	"github.com/illfate2/web-services/dns-server/pkg/dns"
)

type UDPCacheResolver struct {
	udpResolver dns.Resolver
	cache       cache.Cache
}

func NewUDPCacheResolver(udpResolver dns.Resolver, cache cache.Cache) *UDPCacheResolver {
	return &UDPCacheResolver{udpResolver: udpResolver, cache: cache}
}

func (r *UDPCacheResolver) ResolveDNS(msg dnsmessage.Message) (dnsmessage.Message, error) {
	question := msg.Questions[0]
	answers, err := r.cache.Get(question)
	msg.Authoritative = true
	if err == nil {
		msg.Response = true
		msg.Answers = answers
		return msg, nil
	}
	if err != cache.ErrNotFound {
		return dnsmessage.Message{}, err
	}
	resolvedMsg, err := r.udpResolver.ResolveDNS(msg)
	if err != nil {
		return dnsmessage.Message{}, err
	}
	err = r.cache.Set(question, resolvedMsg.Answers, time.Second*time.Duration(resolvedMsg.Answers[0].Header.TTL))
	log.Print(err)
	return resolvedMsg, nil
}

type UDPResolver struct {
	conn           net.Conn
	defaultBufSize int
}

func NewUDPResolver(conn net.Conn) *UDPResolver {
	return &UDPResolver{
		conn:           conn,
		defaultBufSize: 512,
	}
}

func (r *UDPResolver) ResolveDNS(msg dnsmessage.Message) (dnsmessage.Message, error) {
	packedMsg, err := msg.Pack()
	if err != nil {
		return dnsmessage.Message{}, err
	}
	_, err = r.conn.Write(packedMsg)
	if err != nil {
		return dnsmessage.Message{}, err
	}
	resBuf := make([]byte, r.defaultBufSize)
	_, err = r.conn.Read(resBuf)
	if err != nil {
		return dnsmessage.Message{}, err
	}

	var resMsg dnsmessage.Message
	err = resMsg.Unpack(resBuf)
	if err != nil {
		return dnsmessage.Message{}, err
	}
	return resMsg, nil
}
