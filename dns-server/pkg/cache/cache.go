package cache

import (
	"errors"
	"time"

	"github.com/dgraph-io/badger/v3"
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

type BadgerCache struct {
	db *badger.DB
}

func NewBadgerCache(db *badger.DB) *BadgerCache {
	return &BadgerCache{db: db}
}

func (b *BadgerCache) Set(question dnsmessage.Question, answers []dnsmessage.Resource, ttl time.Duration) error {
	questionMsg := dnsmessage.Message{
		Questions: []dnsmessage.Question{question},
	}
	answersMsg := dnsmessage.Message{
		Answers: answers,
	}
	questionBytes, _ := questionMsg.Pack()
	answersBytes, _ := answersMsg.Pack()
	return b.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(questionBytes, answersBytes).WithTTL(ttl)
		err := txn.SetEntry(e)
		return err
	})
}

func (b *BadgerCache) Get(question dnsmessage.Question) ([]dnsmessage.Resource, error) {
	m := dnsmessage.Message{
		Questions: []dnsmessage.Question{question},
	}
	questionBytes, _ := m.Pack()
	err := b.db.View(func(txn *badger.Txn) error {
		answersBytes, err := txn.Get(questionBytes)
		if err == badger.ErrKeyNotFound {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		return answersBytes.Value(func(val []byte) error {
			return m.Unpack(val)
		})
	})
	return m.Answers, err
}
