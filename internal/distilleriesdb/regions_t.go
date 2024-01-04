package distilleriesdb

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Region struct {
	RegionName  string `json:"region_name"`
	Description string `json:"description"`
}

var ErrNoRows = errors.New("sql: no rows in result set")

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

// CreateRegion gets an existing region or creates a new one in the Region table
func (db *DistilleriesDB) GetOrCreateRegion(regionName string, description string) (*Region, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	query := `
        INSERT INTO Regions (region_name, description)
        VALUES ($1, $2)
        ON CONFLICT (region_name) DO NOTHING
        RETURNING region_id;`

	var regionID int
	err := db.Conn.QueryRow(query, regionName, description).Scan(&regionID)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("failed to insert region: %v", err)
		return nil, err
	}

	if err == sql.ErrNoRows {
		// Region already exists, retrieve it
		existingRegion, err := db.GetRegionByName(regionName)
		if err != nil {
			return nil, err
		}

		if existingRegion != nil {
			log.Printf("Region with name '%s' already exists", regionName)
			return existingRegion, nil
		}
	}

	newRegion := &Region{
		RegionName:  regionName,
		Description: description,
	}

	fmt.Printf("Region inserted successfully with ID: %d\n", regionID)
	return newRegion, nil
}

// getRegionByName retrieves a region by its name
func (db *DistilleriesDB) GetRegionByName(regionName string) (*Region, error) {
	query := `
		SELECT region_name, description
		FROM Regions
		WHERE region_name = $1;`

	var region Region
	err := db.Conn.QueryRow(query, regionName).Scan(&region.RegionName, &region.Description)
	if err != nil {
		log.Errorf("failed to get region by name: %v", err)
		return nil, ErrNoRows
	}

	return &region, nil
}
