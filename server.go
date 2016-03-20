package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Initializes and runs the web server
func initServer() {
	router := mux.NewRouter()
	router.HandleFunc("/consumables", consumablesListHandler).Methods("GET")
	router.HandleFunc("/consumables", optionsHandler).Methods("OPTIONS")
	router.HandleFunc("/consumables", consumablesCreateHandler).Methods("POST")
	router.HandleFunc("/ingest", optionsHandler).Methods("OPTIONS")
	router.HandleFunc("/ingest", ingestHandler).Methods("POST")
	router.HandleFunc("/status/now", statusHandler).Methods("GET")
	router.HandleFunc("/status/time", statusTimeHandler).Methods("GET")
	http.Handle("/", router)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("CAFFEINE_PORT")), nil)
}

// Create a Consumable
func consumablesCreateHandler(w http.ResponseWriter, r *http.Request) {
	type consumableJSON struct {
		Name   string `json:"name"`
		Amount string `json:"amount"`
	}

	apiApplyCorsHeaders(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var consumableForm consumableJSON
	err := json.Unmarshal(body, &consumableForm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	amount, err := strconv.Atoi(consumableForm.Amount)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	consumable := Consumable{Name: consumableForm.Name, Amount: uint(amount)}
	db.Create(&consumable)
	w.WriteHeader(http.StatusNoContent)
}

// Index of Consumables
func consumablesListHandler(w http.ResponseWriter, r *http.Request) {
	type consumableResponse struct {
		ID     uint   `json:"id"`
		Name   string `json:"name"`
		Amount uint   `json:"amount"`
	}

	// Find all Consumables
	var consumables []Consumable
	db.Find(&consumables)

	// Build as consumableResponse
	var response []consumableResponse
	for _, consumable := range consumables {
		item := consumableResponse{
			ID:     consumable.ID,
			Name:   consumable.Name,
			Amount: consumable.Amount}
		response = append(response, item)
	}

	// Write JSON to response
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	apiApplyCorsHeaders(w, r)
	w.Write(jsonResponse)
}

func ingestHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ConsumableID uint `json:"consumable_id"`
		EnergyLevel  uint `json:"energy_level"`
	}

	// Read body, parse JSON
	body, _ := ioutil.ReadAll(r.Body)
	var requestBody request
	err := json.Unmarshal(body, &requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Find and ingest Consumable
	var consumable Consumable
	var consumption Consumption

	db.Where("id = ?", requestBody.ConsumableID).Find(&consumable)
	if consumable.ID > 0 {
		consumption = ingest(db, consumable)
	} else {
		// Not found
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Record EnergyLevel
	db.Create(&EnergyLevel{
		Level:       requestBody.EnergyLevel,
		Consumption: consumption})

	apiApplyCorsHeaders(w, r)
	w.WriteHeader(http.StatusNoContent)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		MgInBody float64 `json:"mg_in_body"`
	}

	jsonResponse, _ := json.Marshal(response{MgInBody: mgInBody(db)})
	w.Header().Set("Content-Type", "application/json")
	apiApplyCorsHeaders(w, r)
	w.Write(jsonResponse)
}

func statusTimeHandler(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	startTime, _ := time.Parse(time.RFC3339, r.Form["start_time"][0])
	endTime, _ := time.Parse(time.RFC3339, r.Form["end_time"][0])
	intervalHours, _ := strconv.Atoi(r.Form["interval"][0])

	increment, _ := time.ParseDuration(fmt.Sprintf("%dh", intervalHours))
	points := mgOverTime(db, startTime, endTime, increment)

	jsonResponse, _ := json.Marshal(points)
	w.Header().Set("Content-Type", "application/json")
	apiApplyCorsHeaders(w, r)
	w.Write(jsonResponse)
}
