package dns

import (
	"errors"
	"log"
	"net"

	"golang.org/x/net/dns/dnsmessage"

	resolver2 "github.com/illfate2/web-services/dns-server/pkg/resolver"
)

type DNSServer struct {
	resolver       resolver2.Resolver
	conn           *net.UDPConn
	defaultBufSize int
}

func NewDNSServer(conn *net.UDPConn, resolver resolver2.Resolver) *DNSServer {
	return &DNSServer{
		resolver:       resolver,
		conn:           conn,
		defaultBufSize: 512,
	}
}

func (s *DNSServer) readDNSMsg() (dnsmessage.Message, *net.UDPAddr, error) {
	buf := make([]byte, s.defaultBufSize)
	_, addr, err := s.conn.ReadFromUDP(buf)
	if err != nil {
		return dnsmessage.Message{}, nil, err
	}
	var msg dnsmessage.Message
	err = msg.Unpack(buf)
	if err != nil {
		return dnsmessage.Message{}, nil, err
	}
	return msg, addr, nil
}

func (s *DNSServer) Handle() {
	for {
		err := s.handleIncomingReq()
		if err != nil {
			log.Print(err)
		}
	}
}

var (
	errNotSupportedType              = errors.New("not supported type")
	errNotSupportedAmountOfQuestions = errors.New("not supported amount of questions")
)

func (s *DNSServer) handleIncomingReq() error {
	msg, clientAddr, err := s.readDNSMsg()
	if err != nil {
		return err
	}
	if len(msg.Questions) > 1 || len(msg.Questions) == 0 {
		err = errNotSupportedAmountOfQuestions
		s.responseWithErr(clientAddr, msg, err)
		return err
	}
	resQuestions := make([]dnsmessage.Question, 0, len(msg.Questions))
	for _, q := range msg.Questions {
		if q.Type == dnsmessage.TypeA {
			resQuestions = append(resQuestions, q)
		}
	}
	if len(resQuestions) == 0 {
		err = errNotSupportedType
		s.responseWithErr(clientAddr, msg, err)
		return err
	}
	msg.Questions = resQuestions
	log.Print(msg.Questions)
	resolvedMsg, err := s.resolver.ResolveDNS(msg)
	if err != nil {
		s.responseWithErr(clientAddr, msg, err)
		return err
	}
	log.Print(msg.Answers)
	return s.sendDNSMsg(clientAddr, resolvedMsg)
}

func (s *DNSServer) responseWithErr(clientAddr *net.UDPAddr, msg dnsmessage.Message, err error) {
	switch err {
	case errNotSupportedType:
		msg.Header.RCode = dnsmessage.RCodeNotImplemented
	default:
		msg.Header.RCode = dnsmessage.RCodeRefused
	}
	err = s.sendDNSMsg(clientAddr, msg)
	if err != nil {
		log.Print(err)
	}
}

func (s *DNSServer) sendDNSMsg(addr *net.UDPAddr, message dnsmessage.Message) error {
	packed, err := message.Pack()
	if err != nil {
		return err
	}
	_, err = s.conn.WriteToUDP(packed, addr)
	return err
}
