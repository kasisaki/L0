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

type Service interface {
	Health() map[string]string
	Db() *sql.DB
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
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) Db() *sql.DB {
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
