package main

import (
	"math/rand"
	"time"
)

func GenerateRandomOrder(orderUID string) Order {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	return Order{
		OrderUID:          orderUID,
		TrackNumber:       randomString(12),
		Entry:             randomString(4),
		Delivery:          generateRandomDelivery(),
		Payment:           generateRandomPayment(),
		Items:             generateRandomItems(),
		Locale:            randomString(2),
		InternalSignature: randomString(8),
		CustomerID:        randomString(8),
		DeliveryService:   randomString(5),
		ShardKey:          randomString(1),
		SmID:              rand.Intn(100),
		DateCreated:       time.Now(),
		OofShard:          randomString(1),
	}
}

func generateRandomDelivery() Delivery {
	return Delivery{
		Name:    randomString(10),
		Phone:   randomString(10),
		Zip:     randomString(6),
		City:    randomString(8),
		Address: randomString(15),
		Region:  randomString(7),
		Email:   randomString(10) + "@example.com",
	}
}

func generateRandomPayment() Payment {
	return Payment{
		Transaction:  randomString(10),
		RequestId:    randomString(8),
		Currency:     "USD",
		Provider:     randomString(6),
		Amount:       rand.Intn(10000),
		PaymentDt:    time.Now().Unix(),
		Bank:         randomString(5),
		DeliveryCost: rand.Intn(500),
		GoodsTotal:   rand.Intn(500),
		CustomFee:    rand.Intn(50),
	}
}

func generateRandomItems() []Item {
	itemCount := rand.Intn(5) + 1
	items := make([]Item, itemCount)

	for i := 0; i < itemCount; i++ {
		items[i] = Item{
			ChrtID:      rand.Intn(10000),
			TrackNumber: randomString(12),
			Price:       rand.Intn(500),
			RID:         randomString(10),
			Name:        randomString(8),
			Sale:        rand.Intn(50),
			Size:        randomString(2),
			TotalPrice:  rand.Intn(500),
			NmID:        rand.Intn(10000),
			Brand:       randomString(10),
			Status:      rand.Intn(300),
		}
	}

	return items
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	s := make([]byte, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
