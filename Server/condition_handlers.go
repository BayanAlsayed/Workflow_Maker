package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	database "workflow/database"
)

func ViewConditionsHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch all conditions from the database
	conditions, err := database.GetAllConditions()
	if err != nil {
		fmt.Println("Error fetching conditions: ", err)
		http.Error(w, "Error fetching conditions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conditions)
}

func AddConditionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AddConditionHandler called")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		RenderTemplate(w, "405.html", nil)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fmt.Println("condition to be added: " + string(bodyBytes))
	//condition to be added: {"function_name":"bayan","description":"bayan"}

	type ConditionInput struct {
		FUNC_NAME   string `json:"function_name"`
		DESCRIPTION string `json:"description"`
	}

	var input ConditionInput
	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.FUNC_NAME == "" || input.DESCRIPTION == "" {
		http.Error(w, "FUNCTION_NAME and DESCRIPTION are required", http.StatusBadRequest)
		return
	}

	err = database.AddCondition(input.FUNC_NAME, input.DESCRIPTION)
	if err != nil {
		fmt.Println("Error adding condition:", err)
		http.Error(w, "Failed to add condition", http.StatusInternalServerError)
		return
	}

	fmt.Println("Condition added successfully")

}
