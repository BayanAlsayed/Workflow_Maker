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
