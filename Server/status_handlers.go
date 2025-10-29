package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	database "workflow/database"
)

func AddStatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AddStatusHandler called")
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

	fmt.Println("status to be added: " + string(bodyBytes))

	type StatusInput struct {
		STATUS_NAME           string `json:"status_name"`
		ED_CODE_STATUS_CAT_ID string `json:"ed_code_status_cat_id"`
		ED_CODE_STATUS_ID     string `json:"ed_code_status_id"`
		GS_CODE_REQ_STATUS_ID string `json:"gs_code_req_status_id"`
		IS_TERMINAL           string `json:"is_terminal"`
		SUCCESS_PATH          string `json:"success_path"`
		WF_ID                 int    `json:"workflow_id"`
	}

	var input StatusInput
	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.STATUS_NAME == "" {
		http.Error(w, "STATUS_NAME is required", http.StatusBadRequest)
		return
	} else if input.GS_CODE_REQ_STATUS_ID == "" && (input.ED_CODE_STATUS_ID == "" && input.ED_CODE_STATUS_CAT_ID == "") {

		http.Error(w, "At least one of GS_CODE_REQ_STATUS_ID, ED_CODE_STATUS_ID, or ED_CODE_STATUS_CAT_ID is required", http.StatusBadRequest)
		return
	}

	var ED_CODE_STATUS_CAT = strings.Split(input.ED_CODE_STATUS_CAT_ID, "_")
	var ED_CODE_STATUS = strings.Split(input.ED_CODE_STATUS_ID, "_")
	var GS_CODE_REQ_STATUS = strings.Split(input.GS_CODE_REQ_STATUS_ID, "_")

	// {"status_name":"new","ed_code_status_cat_id":"","ed_code_status_id":"","gs_code_req_status_id":"1_New","is_terminal":"0","success_path":"","workflow_id":5}

	err = database.AddStatus(input.WF_ID, input.STATUS_NAME, input.IS_TERMINAL, input.SUCCESS_PATH, ED_CODE_STATUS_CAT, ED_CODE_STATUS, GS_CODE_REQ_STATUS)
	if err != nil {
		fmt.Println("Error adding status:", err)
		http.Error(w, "Failed to add status", http.StatusInternalServerError)
		return
	}
}

func EditStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET (as you designed)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query params
	q := r.URL.Query()

	workflowID, _ := strconv.Atoi(q.Get("workflow_id"))
	statusID, _ := strconv.Atoi(q.Get("status_id"))
	statusName := q.Get("status_name")

	edCodeStr := q.Get("ed_code_status_id")
	gsCodeStr := q.Get("gs_code_req_status_id")
	isTerminalStr := q.Get("is_terminal")
	successPathStr := q.Get("success_path")

	var (
		edCodeID, gsCodeID, successPath *int
		isTerminal                      bool
	)

	// Parse optional integers safely
	if v, err := strconv.Atoi(edCodeStr); err == nil {
		edCodeID = &v
	}
	if v, err := strconv.Atoi(gsCodeStr); err == nil {
		gsCodeID = &v
	}
	if v, err := strconv.Atoi(successPathStr); err == nil {
		successPath = &v
	}

	// Parse boolean from "0"/"1"
	isTerminal = isTerminalStr == "1"

	fmt.Printf("Received edit status:\n"+
		"WF_ID=%d | STATUS_ID=%d | STATUS_NAME=%s | ED=%v | GS=%v | TERMINAL=%v | PATH=%v\n",
		workflowID, statusID, statusName, edCodeID, gsCodeID, isTerminal, successPath)

	// TODO: Update in DB
	err := database.UpdateStatus(database.Status{
		STATUS_ID:          statusID,
		STATUS_NAME:       statusName,
		ED_CODE_STATUS_ID: edCodeID,
		GS_CODE_REQ_STATUS_ID: gsCodeID,
		IS_TERMINAL:       isTerminal,
		SUCCESS_PATH:      successPath,
	}, workflowID)

	if err != nil {
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	// Respond JSON success
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))

}

func DeleteStatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteStatusHandler called")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	workflowID, _ := strconv.Atoi(q.Get("workflow_id"))
	statusID, _ := strconv.Atoi(q.Get("status_id"))

	err := database.DeleteStatus(statusID, workflowID)
	if err != nil {
		fmt.Println("Error deleting status:", err)
		http.Error(w, "Failed to delete status", http.StatusInternalServerError)
		return
	}

	// Respond JSON success
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
}
