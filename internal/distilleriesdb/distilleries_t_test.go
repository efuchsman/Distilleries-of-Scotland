package distilleriesdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testRegion = &Region{
		RegionName:  "testRegion",
		Description: "test description",
	}
	testDistillery = &Distillery{
		DistilleryName: "Test Distillery",
		RegionName:     testRegion.RegionName,
		Geo:            "12345, -12345",
		Town:           "test town",
		ParentCompany:  "Strickland Propane",
	}
)

func TestCreateDisilleriesTable(t *testing.T) {
	db, err := NewDistilleriesDb(connStr, true, t.Name()) // Use txdb for testing
	require.NoError(t, err)
	defer db.Close()

	distilleriesDB := DistilleriesDB{Conn: db.Conn}
	err = distilleriesDB.CreateRegionsTable()
	require.NoError(t, err)

	err = distilleriesDB.CreateDistilleriesTable()
	require.NoError(t, err)

	rows, err := db.Conn.Query("SELECT table_name FROM information_schema.tables WHERE table_name='distilleries';")
	require.NoError(t, err)
	defer rows.Close()

	var tableName string
	for rows.Next() {
		err := rows.Scan(&tableName)
		require.NoError(t, err)
	}

	assert.Equal(t, "distilleries", tableName, "Expected table 'distilleries' to be created")
}

func TestCreateDistillery(t *testing.T) {
	testCases := []struct {
		description     string
		region          *Region
		distillery      *Distillery
		testDescription string
		expectedOutput  *Distillery
		expectedErr     error
	}{
		{
			description:    "Success: Region added to the DB",
			region:         testRegion,
			distillery:     testDistillery,
			expectedOutput: testDistillery,
			expectedErr:    nil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			db, err := NewDistilleriesDb(connStr, true, t.Name()) // Use txdb for testing
			require.NoError(t, err)
			defer db.Close()

			err = db.CreateRegionsTable()
			require.NoError(t, err)

			err = db.CreateDistilleriesTable()
			require.NoError(t, err)

			_, err = db.CreateRegion(tc.region.RegionName, tc.region.Description)
			require.NoError(t, err)

			dis, err := db.CreateDistillery(tc.distillery.DistilleryName, tc.distillery.RegionName, tc.distillery.Geo, tc.distillery.Town, tc.distillery.ParentCompany)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NotNil(t, dis)
				assert.NoError(t, err, tc.description)
				assert.Equal(t, tc.expectedOutput, dis)
				require.NoError(t, err)
			}

		})
	}
}
