package distilleriesdb

type TestClient struct {
	GetOrCreateRegionData *Region
	GetOrCreateRegionErr  error

	GetRegionByNameData *Region
	GetRegionByNameErr  error

	GetRegionsData []Region
	GetRegionsErr  error
}

func (c TestClient) GetOrCreateRegion(regionName string, description string) (*Region, error) {
	return c.GetOrCreateRegionData, c.GetOrCreateRegionErr
}

func (c TestClient) GetRegionByName(regionName string) (*Region, error) {
	return c.GetRegionByNameData, c.GetRegionByNameErr
}

func (c TestClient) GetRegions() ([]Region, error) {
	return c.GetRegionsData, c.GetRegionsErr
}
