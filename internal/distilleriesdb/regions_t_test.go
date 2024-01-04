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
	// Get the path to the current file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Error getting the current file path.")
	}

	// Get the directory of the current file
	dir := filepath.Dir(filename)

	// Get the project root directory (two levels up from the current file)
	projectRoot := filepath.Join(dir, "..", "..")

	// Construct the path to the testing configuration file
	configPath := filepath.Join(projectRoot, "config", "config_test.yml")

	// Load the testing configuration file
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

	// Create the Region table
	err = distilleriesDB.CreateRegionsTable()
	require.NoError(t, err)

	// Check if the Region table was created successfully
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

			reg, err := db.CreateRegion(tc.testRegionName, tc.testDescription)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			}
			assert.NotNil(t, reg)
			assert.NoError(t, err, tc.description)
			assert.Equal(t, tc.expectedOutput, reg)

			_, err = db.Conn.Exec("DELETE FROM Regions WHERE region_name = $1", tc.testRegionName)
			require.NoError(t, err)

		})
	}
}
