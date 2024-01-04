package distilleriesdb

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var connStr string

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Error getting the current file path.")
	}

	dir := filepath.Dir(filename)
	projectRoot := filepath.Join(dir, "..", "..")
	configPath := filepath.Join(projectRoot, "config", "config_test.yml")

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading testing configuration file: " + err.Error())
	}

	connStr = viper.GetString("environment.test.database.connection_string")
}

func TestCreateRegionTable(t *testing.T) {
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	defer db.Close()

	distilleriesDB := DistilleriesDB{Conn: db}
	err = distilleriesDB.CreateRegionsTable()
	require.NoError(t, err)

	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_name='regions';")
	require.NoError(t, err)

	defer rows.Close()

	var tableName string
	for rows.Next() {
		err := rows.Scan(&tableName)
		require.NoError(t, err)
	}

	assert.Equal(t, "regions", tableName, "Expected table 'regions' to be created")
}

func TestCreateRegion(t *testing.T) {
	testCases := []struct {
		description     string
		testRegionName  string
		testDescription string
		expectedOutput  *Region
		expectedErr     error
	}{
		{
			description:     "Success: Region added to the DB",
			testRegionName:  "testRegion",
			testDescription: "test description",
			expectedOutput: &Region{
				RegionName:  "testRegion",
				Description: "test description",
			},
			expectedErr: nil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			db, err := NewDistilleriesDb(connStr)
			require.NoError(t, err)

			defer db.Close()

			reg, err := db.GetOrCreateRegion(tc.testRegionName, tc.testDescription)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NotNil(t, reg)
				assert.NoError(t, err, tc.description)
				assert.Equal(t, tc.expectedOutput, reg)
				require.NoError(t, err)
			}
		})
	}
}

func TestGetRegionByName(t *testing.T) {
	testCases := []struct {
		description     string
		testRegionName  string
		testDescription string
		expectedOutput  *Region
		expectedErr     error
	}{
		{
			description:     "Success: Region added to the DB",
			testRegionName:  "testRegion",
			testDescription: "test description",
			expectedOutput: &Region{
				RegionName:  "testRegion",
				Description: "test description",
			},
		},
		{
			description:    "Failure: cannot find region",
			testRegionName: "not found",
			expectedOutput: nil,
			expectedErr:    ErrNoRows,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			db, err := NewDistilleriesDb(connStr)
			require.NoError(t, err)

			defer db.Close()

			reg, err := db.GetRegionByName(tc.testRegionName)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NotNil(t, reg)
				assert.NoError(t, err, tc.description)
				assert.Equal(t, tc.expectedOutput, reg)
				require.NoError(t, err)
			}
		})
	}
}
