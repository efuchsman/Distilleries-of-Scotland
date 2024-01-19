package distilleries

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
