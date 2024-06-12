package database

import (
	mod "L0/internal/models"
	"database/sql"
	"encoding/json"
)

func InsertOrder(db *sql.DB, order mod.Order) error {
	deliveryJSON, err := json.Marshal(order.Delivery)
	if err != nil {
		return err
	}

	paymentJSON, err := json.Marshal(order.Payment)
	if err != nil {
		return err
	}

	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO orders (
            order_uid, track_number, entry, delivery_info, payment_info, items_info, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
        )
    `

	_, err = db.Exec(query,
		order.OrderUID, order.TrackNumber, order.Entry, deliveryJSON, paymentJSON, itemsJSON,
		order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey,
		order.SmID, order.DateCreated, order.OofShard,
	)
	return err
}

func GetOrderById(db *sql.DB, orderUID string) (*mod.Order, error) {
	query := `SELECT order_uid, track_number, entry, delivery_info, payment_info, items_info, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders WHERE order_uid = $1`

	row := db.QueryRow(query, orderUID)

	var order mod.Order
	var deliveryJSON, paymentJSON, itemsJSON []byte

	err := row.Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &deliveryJSON, &paymentJSON, &itemsJSON,
		&order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(deliveryJSON, &order.Delivery)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(paymentJSON, &order.Payment)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(itemsJSON, &order.Items)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

//func GetAllFromDB(db *sql.DB) ([]mod.Order, error) {
//	query := `SELECT * from orders`
//
//}
