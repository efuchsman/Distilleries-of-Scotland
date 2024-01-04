package distilleriesdb

import (
	"fmt"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Region struct {
	RegionName  string `json: "region_name"`
	Description string `json: "description"`
}

// Create the Region table
func (db *DistilleriesDB) CreateRegionsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS Regions (
			region_id SERIAL PRIMARY KEY,
			region_name VARCHAR(255) NOT NULL,
			description TEXT,
			UNIQUE (region_name)
		);`

	_, err := db.Conn.Exec(query)
	if err != nil {
		log.Errorf("failed to create Region table: %v", err)
		return err
	}

	fmt.Println("Region table created successfully.")
	return nil
}

// InsertRegion inserts a new region into the Region table
func (db *DistilleriesDB) CreateRegion(regionName string, description string) (*Region, error) {
	query := `
		INSERT INTO Regions (region_name, description)
		VALUES ($1, $2)
		ON CONFLICT (region_name) DO NOTHING
		RETURNING region_id;`

	var regionID int
	err := db.Conn.QueryRow(query, regionName, description).Scan(&regionID)
	if err != nil {
		log.Errorf("failed to insert region: %v", err)
		return nil, err
	}

	newRegion := &Region{
		RegionName:  regionName,
		Description: description,
	}

	fmt.Printf("Region inserted successfully with ID: %d\n", regionID)
	return newRegion, nil
}
