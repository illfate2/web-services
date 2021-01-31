package dns

import (
	"errors"
	"net"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/dns/dnsmessage"
)

type Server struct {
	resolver       Resolver
	conn           *net.UDPConn
	defaultBufSize int
	logger         *logrus.Logger
}

func NewServer(conn *net.UDPConn, resolver Resolver, logger *logrus.Logger) *Server {
	return &Server{
		resolver:       resolver,
		conn:           conn,
		defaultBufSize: 512,
		logger:         logger,
	}
}

func (s *Server) readDNSMsg() (dnsmessage.Message, *net.UDPAddr, error) {
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

func (s *Server) Handle() {
	for {
		err := s.handleIncomingReq()
		if err != nil {
			s.logger.Warn(err)
		}
	}
}

var (
	errNotSupportedType              = errors.New("not supported type")
	errNotSupportedAmountOfQuestions = errors.New("not supported amount of questions")
)

func (s *Server) handleIncomingReq() error {
	msg, clientAddr, err := s.readDNSMsg()
	if err != nil {
		return err
	}
	if len(msg.Questions) > 1 || len(msg.Questions) == 0 {
		err = errNotSupportedAmountOfQuestions
		s.responseWithErr(clientAddr, msg, err)
		return err
	}
	s.logger.Debug("got questions: ", msg.Questions)
	resolvedMsg, err := s.resolver.ResolveDNS(msg)
	if err != nil {
		s.responseWithErr(clientAddr, msg, err)
		return err
	}
	s.logger.Debug("got answers: ", resolvedMsg.Answers)
	return s.sendDNSMsg(clientAddr, resolvedMsg)
}

func (s *Server) responseWithErr(clientAddr *net.UDPAddr, msg dnsmessage.Message, err error) {
	switch err {
	case errNotSupportedType:
		msg.Header.RCode = dnsmessage.RCodeNotImplemented
	default:
		msg.Header.RCode = dnsmessage.RCodeRefused
	}
	err = s.sendDNSMsg(clientAddr, msg)
	if err != nil {
		s.logger.Warn(err)
	}
}

func (s *Server) sendDNSMsg(addr *net.UDPAddr, message dnsmessage.Message) error {
	packed, err := message.Pack()
	if err != nil {
		return err
	}
	_, err = s.conn.WriteToUDP(packed, addr)
	return err
}
