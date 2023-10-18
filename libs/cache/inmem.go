package cache

import (
	"container/list"
	"context"
	"fmt"
	"route256/libs/logger"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	cacheHitCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "server",
		Subsystem: "cache",
		Name:      "hits_total",
	})
	cacheErrorTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "server",
		Subsystem: "cache",
		Name:      "error_total",
	})
	cacheRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "server",
		Subsystem: "cache",
		Name:      "requests_total",
	})
)

type InMemCache[Key comparable, Value any] struct {
	locker  *sync.Mutex
	storage map[Key]Item[Value]
	// Размер хранилища
	size int
	// Время жизни ключа
	ttl time.Duration
	// queue
	queue *list.List
}

type Item[T any] struct {
	value  T
	keyPtr *list.Element
	// Последнее время использования
	lastUsedAt time.Time
	// Время протухания
	deadline time.Time
}

func New[Key comparable, Value any](size int, ttl time.Duration) (c *InMemCache[Key, Value]) {
	c = &InMemCache[Key, Value]{
		locker:  &sync.Mutex{},
		ttl:     ttl,
		size:    size,
		storage: make(map[Key]Item[Value], size),
		queue:   list.New(),
	}
	// Запуск очистки в фоне
	go c.deleteByDeadline(context.Background())
	return
}

// Получить из кэша
func (c *InMemCache[Key, Value]) Get(ctx context.Context, key Key) (value Value, ok bool) {
	cacheRequestsTotal.Inc()
	c.locker.Lock()

	item, ok := c.storage[key]
	if !ok || item.deadline.Before(time.Now()) {
		// Ключ не найден или протух
		cacheErrorTotal.Inc()
		c.locker.Unlock()
		return value, false
	}

	cacheHitCount.Inc()

	// В начало списка и обновляем время использования
	item.lastUsedAt = time.Now()
	c.queue.MoveToFront(item.keyPtr)
	c.storage[key] = item
	c.locker.Unlock()

	value = item.value
	logger.Debug(fmt.Sprintf("cache: get key: %v", key))

	return
}

// Положить в кэш
func (c *InMemCache[Key, Value]) Set(ctx context.Context, key Key, value Value) {
	c.locker.Lock()

	// Если ключа нет и кэш полный, освободим место (удалив LRU ключ)
	item, ok := c.storage[key]
	if !ok {
		if c.IsFull(ctx) {
			c.deleteLRU(ctx)
		}
		item = Item[Value]{
			value:  value,
			keyPtr: c.queue.PushFront(key),
		}
		c.storage[key] = item
	} else {
		item.value = value
		c.queue.MoveToFront(item.keyPtr)
	}

	item.lastUsedAt = time.Now()
	item.deadline = time.Now().Add(c.ttl)

	logger.Debug(fmt.Sprintf("cache: set key: %v, new size: %v", key, len(c.storage)))
	c.locker.Unlock()
}

func (c *InMemCache[Key, Value]) IsFull(ctx context.Context) bool {
	return len(c.storage) >= c.size
}

// Удаляем LRU элемент, O(1)
func (c *InMemCache[Key, Value]) deleteLRU(ctx context.Context) {
	back := c.queue.Back()
	key := c.queue.Remove(back).(Key)
	delete(c.storage, key)

	logger.Debug(fmt.Sprintf("cache: deleted LRU key: %v, new size: %v", key, len(c.storage)))
}

// Удалить протухшие ключи
func (c *InMemCache[Key, Value]) deleteByDeadline(ctx context.Context) {
	logger.Debug("cache: run delete by deadline")

	for now := range time.Tick(time.Second) {
		c.locker.Lock()

		for {
			back := c.queue.Back()
			if back != nil && c.storage[back.Value.(Key)].deadline.Before(now) {
				key := c.queue.Remove(back).(Key)
				delete(c.storage, key)
				logger.Debug(fmt.Sprintf("cache: deleted key by deadline: %v, new size: %v", key, len(c.storage)))
			} else {
				break
			}
		}

		c.locker.Unlock()
	}
}
