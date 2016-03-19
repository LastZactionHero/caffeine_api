package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// Initializes and runs the web server
func initServer() {
	router := mux.NewRouter()
	router.HandleFunc("/consumables", consumablesListHandler).Methods("GET")
	router.HandleFunc("/ingest", ingestHandler).Methods("POST")
	router.HandleFunc("/status/now", statusHandler).Methods("GET")
	http.Handle("/", router)
	http.ListenAndServe(":8181", nil)
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
	w.Write(jsonResponse)
}

func ingestHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ConsumableID uint `json:"consumable_id"`
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
	db.Where("id = ?", requestBody.ConsumableID).Find(&consumable)
	if consumable.ID > 0 {
		ingest(db, consumable)
	} else {
		// Not found
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		MgInBody float64 `json:"mg_in_body"`
	}

	jsonResponse, _ := json.Marshal(response{MgInBody: mgInBody(db)})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
