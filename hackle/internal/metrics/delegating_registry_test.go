package metrics

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/concurrent"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestDelegatingRegistry_NewCounter(t *testing.T) {
	counter := NewDelegatingRegistry().
		Counter("counter", Tags{})
	assert.IsType(t, &DelegatingCounter{}, counter)
}

func TestDelegatingRegistry_NewTimer(t *testing.T) {
	timer := NewDelegatingRegistry().Timer("timer", Tags{})
	assert.IsType(t, &DelegatingTimer{}, timer)
}

func TestDelegatingRegistry_Add(t *testing.T) {

	t.Run("DelegatingRegistry should NOT be added", func(t *testing.T) {
		sut := NewDelegatingRegistry()

		sut.Add(NewDelegatingRegistry())
		sut.Add(NewDelegatingRegistry())
		sut.Add(NewDelegatingRegistry())
		sut.Add(NewDelegatingRegistry())
		sut.Add(NewDelegatingRegistry())

		counter := sut.Counter("counter", Tags{})
		counter.Increment(42)

		assert.Equal(t, int64(0), counter.Count())
	})

	t.Run("already added registry shoud NOT be added", func(t *testing.T) {
		sut := NewDelegatingRegistry()
		registry := NewCumulativeRegistry()

		sut.Add(registry)
		sut.Add(registry)
		sut.Add(registry)
		sut.Add(registry)
		sut.Add(registry)

		counter := sut.Counter("counter", Tags{})
		counter.Increment(42)

		assert.Equal(t, int64(42), counter.Count())
	})

	t.Run("metric before registry add", func(t *testing.T) {
		sut := NewDelegatingRegistry()
		delegatingCounter := sut.Counter("counter", Tags{})
		delegatingCounter.Increment(1)

		assert.Equal(t, int64(0), delegatingCounter.Count())

		cumulative := NewCumulativeRegistry()
		sut.Add(cumulative)

		delegatingCounter.Increment(42)

		assert.Equal(t, int64(42), delegatingCounter.Count())
		assert.Equal(t, int64(42), cumulative.Counter("counter", Tags{}).Count())
	})

	t.Run("registry before metric add", func(t *testing.T) {
		sut := NewDelegatingRegistry()
		cumulative := NewCumulativeRegistry()
		sut.Add(cumulative)

		sut.Counter("counter", Tags{}).Increment(42)

		assert.Equal(t, int64(42), cumulative.Counter("counter", Tags{}).Count())
	})
}

func TestDelegatingRegistry_Close(t *testing.T) {
	sut := NewDelegatingRegistry()
	registry1 := newMockRegistry()
	registry2 := newMockRegistry()

	sut.Add(registry1)
	sut.Add(registry2)

	assert.Equal(t, false, registry1.IsClosed())
	assert.Equal(t, false, registry2.IsClosed())

	sut.Close()

	assert.Equal(t, true, registry1.IsClosed())
	assert.Equal(t, true, registry2.IsClosed())
}

func TestDelegatingRegistry_concurrency(t *testing.T) {

	c := 256

	registry := NewDelegatingRegistry()

	cumulativeRegistries := make([]*CumulativeRegistry, c)
	for i := 0; i < c; i++ {
		cumulativeRegistries[i] = NewCumulativeRegistry()
	}

	task := func(n int) {
		for i := 0; i < c; i++ {
			index := (i / 2) + ((c / 2) * n)
			if i%2 == 0 {
				registry.Counter(strconv.Itoa(index), Tags{})
			} else {
				registry.Add(cumulativeRegistries[index])
			}
		}
	}

	executor := concurrent.NewExecutor()
	for i := 0; i < 2; i++ {
		n := i
		executor.Go(func() {
			task(n)
		})
	}
	executor.Wait()

	assert.Equal(t, c, len(registry.Metrics()))
	for i := 0; i < c; i += 2 {
		name := strconv.Itoa(i)
		registry.Counter(name, Tags{}).Increment(1)
		count := int64(0)
		for _, cumulativeRegistry := range cumulativeRegistries {
			count += cumulativeRegistry.Counter(name, Tags{}).Count()
		}
		assert.Equal(t, int64(c), count)
	}
}
