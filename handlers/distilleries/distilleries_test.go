package distilleries

import (
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/efuchsman/distilleries_of_scotland/internal/distilleries"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	projectRoot = ".."
	testRegion1 = &distilleries.Region{
		RegionName:  "test1",
		Description: "test1",
	}
	distillery1 = &distilleries.RegionalDistillery{
		DistilleryName: "test1",
		RegionName:     testRegion1.RegionName,
		Geo:            "1234, -1234",
		Town:           "test",
		ParentCompany:  "test",
	}
	distillery2 = &distilleries.RegionalDistillery{
		DistilleryName: "test2",
		RegionName:     testRegion1.RegionName,
		Geo:            "12345, -12345",
		Town:           "test",
		ParentCompany:  "test",
	}
)

func init() {
	_, currentFile, _, ok := runtime.Caller(0)
	if ok {
		projectRoot = filepath.Join(filepath.Dir(currentFile), "..", "..")
	}

	configPath := filepath.Join(projectRoot, "config", "config_test.yml")
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading testing configuration file: " + err.Error())
	}
}
func TestGetRegionalDistilleries(t *testing.T) {
	testCases := []struct {
		description  string
		disClient    *distilleries.TestClient
		regionName   string
		expectedCode int
		expectedBody string
	}{
		{
			description: "Success: Region is returned",
			regionName:  testRegion1.RegionName,
			disClient: &distilleries.TestClient{
				GetRegionByNameData: testRegion1,

				GetRegionalDistilleriesData: &distilleries.RegionalDistilleries{
					RegionalDistilleries: []*distilleries.RegionalDistillery{
						distillery1,
						distillery2,
					},
				},
			},
			expectedCode: 200,
			expectedBody: `{"distilleries":[{"distillery_name":"test1","region_name":"test1","geo":"1234, -1234","town":"test","parent_company":"test"},{"distillery_name":"test2","region_name":"test1","geo":"12345, -12345","town":"test","parent_company":"test"}]}` + "\n",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Log(tc.description)
			t.Parallel()

			h := NewHandler(tc.disClient)

			url := fmt.Sprintf("/regions/%s/distilleries", tc.regionName)
			r := httptest.NewRequest("GET", url, nil)
			r = mux.SetURLVars(r, map[string]string{"region_name": tc.regionName})
			w := httptest.NewRecorder()
			h.GetRegionalDistilleries(w, r)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}
