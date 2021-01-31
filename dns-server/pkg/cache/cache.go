package cache

import (
	"errors"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

type Cache interface {
	Set(question dnsmessage.Question, answers []dnsmessage.Resource, ttl time.Duration) error
	Get(question dnsmessage.Question) ([]dnsmessage.Resource, error)
}

type inMemoryCache struct {
	questionToAnswer map[dnsmessage.Question][]dnsmessage.Resource
}

func NewInMemoryCache() *inMemoryCache {
	return &inMemoryCache{
		questionToAnswer: make(map[dnsmessage.Question][]dnsmessage.Resource),
	}
}

func (i *inMemoryCache) Set(question dnsmessage.Question, answers []dnsmessage.Resource, ttl time.Duration) error {
	i.questionToAnswer[question] = answers
	return nil
}

var ErrNotFound = errors.New("err not found")

func (i *inMemoryCache) Get(question dnsmessage.Question) ([]dnsmessage.Resource, error) {
	answers, ok := i.questionToAnswer[question]
	if !ok {
		return nil, ErrNotFound
	}
	return answers, nil
}
