package server

import (
	"encoding/json"
	"net/http"

	database "workflow/database"
)

func ViewConditionsHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch all conditions from the database
	conditions, err := database.GetAllConditions()
	if err != nil {
		http.Error(w, "Error fetching conditions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conditions)
}
