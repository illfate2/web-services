package main

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

func newInMemoryCache() *inMemoryCache {
	return &inMemoryCache{
		questionToAnswer: make(map[dnsmessage.Question][]dnsmessage.Resource),
	}
}

func (i *inMemoryCache) Set(question dnsmessage.Question, answers []dnsmessage.Resource, ttl time.Duration) error {
	i.questionToAnswer[question] = answers
	return nil
}

var errNotFound = errors.New("err not found")

func (i *inMemoryCache) Get(question dnsmessage.Question) ([]dnsmessage.Resource, error) {
	answers, ok := i.questionToAnswer[question]
	if !ok {
		return nil, errNotFound
	}
	return answers, nil
}
