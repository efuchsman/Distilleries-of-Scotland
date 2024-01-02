package distilleriesdb

import (
	"database/sql"
	"fmt"
)

type DistilleriesDB struct {
	Conn *sql.DB
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
	return &DistilleriesDB{Conn: db}, nil
}

// Close closes the database connection
func (db *DistilleriesDB) Close() {
	db.Conn.Close()
	fmt.Println("Closed the database connection")
}
