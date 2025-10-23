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

func ViewDetailsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ViewDetailsHandler called")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		RenderTemplate(w, "405.html", nil)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/view_workflow/")
	WF_ID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
		return
	}

	statuses, err := database.ViewStatuses(WF_ID)
	if err != nil {
		http.Error(w, "Failed to retrieve statuses", http.StatusInternalServerError)
		return
	}

	rules, err := database.ViewRules(WF_ID)
	if err != nil {
		http.Error(w, "Failed to retrieve rules", http.StatusInternalServerError)
		return
	}

	details := struct {
		STATUSES []database.Status `json:"statuses"`
		RULES    []database.Rule   `json:"rules"`
	}{
		STATUSES: statuses,
		RULES:    rules,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}

func AddRuleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AddRuleHandler called")
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

	fmt.Println("rule to be added: " + string(bodyBytes))

	type RuleInput struct {
		FROM_STATUS string `json:"from_status"`
		TO_STATUS   string `json:"to_status"`
		SE_CODE_USER_TYPE      string `json:"se_code_user_type"`
		SE_ACCNT        string `json:"se_accnt"`
		ACTION_BUTTON     string `json:"action_button"`
		ACTION_FUNCTION   string `json:"action_function"`
		WF_ID       int    `json:"workflow_id"`
	}

	var input RuleInput
	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.SE_ACCNT == "" && input.SE_CODE_USER_TYPE == "" {
		fmt.Println("atleast one of SE_ACCNT or SE_CODE_USER_TYPE is required")
		http.Error(w, "atleast one of SE_ACCNT or SE_CODE_USER_TYPE is required", http.StatusBadRequest)
		return
	}

	var FROM_STATUS_SLICE = strings.Split(input.FROM_STATUS, "_")
	var TO_STATUS_SLICE = strings.Split(input.TO_STATUS, "_")
	var SE_CODE_USER_TYPE_SLICE = strings.Split(input.SE_CODE_USER_TYPE, "_")
	var SE_ACCNT_SLICE = strings.Split(input.SE_ACCNT, "_")

	err = database.AddRule(input.WF_ID, FROM_STATUS_SLICE, TO_STATUS_SLICE, SE_CODE_USER_TYPE_SLICE, SE_ACCNT_SLICE, input.ACTION_BUTTON, input.ACTION_FUNCTION)
	if err != nil {
		fmt.Println("Error adding rule:", err)
		http.Error(w, "Failed to add rule", http.StatusInternalServerError)
		return
	}

}

func LookupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		RenderTemplate(w, "405.html", nil)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/lookups/")
	WF_ID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
		return
	}

	userTypes, err := database.GetUserTypes()
	if err != nil {
		http.Error(w, "Failed to retrieve user types", http.StatusInternalServerError)
		return
	}

	accounts, err := database.GetAccounts()
	if err != nil {
		http.Error(w, "Failed to retrieve accounts", http.StatusInternalServerError)
		return
	}

	edStatusCodesCat, err := database.GetEdCodeStatusesCat()
	if err != nil {
		http.Error(w, "Failed to retrieve ED status codes", http.StatusInternalServerError)
		return
	}

	edStatusCodes, err := database.GetEdCodeStatuses()
	if err != nil {
		http.Error(w, "Failed to retrieve ED status codes", http.StatusInternalServerError)
		return
	}

	gsStatusCodes, err := database.GetGsCodeReqStatuses()
	if err != nil {
		http.Error(w, "Failed to retrieve GS status codes", http.StatusInternalServerError)
		return
	}

	workflowStatuses, err := database.ViewStatuses(WF_ID)
	if err != nil {
		http.Error(w, "Failed to retrieve workflow statuses", http.StatusInternalServerError)
		return
	}

	lookups := struct {
		UserTypes        []database.UserType        `json:"user_types"`
		Accounts         []database.Account         `json:"accounts"`
		EDStatusCodesCat []database.EDCodeStatusCat `json:"ed_status_codes_cat"`
		EDStatusCodes    []database.EDCodeStatus    `json:"ed_status_codes"`
		GSStatusCodes    []database.GSCodeReqStatus `json:"gs_status_codes"`
		WorkflowStatuses []database.Status          `json:"workflow_statuses"`
	}{
		UserTypes:        userTypes,
		Accounts:         accounts,
		EDStatusCodesCat: edStatusCodesCat,
		EDStatusCodes:    edStatusCodes,
		GSStatusCodes:    gsStatusCodes,
		WorkflowStatuses: workflowStatuses,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lookups)
}
