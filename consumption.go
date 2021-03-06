package main

import (
	"math"
	"time"

	"github.com/jinzhu/gorm"
)

const caffeineHalfLife float64 = 5.7
const daysToTrack int = 2

// Ingest a single Consumable
func ingest(db gorm.DB, c Consumable) Consumption {
	consumption := Consumption{Consumable: c, ConsumableID: c.ID}
	db.Create(&consumption)
	return consumption
}

// Find all caffeine ingestions since a point in time
func ingestionsSince(db gorm.DB, findSince time.Time) []Consumption {
	var consumptions []Consumption
	db.Where("created_at > ?", findSince).Find(&consumptions)
	return consumptions
}

func amountRemainingAtTime(consumption Consumption, time time.Time) float64 {
	mgConsumed := consumption.Consumable.Amount
	consumedAt := consumption.CreatedAt

	timeElapsed := time.Sub(consumedAt)
	hoursElapsed := timeElapsed.Hours()
	if hoursElapsed < 0 {
		return 0
	}

	remaining := float64(mgConsumed) * math.Pow(0.5, (hoursElapsed/float64(caffeineHalfLife)))
	return remaining
}

func mgInBody(db gorm.DB) float64 {
	return mgInBodyAtTime(db, time.Now())
}

func mgInBodyAtTime(db gorm.DB, time time.Time) float64 {
	var totalMg float64

	findSince := time.AddDate(0, 0, -daysToTrack)
	consumptions := ingestionsSince(db, findSince)

	for _, consumption := range consumptions {
		db.Model(&consumption).Related(&consumption.Consumable)
		mgRemaining := amountRemainingAtTime(consumption, time)
		totalMg += mgRemaining
	}
	return totalMg
}

func mgOverTime(db gorm.DB, startTime time.Time, endTime time.Time, increment time.Duration) []MgAtTime {
	var points []MgAtTime

	currentTime := startTime
	for currentTime.Before(endTime) {
		points = append(points, MgAtTime{
			Time:   currentTime,
			Amount: mgInBodyAtTime(db, currentTime)})
		currentTime = currentTime.Add(increment)
	}

	return points
}
