package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Consumable represents a serving of a beverage or food with caffeine
type Consumable struct {
	gorm.Model
	Name   string
	Amount uint
}

// Consumption represents a the ingestion of a Consumable
type Consumption struct {
	gorm.Model
	Consumable   Consumable
	ConsumableID uint
}

// EnergyLevel represents your feeling at a time
type EnergyLevel struct {
	gorm.Model
	Consumption   Consumption
	ConsumptionID uint
	Level         uint
}

// MgAtTime represents the amount of caffeine in the body at a particular time
type MgAtTime struct {
	Amount float64   `json:"mg"`
	Time   time.Time `json:"time"`
}

var db gorm.DB

func main() {
	fmt.Println("starting up")

	db = initDb()

	// TODO: Init consumables on the first run
	// dbInitSeedConsumables(db)
	// db.Create(&Consumable{Name: "Small Coffee", Amount: 95})

	initServer()
}

func initDb() gorm.DB {
	dbUsername := os.Getenv("CAFFEINE_DB_USERNAME")
	dbPassword := os.Getenv("CAFFEINE_DB_PASSWORD")
	dbName := os.Getenv("CAFFEINE_DB_NAME")

	dbConnStr := fmt.Sprintf(
		"%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		dbUsername,
		dbPassword,
		dbName)
	db, err := gorm.Open("mysql", dbConnStr)
	if err != nil {
		panic("failed to connect to database")
	}
	dbCreateTables(db)

	return *db
}

func dbCreateTables(db *gorm.DB) {
	if !db.HasTable("consumables") {
		db.CreateTable(&Consumable{})
	}
	if !db.HasTable("consumptions") {
		db.CreateTable(&Consumption{})
	}
	if !db.HasTable("energy_levels") {
		db.CreateTable(&EnergyLevel{})
	}
}
