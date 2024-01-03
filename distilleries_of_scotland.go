package main

import (
	"fmt"
	"log"
	"os"

	"github.com/efuchsman/distilleries_of_scotland/internal/distilleriesdb"
)

func main() {
	connStr := os.Getenv("CONN_STR")
	fmt.Println("Connection String:", connStr)

	if connStr == "" {
		log.Fatal("CONN_STR environment variable is not set.")
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
