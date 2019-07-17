package memstore

import (
	"sync"
	"time"
)

// 缓存接口
type Store interface {
	Set(key string, value interface{})

	Get(key string) (value interface{})

	Del(key string)
}

type memoryStore struct {
	sync.RWMutex
	hash       map[string]interface{}
	hashTime   map[string]time.Time
	expiration time.Duration
}

func NewMemoryStore(expiration time.Duration) Store {
	s := new(memoryStore)
	s.hash = make(map[string]interface{})
	s.hashTime = make(map[string]time.Time)
	s.expiration = expiration
	return s
}

func (s *memoryStore) Set(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()

	s.hash[key] = value
	s.hashTime[key] = time.Now()
}

func (s *memoryStore) Get(key string) (value interface{}) {
	s.Lock()
	defer s.Unlock()

	value, ok := s.hash[key]
	if !ok {
		return
	}
	if time.Now().After(s.hashTime[key].Add(s.expiration)) {
		delete(s.hash, key)
		delete(s.hashTime, key)
		return nil
	}
	return
}

func (s *memoryStore) Del(key string) {
	s.Lock()
	defer s.Unlock()

	delete(s.hash, key)
	delete(s.hashTime, key)
}
