package memory

import (
	"math"
	"math/rand"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemory(t *testing.T) {
	inputCount := 1_000_000
	testInput := make(map[int]int, inputCount)

	for i := 0; i < inputCount; i++ {
		testInput[i] = rand.Intn(math.MaxInt)
	}

	m := New[int, int]()

	for i := 0; i < inputCount; i++ {
		testValue := testInput[i]
		m.Set(i, testValue)
	}

	for i := 0; i < inputCount; i++ {
		expectedValue := testInput[i]
		require.Equal(t, expectedValue, m.Get(i))
		_, ok := m.GetWithCheck(i)
		require.True(t, ok)
	}

	for i := 0; i < inputCount; i++ {
		if i%2 == 0 {
			m.Delete(i)
		}
	}

	for i := 0; i < inputCount; i++ {
		if i%2 == 0 {
			expectedValue := 0
			require.Equal(t, expectedValue, m.Get(i))
			_, ok := m.GetWithCheck(i)
			require.False(t, ok)
		} else {
			expectedValue := testInput[i]
			require.Equal(t, expectedValue, m.Get(i))
			_, ok := m.GetWithCheck(i)
			require.True(t, ok)
		}
	}
}

func TestMemoryParallel(t *testing.T) {
	inputCount := 1_000_000
	testInput := make(map[int]int, inputCount)

	for i := 0; i < inputCount; i++ {
		testInput[i] = rand.Intn(math.MaxInt)
	}

	m := New[int, int]()

	poolSize := runtime.NumCPU()
	workManager := make(chan struct{}, poolSize)

	var wg sync.WaitGroup
	for i := 0; i < inputCount; i++ {
		workManager <- struct{}{}
		wg.Add(1)
		go func(i int) {
			testValue := testInput[i]
			m.Set(i, testValue)
			wg.Done()
			<-workManager
		}(i)
	}
	wg.Wait()

	for i := 0; i < inputCount; i++ {
		workManager <- struct{}{}
		wg.Add(1)
		go func(i int) {
			expectedValue := testInput[i]
			require.Equal(t, expectedValue, m.Get(i))
			_, ok := m.GetWithCheck(i)
			require.True(t, ok)
			wg.Done()
			<-workManager
		}(i)
	}
	wg.Wait()
}

func BenchmarkSetGetDelete(b *testing.B) {
	inputCount := 1_000_000
	testInput := make(map[int]int, inputCount)

	for i := 0; i < inputCount; i++ {
		testInput[i] = rand.Intn(math.MaxInt)
	}

	m := New[int, int]()

	b.ResetTimer()

	for k := 0; k < b.N; k++ {
		m.Set(k, k)
		m.Get(k)
		m.Delete(k)
	}
}
