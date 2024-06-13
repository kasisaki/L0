package database

import (
	"L0/internal/models"
	"log"
	"sync"
	"time"
)

var migrationStarted bool

type InMemoryStorage struct {
	mu            sync.RWMutex
	Storage       map[string]models.Order
	BackUpStorage map[string]models.Order
}

func (s *InMemoryStorage) Save(order models.Order) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Storage[order.OrderUID]; ok {
		return false
	}
	s.Storage[order.OrderUID] = order
	return true
}

func (s *InMemoryStorage) SaveToBackup(order models.Order) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.BackUpStorage[order.OrderUID]; ok {
		return false
	}
	s.BackUpStorage[order.OrderUID] = order
	return true
}

func (s *InMemoryStorage) Remove(orderUid string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.Storage, orderUid)
	return true
}

func (s *InMemoryStorage) RemoveFromBackup(orderUID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.BackUpStorage, orderUID)
}

func (s *InMemoryStorage) Get(orderUid string) (models.Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.Storage[orderUid]
	return order, ok
}

func (s *InMemoryStorage) GetAllBackupOrders() ([]models.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]models.Order, 0, len(s.BackUpStorage))
	for _, order := range s.BackUpStorage {
		orders = append(orders, order)
	}
	return orders, nil
}

func NewInMemory() *InMemoryStorage {
	return &InMemoryStorage{
		Storage:       map[string]models.Order{},
		BackUpStorage: map[string]models.Order{},
	}
}

func (s *InMemoryStorage) LoadOrders(orders []models.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, order := range orders {
		s.Storage[order.OrderUID] = order
	}
}

func (s *InMemoryStorage) MigrateBackUpToDB() {
	if migrationStarted {
		return
	}
	migrationStarted = true
	for {
		time.Sleep(30 * time.Second) // Раз в 30 секунд проверяем состояние базы данных

		orders, err := s.GetAllBackupOrders()
		if err != nil {
			continue
		}

		if len(orders) == 0 {
			migrationStarted = false
			log.Println("Backed up data migrated successfully...")
			return
		}

		log.Println("Attempting migration of backed up data to DB...")
		for _, order := range orders {
			err := InsertOrder(Db(), order)
			if err != nil {
				log.Printf("Failed to insert order %s: %v", order.OrderUID, err)
				continue
			}
			s.RemoveFromBackup(order.OrderUID)
		}
	}
}
