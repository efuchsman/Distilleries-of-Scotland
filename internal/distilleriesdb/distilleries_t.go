package distilleriesdb

import (
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
