package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	database "workflow/database"
)

func AddWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		RenderTemplate(w, "405.html", nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Println("Error parsing form: ", err)
		w.WriteHeader(http.StatusBadRequest)
		RenderTemplate(w, "400.html", nil)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		RenderTemplate(w, "400.html", "Workflow name is required")
		return
	}

	newWorkflow := database.Workflow{
		WORKFLOW_NAME:        name,
		WORKFLOW_DESCRIPTION: description,
	}

	if _, err := database.AddWorkflow(newWorkflow); err != nil {
		fmt.Println("Error adding workflow: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		RenderTemplate(w, "500.html", nil)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteWorkflowHandler called")
	// if r.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// 	RenderTemplate(w, "405.html", nil)
	// 	return
	// }
	fmt.Println("Request method:", r.Method)

	if r.Method == "GET" {
		id := r.URL.Query().Get("Workflow_ID")

		fmt.Println("Received request to delete workflow with ID:", id)
		if err := database.DeleteWorkflowByID(id); err != nil {
			fmt.Println("Error deleting workflow: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			RenderTemplate(w, "500.html", nil)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func EditWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		RenderTemplate(w, "405.html", nil)
		return
	}

	var wf database.Workflow // use your existing struct

	// Decode the JSON body into wf
	if err := json.NewDecoder(r.Body).Decode(&wf); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Println("Editing workflow ID:", wf.WORKFLOW_ID, "with name:", wf.WORKFLOW_NAME, "and description:", wf.WORKFLOW_DESCRIPTION)

	if strconv.Itoa(wf.WORKFLOW_ID) == "" || wf.WORKFLOW_NAME == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Workflow ID or name is missing")
		RenderTemplate(w, "400.html", "Workflow ID and name are required")
		return
	}

	updatedWorkflow := database.Workflow{
		WORKFLOW_ID:          wf.WORKFLOW_ID,
		WORKFLOW_NAME:        wf.WORKFLOW_NAME,
		WORKFLOW_DESCRIPTION: wf.WORKFLOW_DESCRIPTION,
	}

	if err := database.UpdateWorkflow(updatedWorkflow); err != nil {
		fmt.Println("Error updating workflow: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		RenderTemplate(w, "500.html", nil)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
