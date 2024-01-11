package distilleries

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/efuchsman/distilleries_of_scotland/internal/distilleriesdb"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testDbRegion = &distilleriesdb.Region{
		RegionName:  "Test Region",
		Description: "Test Description",
	}
)

var projectRoot = ".."

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

func TestAddRegion(t *testing.T) {
	testCases := []struct {
		description       string
		dbclient          distilleriesdb.TestClient
		regionName        string
		regionDescription string
		expectedOutput    *Region
		expectedErr       error
	}{
		{
			description:       "Success: New Region added to the db if it does not already exist",
			regionName:        "Test Region",
			regionDescription: "Test Description",
			dbclient: distilleriesdb.TestClient{
				GetOrCreateRegionData: testDbRegion,
			},
			expectedOutput: &Region{RegionName: "Test Region", Description: "Test Description"},
			expectedErr:    nil,
		},
		{
			description:       "Failure: db returns an error",
			regionName:        "Test Region",
			regionDescription: "Test Description",
			dbclient: distilleriesdb.TestClient{
				GetOrCreateRegionErr: distilleriesdb.ErrNoRows,
			},
			expectedErr: distilleriesdb.ErrNoRows,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Log(tc.description)
			t.Parallel()
			c := NewClient(tc.dbclient)

			newRegion, err := c.AddRegion(tc.regionName, tc.regionDescription)
			if tc.expectedErr != nil {
				assert.Error(t, err, tc.description)
				return
			}

			assert.NoError(t, err, tc.description)
			assert.Equal(t, tc.expectedOutput, newRegion)
		})
	}
}

func TestBuildRegions(t *testing.T) {
	testCases := []struct {
		description    string
		filePath       string
		expectedOutput []*Region
		expectedCount  int
		expectedErr    error
	}{
		{
			description:    "Success: Regions slice is created from json file",
			filePath:       filepath.Join(projectRoot, "data/mocks/mock_regions_good.json"),
			expectedOutput: []*Region{},
			expectedCount:  2,
			expectedErr:    nil,
		},
		{
			description:    "Failure: Bad Json",
			filePath:       filepath.Join(projectRoot, "data/mocks/mock_bad_json.json"),
			expectedOutput: []*Region{},
			expectedCount:  0,
			expectedErr:    errors.New(""),
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Log(tc.description)
			t.Parallel()

			regions, err := buildRegions(tc.filePath)
			if tc.expectedErr != nil {
				assert.Error(t, err, tc.description)
				return
			}

			assert.NoError(t, err, tc.description)
			assert.Equal(t, tc.expectedCount, len(regions))
		})
	}
}

func TestGetRegionByName(t *testing.T) {
	testCases := []struct {
		description    string
		regionName     string
		dbclient       distilleriesdb.TestClient
		expectedOutput *Region
		expectedErr    error
	}{
		{
			description: "Success: Region is returned",
			regionName:  "Test Region",
			dbclient: distilleriesdb.TestClient{
				GetRegionByNameData: testDbRegion,
			},
			expectedOutput: &Region{RegionName: "Test Region", Description: "Test Description"},
			expectedErr:    nil,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Log(tc.description)
			t.Parallel()

			c := NewClient(tc.dbclient)

			newRegion, err := c.GetRegionByName(tc.regionName)
			if tc.expectedErr != nil {
				assert.Error(t, err, tc.description)
				return
			}

			assert.NoError(t, err, tc.description)
			assert.Equal(t, tc.expectedOutput, newRegion)
		})
	}
}
