package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var InMemory = NewInMemory()

type Service interface {
	Health() map[string]string
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}

	err = runSQLScript(db)
	if err != nil {
		fmt.Println("Error running script")
		log.Fatal(err)
	}

	orders, err := FetchAllOrders(Db())
	if err != nil {
		log.Fatalf("Failed to fetch orders from database: %v", err)
	}

	InMemory.LoadOrders(orders)
	log.Println("Loaded orders into in-memory storage")

	return dbInstance
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer func() {
		fmt.Println("Calling context cancel")
		cancel()
	}()
	err := s.db.PingContext(ctx)
	if err != nil {
		log.Printf(fmt.Sprintf("db down: %v", err))
		return map[string]string{
			"message": "DB is down",
		}
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func Db() *sql.DB {
	return dbInstance.db
}

func runSQLScript(db *sql.DB) error {
	sqlBytes, err := os.ReadFile(filepath.Join("internal/database/schema.sql"))
	if err != nil {
		return err
	}

	requests := strings.Split(string(sqlBytes), ";\n")

	for _, request := range requests {
		_, err := db.Exec(request)
		if err != nil {
			log.Fatalf("DATABASE setup completed with an error: %s\n", err)
		}
	}

	if _, err := db.Exec(string(sqlBytes)); err != nil {
		return err
	}

	return nil
}
