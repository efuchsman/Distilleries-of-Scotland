package distilleriesdb

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Distillery struct {
	DistilleryName string `json:"distillery_name"`
	RegionName     string `json:"region_name"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
	Geo            string `json:"geo"`
	City           string `json:"city"`
	ParentCompany  string `json:"parent_company"`
}

func (db *DistilleriesDB) CreateDistilleriesTable() error {
	query := `
			CREATE TABLE IF NOT EXISTS Distilleries (
					distillery_id SERIAL PRIMARY KEY,
					distillery_name VARCHAR(255) NOT NULL,
					region_name VARCHAR(255) REFERENCES Regions(region_name),
					latitude VARCHAR(255),
					longitude VARCHAR(255),
					geo VARCHAR(255),
					city VARCHAR(255),
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
