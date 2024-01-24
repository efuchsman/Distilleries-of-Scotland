package distilleries

type TestClient struct {
	DistilleriesClient
	GetRegionByNameData *Region
	GetRegionByNameErr  error

	GetRegionsData *Regions
	GetRegionsErr  error

	GetRegionalDistilleriesData *RegionalDistilleries
	GetRegionalDistilleriesErr  error
}

func (c TestClient) GetRegionByName(regionName string) (*Region, error) {
	return c.GetRegionByNameData, c.GetRegionByNameErr
}

func (c TestClient) GetRegions() (*Regions, error) {
	return c.GetRegionsData, c.GetRegionsErr
}

func (c TestClient) GetRegionalDistilleries(regionName string) (*RegionalDistilleries, error) {
	return c.GetRegionalDistilleriesData, c.GetRegionalDistilleriesErr
}
