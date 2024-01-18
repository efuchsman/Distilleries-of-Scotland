package distilleriesdb

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Distillery struct {
	DistilleryName string `json:"distillery_name"`
	RegionName     string `json:"region_name"`
	Geo            string `json:"geo"`
	Town           string `json:"town"`
	ParentCompany  string `json:"parent_company"`
}

func (db *DistilleriesDB) CreateDistilleriesTable() error {
	query := `
			CREATE TABLE IF NOT EXISTS Distilleries (
					distillery_id SERIAL PRIMARY KEY,
					distillery_name VARCHAR(255) NOT NULL,
					region_name VARCHAR(255) REFERENCES Regions(region_name),
					geo VARCHAR(255),
					town VARCHAR(255),
					parent_company VARCHAR(255),
					UNIQUE (distillery_name)
			);`

	_, err := db.Conn.Exec(query)
	if err != nil {
		log.Errorf("failed to create Distilleries table: %v", err)
		return err
	}

	fmt.Println("Distilleries table created successfully.")
	return nil
}

func (db *DistilleriesDB) CreateDistillery(distilleryName, regionName, geo, town, parentCompany string) (*Distillery, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	res, err := db.GetDistilleryByName(distilleryName)
	if err != nil {
		if err != ErrNoRows {
			log.Errorf("error checking distillery existence: %v", err)
			return nil, err
		}
	}
	if res != nil {
		log.Errorf("distillery %+v already exists", res)
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

	dis := &Distillery{
		DistilleryName: distilleryName,
		RegionName:     regionName,
		Geo:            geo,
		Town:           town,
		ParentCompany:  parentCompany,
	}

	newDistillery, err := db.createDistilleryTx(tx, dis)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("failed to commit transaction: %v", err)
		return nil, err
	}

	return newDistillery, nil
}

func (db *DistilleriesDB) createDistilleryTx(tx *sql.Tx, d *Distillery) (*Distillery, error) {

	query := `
	INSERT INTO Distilleries (distillery_name, region_name, geo, town, parent_company)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (distillery_name) DO NOTHING
			RETURNING distillery_id, distillery_name, region_name, geo, town, parent_company;`

	var distilleryId int
	var distilleryName, regionName, geo, town, parentCompany string

	err := tx.QueryRow(query, d.DistilleryName, d.RegionName, d.Geo, d.Town, d.ParentCompany).Scan(&distilleryId, &distilleryName, &regionName, &geo, &town, &parentCompany)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("failed to insert transactional distillery: %v", err)
		return nil, err
	}

	newDistillery := &Distillery{
		DistilleryName: d.DistilleryName,
		RegionName:     d.RegionName,
		Geo:            d.Geo,
		Town:           d.Town,
		ParentCompany:  d.ParentCompany,
	}

	fmt.Printf("distillery inserted successfully with ID: %d\n", distilleryId)
	return newDistillery, nil
}

// getRegionByName retrieves a region by its name
func (db *DistilleriesDB) GetDistilleryByName(distilleryName string) (*Distillery, error) {
	query := `
		SELECT distillery_name, region_name, geo, town, parent_company
		FROM Distilleries
		WHERE distillery_name = $1;`

	var distillery Distillery
	err := db.Conn.QueryRow(query, distilleryName).Scan(&distillery.DistilleryName, &distillery.RegionName, &distillery.Geo, &distillery.Town, &distillery.ParentCompany)
	if err != nil {
		log.Errorf("failed to get distillery by name: %v", err)
		return nil, ErrNoRows
	}

	return &distillery, nil
}
