package distilleries

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type RegionalDistillery struct {
	DistilleryName string `json:"distillery_name"`
	RegionName     string `json:"region_name"`
	Geo            string `json:"geo"`
	Town           string `json:"town"`
	ParentCompany  string `json:"parent_company"`
}

type RegionalDistilleries struct {
	RegionalDistilleries []*RegionalDistillery `json:"distilleries"`
}

func (c *DistilleriesClient) AddDistillery(distilleryName, regionName, geo, town, parentCompany string) (*RegionalDistillery, error) {
	fields := log.Fields{"Distillery Name": distilleryName, "Region Name": regionName, "Geo": geo, "Town": town, "Parent Company": parentCompany}

	distillery, err := c.db.CreateDistillery(distilleryName, regionName, geo, town, parentCompany)
	if err != nil {
		log.WithFields(fields).Errorf("Error creating region: %+v", err)
		return nil, errors.WithStack(err)
	}
	newDistillery := &RegionalDistillery{
		DistilleryName: distillery.DistilleryName,
		RegionName:     distillery.RegionName,
		Geo:            distillery.Geo,
		Town:           distillery.Town,
		ParentCompany:  distillery.ParentCompany,
	}

	return newDistillery, nil
}

// Helper which builds a regional distilleries slice
func buildDistilleries(filePath string) (*RegionalDistilleries, error) {
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var distilleriesData RegionalDistilleries
	err = json.Unmarshal(jsonData, &distilleriesData.RegionalDistilleries)
	if err != nil {
		log.Errorf("Error building regions slice from file: %v", err)
		return nil, errors.WithStack(err)
	}

	return &distilleriesData, nil
}

func (c *DistilleriesClient) SeedDistilleries(filePath string) error {
	distilleries, err := buildDistilleries(filePath)
	if err != nil {
		log.Errorf("Error seeding Regions: %v", err)
		return errors.WithStack(err)
	}

	for _, dis := range distilleries.RegionalDistilleries {
		resp, err := c.db.GetDistilleryByName(dis.DistilleryName)
		if err != nil {
			log.Errorf("error checking region existence: %v", err)
			return err
		}

		if resp != nil {
			log.Printf("region %+v does not need to be seeded as it already exists", resp)
		} else {
			_, err = c.AddDistillery(dis.DistilleryName, dis.RegionName, dis.Geo, dis.Town, dis.ParentCompany)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}

func (c *DistilleriesClient) GetRegionalDistilleries(regionName string) (*RegionalDistilleries, error) {
	fields := log.Fields{"Region Name": regionName}

	regDis, err := c.db.GetRegionalDistilleries(regionName)
	if err != nil {
		log.WithFields(fields).Errorf("Error fetching regions: %+v", err)
		return nil, errors.WithStack(err)
	}

	regionalDistilleries := &RegionalDistilleries{RegionalDistilleries: make([]*RegionalDistillery, 0)}
	for _, dis := range regDis {
		convertedDis := &RegionalDistillery{
			DistilleryName: dis.DistilleryName,
			RegionName:     dis.RegionName,
			Geo:            dis.Geo,
			Town:           dis.Town,
			ParentCompany:  dis.ParentCompany,
		}
		regionalDistilleries.RegionalDistilleries = append(regionalDistilleries.RegionalDistilleries, convertedDis)
	}
	return regionalDistilleries, nil
}
