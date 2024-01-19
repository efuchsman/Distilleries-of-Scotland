package distilleries

import (
	disdb "github.com/efuchsman/distilleries_of_scotland/internal/distilleriesdb"
)

type Client interface {
	SeedRegions(filePath string) error
	GetRegions() (*Regions, error)
	GetRegionByName(regionName string) (*Region, error)
	SeedDistilleries(filePath string) error
	GetRegionalDistilleries(regionName string) (*RegionalDistilleries, error)
}

type DistilleriesClient struct {
	db disdb.Client
}

func NewDistilleriesClient(db disdb.Client) *DistilleriesClient {
	return &DistilleriesClient{
		db: db,
	}
}
