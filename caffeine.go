package main

import (
	"fmt"
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
	ConsumptionID Consumption
	level         uint
}

func main() {
	fmt.Println("starting up")

	db := initDb()

	// TODO: Init consumables on the first run
	// dbInitSeedConsumables(db)
	// db.Create(&Consumable{Name: "Small Coffee", Amount: 95})
	// Demo, drink coffee
	// var coffee Consumable
	// db.Where("name = ?", "Small Coffee").Find(&coffee)
	// ingest(db, coffee)

	// Demo, calculate amount in body
	mg := mgInBody(db)
	fmt.Println(mg)

	duration, _ := time.ParseDuration("1h")
	fmt.Println(duration)
	timeAtDuration := time.Now().Add(duration)
	fmt.Println(timeAtDuration)
	fmt.Println(mgInBodyAtTime(db, timeAtDuration))
}

func initDb() gorm.DB {
	// TODO: Add to env variables
	db, err := gorm.Open("mysql", "root@/caffeine?charset=utf8&parseTime=True&loc=Local")
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
