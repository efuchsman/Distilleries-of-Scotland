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

// CreateRegion creates a new region in the Regions table
func (db *DistilleriesDB) CreateRegion(regionName string, description string) (*Region, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	res, err := db.GetRegionByName(regionName)
	if err != nil {
		if err != ErrNoRows {
			log.Errorf("error checking region existence: %v", err)
			return nil, err
		}
	}
	if res != nil {
		log.Errorf("region %+v already exists", res)
		return nil, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		log.Errorf("failed to begin transaction: %v", err)
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != sql.ErrTxDone && err != nil {
			log.Errorf("failed to rollback transaction: %v", err)
		}
	}()

	reg := &Region{
		RegionName:  regionName,
		Description: description,
	}

	newRegion, err := db.createRegionTx(tx, reg)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("failed to commit transaction: %v", err)
		return nil, err
	}

	return newRegion, nil
}

func (db *DistilleriesDB) createRegionTx(tx *sql.Tx, r *Region) (*Region, error) {
	query := `
			INSERT INTO Regions (region_name, description)
			VALUES ($1, $2)
			ON CONFLICT (region_name) DO NOTHING
			RETURNING region_id, region_name, description;`

	var regionID int
	var regionName, description string
	err := tx.QueryRow(query, r.RegionName, r.Description).Scan(&regionID, &regionName, &description)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("failed to insert transactional region: %v", err)
		return nil, err
	}

	newRegion := &Region{
		RegionName:  r.RegionName,
		Description: r.Description,
	}

	fmt.Printf("Region inserted successfully with ID: %d\n", regionID)
	return newRegion, nil
}

// getRegionByName retrieves a region by its name
func (db *DistilleriesDB) GetRegionByName(regionName string) (*Region, error) {

	query := `
		SELECT region_name, description
		FROM Regions
		WHERE LOWER(region_name) = LOWER($1);`

	var region Region
	err := db.Conn.QueryRow(query, regionName).Scan(&region.RegionName, &region.Description)
	if err != nil {
		log.Errorf("failed to get region by name: %v", err)
		return nil, ErrNoRows
	}

	return &region, nil
}

func (db *DistilleriesDB) GetRegions() ([]Region, error) {
	query := `
			SELECT region_name, description
			FROM Regions;`

	rows, err := db.Conn.Query(query)
	if err != nil {
		log.Errorf("failed to get regions: %v", err)
		return nil, err
	}
	defer rows.Close()

	var regions []Region

	for rows.Next() {
		var region Region
		if err := rows.Scan(&region.RegionName, &region.Description); err != nil {
			log.Errorf("failed to scan region row: %v", err)
			return nil, err
		}
		regions = append(regions, region)
	}

	if err := rows.Err(); err != nil {
		log.Errorf("error iterating over regions rows: %v", err)
		return nil, err
	}

	return regions, nil
}
