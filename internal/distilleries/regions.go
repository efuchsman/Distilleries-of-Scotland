package distilleries

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Region struct {
	RegionName  string `json:"region_name"`
	Description string `json:"description"`
}

func (c *Client) AddRegion(regionName string, regionDescription string) (*Region, error) {
	fields := log.Fields{"Region Name": regionName, "Region Description": regionDescription}

	region, err := c.db.GetOrCreateRegion(regionName, regionDescription)
	if err != nil {
		log.WithFields(fields).Errorf("Error creating region: %+v", err)
		return nil, errors.WithStack(err)
	}
	distilleryRegion := &Region{
		RegionName:  region.RegionName,
		Description: region.Description,
	}

	return distilleryRegion, nil
}

// Helper which builds a cities slice
func buildRegions(filePath string) ([]*Region, error) {
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var regions []*Region
	err = json.Unmarshal(jsonData, &regions)
	if err != nil {
		log.Errorf("Error building regions slice from file: %v", err)
		return nil, errors.WithStack(err)
	}

	return regions, nil
}

func (c *Client) SeedRegions(filePath string) error {
	regions, err := buildRegions(filePath)
	if err != nil {
		log.Errorf("Error seeding Regions: %v", err)
		return errors.WithStack(err)
	}

	for _, region := range regions {
		_, err := c.AddRegion(region.RegionName, region.Description)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
