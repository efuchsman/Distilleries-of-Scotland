package main

import (
	"fmt"
	"log"
	"os"

	"github.com/efuchsman/distilleries_of_scotland/internal/distilleriesdb"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config/config_dev.yml")

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	// Check if the connection string is provided as an environment variable
	if connStr := os.Getenv("CONN_STR"); connStr != "" {
		fmt.Println("Using connection string from environment variable")
		return
	}

	// Get the connection string from the configuration
	connStr := viper.GetString("environment.development.database.connection_string")

	fmt.Println("Connection String:", connStr)
	if connStr == "" {
		log.Fatal("Connection string not found in the configuration.")
	}

	db, err := distilleriesdb.NewDistilleriesDb(connStr)
	if err != nil {
		log.Fatalf("FAILURE OPENING DATABASE CONNECTION: %v", err)
	}

	defer db.Close()

	// Create the Region table
	err = db.CreateRegionTable()
	if err != nil {
		log.Fatalf("FAILURE TO CREATE REGION TABLE: %v", err)
	}
}
