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

type Regions struct {
	Regions []*Region `json:"regions"`
}

func (c *DistilleriesClient) AddRegion(regionName string, regionDescription string) (*Region, error) {
	fields := log.Fields{"Region Name": regionName, "Region Description": regionDescription}

	region, err := c.db.CreateRegion(regionName, regionDescription)
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

// Helper which builds a regions slice
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

func (c *DistilleriesClient) SeedRegions(filePath string) error {
	regions, err := buildRegions(filePath)
	if err != nil {
		log.Errorf("Error seeding Regions: %v", err)
		return errors.WithStack(err)
	}

	for _, region := range regions {
		_, err := c.AddRegion(region.RegionName, region.Description)
		if err != nil {
			log.Errorf("error seeding region %+v: %v", region, err)
			return errors.WithStack(err)
		}
	}

	return nil
}

func (c *DistilleriesClient) GetRegionByName(regionName string) (*Region, error) {
	fields := log.Fields{"Region Name": regionName}

	region, err := c.db.GetRegionByName(regionName)
	if err != nil {
		log.WithFields(fields).Errorf("Error fetching region: %+v", err)
		return nil, errors.WithStack(err)
	}

	distilleryRegion := &Region{
		RegionName:  region.RegionName,
		Description: region.Description,
	}

	return distilleryRegion, nil
}

func (c *DistilleriesClient) GetRegions() (*Regions, error) {
	regions, err := c.db.GetRegions()
	if err != nil {
		log.Errorf("Error fetching regions: %+v", err)
		return nil, errors.WithStack(err)
	}

	distilleryRegions := &Regions{Regions: make([]*Region, 0)}

	for _, r := range regions {
		convertedRegion := &Region{
			RegionName:  r.RegionName,
			Description: r.Description,
		}
		distilleryRegions.Regions = append(distilleryRegions.Regions, convertedRegion)
	}

	return distilleryRegions, nil
}
