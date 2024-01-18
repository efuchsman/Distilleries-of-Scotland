package distilleries

type TestClient struct {
	GetRegionByNameData *Region
	GetRegionByNameErr  error

	GetRegionsData *Regions
	GetRegionsErr  error
}

func (c TestClient) GetRegionByName(regionName string) (*Region, error) {
	return c.GetRegionByNameData, c.GetRegionByNameErr
}

func (c TestClient) GetRegions() (*Regions, error) {
	return c.GetRegionsData, c.GetRegionsErr
}
