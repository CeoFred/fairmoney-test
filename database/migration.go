package database

import (
	"context"
	"gorm.io/gorm"
	"log"
	"time"
)

const (
	dbTimeout = 30 * time.Second
)

func RunManualMigration(db *gorm.DB) {

	query1 := `CREATE TABLE IF NOT EXISTS accounts (
			id UUID NOT NULL,
			balance DECIMAL(10, 2),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);

			---- INSERT INTO accounts (id, balance) VALUES 
 --- ('018ec333-5c51-7f5d-b3fc-218d742e9a02', 100.00),
 --- ('018ec333-7e25-7d33-95fd-9c10bed5d388', 200.00);
	`

	query2 := `
		CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    amount FLOAT NOT NULL,
    type VARCHAR(255) NOT NULL,
    account_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
	`

	migrationQueries := []string{
		query1,
		query2,
	}

	log.Println("running db migration :::::::::::::")

	_, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	for _, query := range migrationQueries {
		err := db.Exec(query).Error
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("complete db migration")
}
