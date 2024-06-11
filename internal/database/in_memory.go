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

func (s *InMemoryStorage) Get(orderUid string, wg *sync.WaitGroup) models.Order {
	defer wg.Done()

	mu.Lock()
	order := s.Storage[orderUid]
	mu.Unlock()
	return order
}

func (s *InMemoryStorage) New() map[string]models.Order {

}
