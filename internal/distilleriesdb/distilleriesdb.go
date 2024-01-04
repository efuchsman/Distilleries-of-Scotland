package distilleriesdb

import (
	"database/sql"
	"fmt"
	"sync"
)

type Client interface {
	GetOrCreateRegion(regionName string, description string) (*Region, error)
	GetRegionByName(regionName string) (*Region, error)
}

type DistilleriesDB struct {
	Conn *sql.DB
	mu   sync.Mutex
}

func NewDistilleriesDb(connStr string) (*DistilleriesDB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database")
	return &DistilleriesDB{
		Conn: db,
		mu:   sync.Mutex{},
	}, nil
}

// Close closes the database connection
func (db *DistilleriesDB) Close() {
	db.Conn.Close()
	fmt.Println("Closed the database connection")
}
