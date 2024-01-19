package distilleries

import (
	"fmt"
	"testing"

	"github.com/efuchsman/distilleries_of_scotland/internal/distilleriesdb"
	"github.com/stretchr/testify/assert"
)

var (
	testDbDistillery = distilleriesdb.Distillery{
		DistilleryName: "Test",
		RegionName:     "test",
		Geo:            "12345, -12345",
		Town:           "test",
		ParentCompany:  "test",
	}
)

func TestAddDistillery(t *testing.T) {
	testCases := []struct {
		description    string
		distillery     *distilleriesdb.Distillery
		dbclient       distilleriesdb.TestClient
		expectedOutput *RegionalDistillery
		expectedErr    error
	}{
		{
			description: "Success: New Region added to the db if it does not already exist",
			distillery:  &testDbDistillery,
			dbclient: distilleriesdb.TestClient{
				CreateDistilleryData: &testDbDistillery,
			},
			expectedOutput: &RegionalDistillery{DistilleryName: "Test", RegionName: "test", Geo: "12345, -12345", Town: "test", ParentCompany: "test"},
			expectedErr:    nil,
		},
		{
			description: "Failure: db returns an error",
			distillery:  &testDbDistillery,
			dbclient: distilleriesdb.TestClient{
				CreateDistilleryErr: distilleriesdb.ErrNoRows,
			},
			expectedOutput: nil,
			expectedErr:    distilleriesdb.ErrNoRows,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Log(tc.description)
			t.Parallel()
			c := NewDistilleriesClient(tc.dbclient)

			newDistillery, err := c.AddDistillery(tc.distillery.DistilleryName, tc.distillery.RegionName, tc.distillery.Geo, tc.distillery.Town, tc.distillery.ParentCompany)
			if tc.expectedErr != nil {
				assert.Error(t, err, tc.description)
				return
			}

			assert.NoError(t, err, tc.description)
			assert.Equal(t, tc.expectedOutput, newDistillery)
		})
	}
}

func TestGetRegionalDistilleries(t *testing.T) {
	testCases := []struct {
		description    string
		regionName     string
		dbclient       distilleriesdb.TestClient
		expectedOutput *RegionalDistillery
		expectedErr    error
	}{
		{
			description: "Success: New Region added to the db if it does not already exist",
			regionName:  testDbDistillery.RegionName,
			dbclient: distilleriesdb.TestClient{
				GetRegionalDistilleriesData: []distilleriesdb.Distillery{testDbDistillery},
			},
			expectedOutput: &RegionalDistillery{DistilleryName: "Test", RegionName: "test", Geo: "12345, -12345", Town: "test", ParentCompany: "test"},
			expectedErr:    nil,
		},
		{
			description: "Failure: db returns an error",
			regionName:  testDbDistillery.RegionName,
			dbclient: distilleriesdb.TestClient{
				GetRegionalDistilleriesErr: distilleriesdb.ErrNoRows,
			},
			expectedOutput: nil,
			expectedErr:    distilleriesdb.ErrNoRows,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Log(tc.description)
			t.Parallel()
			c := NewDistilleriesClient(tc.dbclient)

			newDistillery, err := c.GetRegionalDistilleries(tc.regionName)
			if tc.expectedErr != nil {
				assert.Error(t, err, tc.description)
				return
			}

			assert.NoError(t, err, tc.description)
			assert.Equal(t, tc.expectedOutput, newDistillery)
		})
	}
}
