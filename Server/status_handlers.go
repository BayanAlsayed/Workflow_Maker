package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
