package distilleries

type TestClient struct {
	GetRegionByNameData *Region
	GetRegionByNameErr  error
}

func (c TestClient) GetRegionByName(regionName string) (*Region, error) {
	return c.GetRegionByNameData, c.GetRegionByNameErr
}
