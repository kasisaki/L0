package database

import (
	"L0/internal/models"
	"sync"
)

var mu sync.Mutex

type InMemoryStorage struct {
	Storage map[string]models.Order
}

func (s *InMemoryStorage) Save(order models.Order, wg *sync.WaitGroup) bool {
	defer wg.Done()

	mu.Lock()
	if _, ok := s.Storage[order.OrderUID]; ok {
		return false
	}
	s.Storage[order.OrderUID] = order
	mu.Unlock()
	return true
}

func (s *InMemoryStorage) Remove(orderUid string, wg *sync.WaitGroup) bool {
	defer wg.Done()

	mu.Lock()
	delete(s.Storage, orderUid)
	mu.Unlock()
	return true
}

func (s *InMemoryStorage) Get(orderUid string, wg *sync.WaitGroup) (models.Order, bool) {
	defer wg.Done()

	mu.Lock()
	order, ok := s.Storage[orderUid]
	mu.Unlock()
	return order, ok
}

func NewInMemory() *InMemoryStorage {
	return &InMemoryStorage{Storage: map[string]models.Order{}}
}
