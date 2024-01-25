package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/efuchsman/distilleries_of_scotland/cors"
	disH "github.com/efuchsman/distilleries_of_scotland/handlers/distilleries"
	"github.com/efuchsman/distilleries_of_scotland/handlers/regions"
	"github.com/efuchsman/distilleries_of_scotland/internal/distilleries"
	"github.com/efuchsman/distilleries_of_scotland/internal/distilleriesdb"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	fmt.Println("Starting the application")

	// Read the configuration file
	viper.SetConfigFile("config/config_dev.yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}

	// Check if the connection string is provided as an environment variable
	connStr := os.Getenv("CONN_STR")
	if connStr != "" {
		fmt.Println("Using connection string from environment variable")
	} else {
		// Get the connection string from the configuration
		connStr = viper.GetString("environment.development.database.connection_string")

		fmt.Println("Connection String:", connStr)
		if connStr == "" {
			log.Fatal("Connection string not found in the configuration.")
		}
	}

	// Establish the database connection
	db, err := distilleriesdb.NewDistilleriesDb(connStr, false, "")
	if err != nil {
		log.Fatalf("FAILURE OPENING DATABASE CONNECTION: %v", err)
	}
	defer db.Close()

	// Create the Regions table
	err = db.CreateRegionsTable()
	if err != nil {
		log.Fatalf("FAILURE TO CREATE REGION TABLE: %v", err)
	}

	// Create the Distilleries table
	err = db.CreateDistilleriesTable()
	if err != nil {
		log.Fatalf("FAILURE TO CREATE DISTILLERIES TABLE: %v", err)
	}

	// Seed regions to the database
	dis := distilleries.NewDistilleriesClient(db)
	regionsFilePath := "data/regions.json"
	if err = dis.SeedRegions(regionsFilePath); err != nil {
		log.Fatalf("Error seeding regions to the database: %v", err)
	}

	distilleriesFilePath := "data/distilleries.json"
	if err = dis.SeedDistilleries(distilleriesFilePath); err != nil {
		log.Fatalf("Error seeding regions to the database: %v", err)
	}

	// Setup the HTTP server and router
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Distilleries of Scotland API")
	})

	distilleryRegions := regions.NewHandler(dis)
	router.HandleFunc("/regions", distilleryRegions.GetRegions).Methods("GET")
	router.HandleFunc("/regions/{region_name}", distilleryRegions.GetRegion).Methods("GET")

	regionalDistilleries := disH.NewHandler(dis)
	router.HandleFunc("/regions/{region_name}/distilleries", regionalDistilleries.GetRegionalDistilleries).Methods("GET")

	handler := cors.SetCORSHeader(router)

	// Start the HTTP server
	port := 8000
	fmt.Printf("Server is running on :%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	fmt.Println("Application started successfully")
}
