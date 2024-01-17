package distilleriesdb

type TestClient struct {
	CreateRegionData *Region
	CreateRegionErr  error

	GetRegionByNameData *Region
	GetRegionByNameErr  error

	GetRegionsData []Region
	GetRegionsErr  error
}

func (c TestClient) CreateRegion(regionName string, description string) (*Region, error) {
	return c.CreateRegionData, c.CreateRegionErr
}

func (c TestClient) GetRegionByName(regionName string) (*Region, error) {
	return c.GetRegionByNameData, c.GetRegionByNameErr
}

func (c TestClient) GetRegions() ([]Region, error) {
	return c.GetRegionsData, c.GetRegionsErr
}
