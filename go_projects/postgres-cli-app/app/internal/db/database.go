package db

import (
	"database/sql"
	"fmt"
	"mycli/internal/models"

	"github.com/lib/pq"
)

type DataBase struct {
	conn *sql.DB
}

func NewClient(connStr string) (*DataBase, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}
	return &DataBase{conn: db}, nil
}

func (db *DataBase) Close() {
	if db.conn != nil {
		db.conn.Close()
	}
}

func (db *DataBase) UpsertConfig(cfg models.Config) (bool, error) {
	query := `
		INSERT INTO postgres.configs (name, description, version, author, tags)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (name) DO UPDATE 
		SET description = EXCLUDED.description,
		    version = EXCLUDED.version,
		    author = EXCLUDED.author,
		    tags = EXCLUDED.tags
		RETURNING (xmax = 0) AS inserted; 
	`
	var inserted bool
	err := db.conn.QueryRow(query,
		cfg.Name,
		cfg.Description,
		cfg.Version,
		cfg.Author,
		pq.Array(cfg.Tags),
	).Scan(&inserted)

	if err != nil {
		return false, fmt.Errorf("failed to upsert config %s: %w", cfg.Name, err)
	}
	return inserted, nil
}
