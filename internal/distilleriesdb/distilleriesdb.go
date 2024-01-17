package distilleriesdb

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/lib/pq"
)

type Client interface {
	CreateRegion(regionName string, description string) (*Region, error)
	GetRegionByName(regionName string) (*Region, error)
	GetRegions() ([]Region, error)
}

type DistilleriesDB struct {
	Conn  *sql.DB
	TxDB  bool   // Flag to indicate whether to use txdb (only use for testing)
	TxDrv string // Unique name for txdb registration
	mu    sync.Mutex
}

func NewDistilleriesDb(connStr string, useTxDB bool, TxDrv string) (*DistilleriesDB, error) {
	var db *sql.DB
	var err error

	if useTxDB {
		txdb.Register(TxDrv, "postgres", connStr)
		db, err = sql.Open(TxDrv, "")
	} else {
		db, err = sql.Open("postgres", connStr)
	}

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database")
	return &DistilleriesDB{
		Conn:  db,
		TxDB:  useTxDB,
		TxDrv: TxDrv,
		mu:    sync.Mutex{},
	}, nil
}

func (db *DistilleriesDB) Close() {
	db.Conn.Close()
	fmt.Println("Closed the database connection")
}
