package regions

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/efuchsman/distilleries_of_scotland/internal/distilleries"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	projectRoot = ".."
	testRegion1 = &distilleries.Region{
		RegionName:  "test1",
		Description: "test1",
	}
	testRegion2 = &distilleries.Region{
		RegionName:  "test2",
		Description: "test2",
	}

	testRegions = &distilleries.Regions{
		Regions: []*distilleries.Region{
			testRegion1,
			testRegion2,
		},
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

func TestGetRegions(t *testing.T) {
	testCases := []struct {
		description  string
		disClient    *distilleries.TestClient
		expectedCode int
		expectedBody string
	}{
		{
			description: "Success: Region is returned",
			disClient: &distilleries.TestClient{
				GetRegionsData: testRegions,
			},
			expectedCode: 200,
			expectedBody: `{"regions":[{"region_name":"test1","description":"test1"},{"region_name":"test2","description":"test2"}]}` + "\n",
		},
		{
			description: "Error: Internal Error",
			disClient: &distilleries.TestClient{
				GetRegionsErr: errors.New("Internal Error"),
			},
			expectedCode: 500,
			expectedBody: `{"message":"INTERNAL_ERROR","resource":"region","description":"An internal error occurred."}` + "\n",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Log(tc.description)
			t.Parallel()

			h := NewHandler(tc.disClient)

			r := httptest.NewRequest("GET", "/regions", nil)
			w := httptest.NewRecorder()
			h.GetRegions(w, r)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestGetRegion(t *testing.T) {
	testCases := []struct {
		description  string
		disClient    *distilleries.TestClient
		expectedCode int
		expectedBody string
	}{
		{
			description: "Success: Region is returned",
			disClient: &distilleries.TestClient{
				GetRegionByNameData: testRegion1,
			},
			expectedCode: 200,
			expectedBody: `{"region_name":"test1","description":"test1"}` + "\n",
		},
		{
			description: "Error: Region is cannot be found",
			disClient: &distilleries.TestClient{
				GetRegionByNameErr: errors.New("Error fetching region"),
			},
			expectedCode: 404,
			expectedBody: `{"message":"NOT_FOUND","resource":"region","description":"What you are looking for cannot be found."}` + "\n",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Log(tc.description)
			t.Parallel()

			h := NewHandler(tc.disClient)

			r := httptest.NewRequest("GET", "/regions/test1", nil)
			w := httptest.NewRecorder()
			h.GetRegion(w, r)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}
