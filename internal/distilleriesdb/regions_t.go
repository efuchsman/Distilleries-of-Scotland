package distilleriesdb

import (
	"fmt"

	_ "github.com/lib/pq"
)

type Region struct {
	RegionName  string `json: "region_name"`
	Description string `json: "description"`
}

// Create the Region table
func (db *DistilleriesDB) CreateRegionTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS Region (
			region_id SERIAL PRIMARY KEY,
			region_name VARCHAR(255) NOT NULL,
			description TEXT,
			UNIQUE (region_name)
		);`

	_, err := db.Conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create Region table: %v", err)
	}

	fmt.Println("Region table created successfully.")
	return nil
}
