package storage

import (
	"maps"
	"sync"
	"time"
)

const memoryCleanerInterval = time.Hour

type Storage struct {
	data map[string]string
	rwmu sync.RWMutex
}

func NewStorage() *Storage {
	s := &Storage{
		data: make(map[string]string),
		rwmu: sync.RWMutex{},
	}

	go s.startMemoryCleanerProcess()

	return s
}

func (s *Storage) Set(key, value string) {
	s.rwmu.Lock()
	s.data[key] = value
	s.rwmu.Unlock()
}

func (s *Storage) Get(key string) string {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()
	return s.data[key]
}

func (s *Storage) Delete(key string) {
	s.rwmu.Lock()
	delete(s.data, key)
	s.rwmu.Unlock()
}

// startMemoryCleanerProcess cleans memory. It's need, cause map in Go don't release completely memory when delete.
// see more: https://github.com/golang/go/issues/20135
func (s *Storage) startMemoryCleanerProcess() {
	for range time.Tick(memoryCleanerInterval) {
		s.rwmu.Lock()
		tmp := s.data
		s.data = make(map[string]string, len(tmp))
		maps.Copy(s.data, tmp)
		s.rwmu.Unlock()
	}
}
