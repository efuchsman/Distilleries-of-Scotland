package distilleries

import (
	disdb "github.com/efuchsman/distilleries_of_scotland/internal/distilleriesdb"
)

type Service interface {
	SeedRegions(filePath string) error
	GetRegions() (*Regions, error)
	GetRegionByName(regionName string) (*Region, error)
}

type Client struct {
	db disdb.Client
}

func NewClient(db disdb.Client) *Client {
	return &Client{
		db: db,
	}
}
