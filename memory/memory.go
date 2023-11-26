package memory

import (
	"fmt"
	"hash/maphash"
	"maps"
	"sync"
	"time"
)

const (
	memoryCleanerInterval = time.Hour

	countOfShards = 10
)

type Memory[K comparable, V any] struct {
	shardedData [countOfShards]map[K]V

	shardedRWMu [countOfShards]*sync.RWMutex

	hashSeed maphash.Seed

	onceMemoryCleanerProcess *sync.Once
}

func New[K comparable, V any]() *Memory[K, V] {
	var shardedData [countOfShards]map[K]V
	for i := range shardedData {
		shardedData[i] = make(map[K]V)
	}

	var shardedRWMu [countOfShards]*sync.RWMutex
	for i := range shardedRWMu {
		shardedRWMu[i] = new(sync.RWMutex)
	}

	s := &Memory[K, V]{
		shardedData: shardedData,
		shardedRWMu: shardedRWMu,

		hashSeed: maphash.MakeSeed(),

		onceMemoryCleanerProcess: new(sync.Once),
	}

	go s.startMemoryCleanerProcess()

	return s
}

func (m *Memory[K, V]) Set(key K, value V) {
	shardKey := m.shardKey(key)

	shardRWMu := m.shardedRWMu[shardKey]

	withLock(shardRWMu, func() {
		m.shardedData[shardKey][key] = value
	})
}

func (m *Memory[K, V]) Get(key K) (value V) {
	shardKey := m.shardKey(key)

	shardRWMu := m.shardedRWMu[shardKey]

	withRLock(shardRWMu, func() {
		value = m.shardedData[shardKey][key]
	})

	return value
}

func (m *Memory[K, V]) GetWithCheck(key K) (value V, isThere bool) {
	shardKey := m.shardKey(key)

	shardRWMu := m.shardedRWMu[shardKey]

	withRLock(shardRWMu, func() {
		value, isThere = m.shardedData[shardKey][key]
	})

	return value, isThere
}

func (m *Memory[K, V]) Delete(key K) {
	shardKey := m.shardKey(key)

	shardRWMu := m.shardedRWMu[shardKey]

	withLock(shardRWMu, func() {
		delete(m.shardedData[shardKey], key)
	})
}

// startMemoryCleanerProcess cleans memory. It's need, cause map in Go don't release completely memory when deleted.
// see more: https://github.com/golang/go/issues/20135
func (m *Memory[K, V]) startMemoryCleanerProcess() {
	m.onceMemoryCleanerProcess.Do(func() {

		intervalBetweenStartsShardsClean := memoryCleanerInterval / countOfShards

		for i := 0; i < countOfShards; i++ {
			go func(shardKey int) {
				for range time.Tick(memoryCleanerInterval) {
					shardRWMu := m.shardedRWMu[shardKey]

					withLock(shardRWMu, func() {
						tmp := m.shardedData[shardKey]

						m.shardedData[shardKey] = make(map[K]V, len(tmp))

						maps.Copy(m.shardedData[shardKey], tmp)
					})
				}
			}(i)

			// this is necessary so that each shard is cleaned at different times
			<-time.After(intervalBetweenStartsShardsClean)
		}
	})
}

func (m *Memory[K, V]) shardKey(key K) uint64 {
	return maphash.String(m.hashSeed, fmt.Sprint(key)) % countOfShards
}

func withLock(mu *sync.RWMutex, do func()) {
	mu.Lock()
	defer mu.Unlock()
	do()
}

func withRLock(mu *sync.RWMutex, do func()) {
	mu.RLock()
	defer mu.RUnlock()
	do()
}
