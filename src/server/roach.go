package server

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
	"fmt"
	"pscoin/src/model"
	"os"
)

var (
	COCKROACHDB_URL = os.Getenv("COCKROACHDB_URL")
	CockroachClient *gorm.DB

)

func init() {
	if COCKROACHDB_URL == "" {
		COCKROACHDB_URL = "postgresql://pscoinroach@localhost:26257/pscoin?sslmode=disable"
	}
}

func SetupCockroachDB() *gorm.DB {
	db, err := gorm.Open("postgres", COCKROACHDB_URL)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	// Auto Migrate the schemas
	result := db.AutoMigrate(&model.PeerNode{})
	if result.Error != nil {
		panic(fmt.Sprintf("failed to AutoMigrate database: %v", result.Error))
	}

	return db
}
