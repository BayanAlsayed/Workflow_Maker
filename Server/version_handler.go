package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	database "workflow/database"
)

func GetWorkflowVersionsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetWorkflowVersionsHandler called")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		RenderTemplate(w, "405.html", nil)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/get_workflow_versions/")
	WF_ID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
		return
	}

	versions, err := database.GetWorkflowVersions(WF_ID)
	if err != nil {
		http.Error(w, "Failed to fetch workflow versions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versions)

}

func ActivateWorkflowVersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ActivateWorkflowVersionHandler called")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		RenderTemplate(w, "405.html", nil)
		return
	}

	var input struct {
		WORKFLOW_ID int `json:"workflow_id"`
		VERSION int `json:"version"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := database.ActivateWorkflowVersion(input.WORKFLOW_ID, input.VERSION)
	if err != nil {
		http.Error(w, "Failed to activate workflow version", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}


func ApproveWorkflowVersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ApproveWorkflowVersionHandler called")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		RenderTemplate(w, "405.html", nil)
		return
	}

	var input struct {
		WORKFLOW_ID int `json:"workflow_id"`
		VERSION int `json:"version"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := database.ApproveWorkflowVersion(input.WORKFLOW_ID, input.VERSION)
	if err != nil {
		http.Error(w, "Failed to approve workflow version", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func DeleteWorkflowVersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteWorkflowHandler called")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		RenderTemplate(w, "405.html", nil)
		return
	}


	var input struct {
		WORKFLOW_ID int `json:"workflow_id"`
		VERSION int `json:"version"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := database.DeleteWorkflowVersion(input.WORKFLOW_ID, input.VERSION); err != nil {
		fmt.Println("Error deleting workflow version: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		RenderTemplate(w, "500.html", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func CreateWorkflowVersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateWorkflowVersionHandler called")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		RenderTemplate(w, "405.html", nil)
		return
	}

	var WORKFLOW_ID int

	if err := json.NewDecoder(r.Body).Decode(&WORKFLOW_ID); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newVersion, err := database.AddVersion(WORKFLOW_ID)
	if err != nil {
		http.Error(w, "Failed to create workflow version", http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newVersion)
}

func DuplicateWorkflowVersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DuplicateWorkflowVersionHandler")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		RenderTemplate(w, "405.html", nil)
		return
	}


	var input struct {
		WORKFLOW_ID int `json:"workflow_id"`
		VERSION int `json:"version"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := database.DuplicateWorkflowVersion(input.WORKFLOW_ID, input.VERSION); err != nil {
		fmt.Println("Error duplicating workflow version: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		RenderTemplate(w, "500.html", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
