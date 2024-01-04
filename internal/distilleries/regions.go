package distilleries

import (
	"encoding/json"
	"io/ioutil"

	disdb "github.com/efuchsman/distilleries_of_scotland/internal/distilleriesdb"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Region struct {
	RegionName  string `json:"region_name"`
	Description string `json:"description"`
}

func (c *Client) AddRegion(newRegion *Region) (*disdb.Region, error) {
	fields := log.Fields{"Region Name": newRegion.RegionName, "Region Description": newRegion.Description}

	region, err := c.db.GetOrCreateRegion(newRegion.RegionName, newRegion.Description)
	if err != nil {
		log.WithFields(fields).Errorf("Error creating region: %+v", err)
	}

	return region, nil
}

// Helper which builds a cities slice
func (c *Client) buildRegions(filePath string) ([]*Region, error) {
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Create a slice to hold the cities
	var regions []*Region

	// Unmarshal the JSON data into the cities slice
	err = json.Unmarshal(jsonData, &regions)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return regions, nil
}

func (c *Client) SeedRegions(filePath string) error {
	regions, err := c.buildRegions(filePath)
	if err != nil {
		log.Errorf("Error seeding Regions: %v", err)
	}

	for _, region := range regions {
		_, err := c.AddRegion(region)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
