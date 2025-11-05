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

	Str := strings.Split(strings.TrimPrefix(r.URL.Path, "/view_workflow/"), "/")
	WF_ID, err := strconv.Atoi(Str[0])
	if err != nil {
		http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
		return
	}
	version, err := strconv.Atoi(Str[1])
	if err != nil {
		http.Error(w, "Invalid workflow version", http.StatusBadRequest)
		return
	}

	statuses, err := database.ViewStatuses(WF_ID, version)
	if err != nil {
		http.Error(w, "Failed to retrieve statuses", http.StatusInternalServerError)
		return
	}

	rules, err := database.ViewRules(WF_ID, version)
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
		VERSION int `json:"version"`
	}

	var input RuleInput
	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	fmt.Println("rule input: ", input)

	if input.SE_ACCNT == "" && input.SE_CODE_USER_TYPE == "" {
		fmt.Println("atleast one of SE_ACCNT or SE_CODE_USER_TYPE is required")
		http.Error(w, "atleast one of SE_ACCNT or SE_CODE_USER_TYPE is required", http.StatusBadRequest)
		return
	}

	var FROM_STATUS_SLICE = strings.Split(input.FROM_STATUS, "_")
	var TO_STATUS_SLICE = strings.Split(input.TO_STATUS, "_")
	var SE_CODE_USER_TYPE_SLICE = strings.Split(input.SE_CODE_USER_TYPE, "_")
	var SE_ACCNT_SLICE = strings.Split(input.SE_ACCNT, "_")

	err = database.AddRule(input.WF_ID, input.VERSION, FROM_STATUS_SLICE, TO_STATUS_SLICE, SE_CODE_USER_TYPE_SLICE, SE_ACCNT_SLICE, input.ACTION_BUTTON, input.ACTION_FUNCTION)
	if err != nil {
		fmt.Println("Error adding rule:", err)
		http.Error(w, "Failed to add rule", http.StatusInternalServerError)
		return
	}

}

func EditRuleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("EditRuleHandler called")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query params
	q := r.URL.Query()

	workflowID, _ := strconv.Atoi(q.Get("workflow_id"))
	version, _:= strconv.Atoi(q.Get("version"))
	ruleID, _ := strconv.Atoi(q.Get("rule_id"))
	fromStatusId, _ := strconv.Atoi(q.Get("from_status_id"))
	toStatusId, _ := strconv.Atoi(q.Get("to_status_id"))
	seCodeUserType, _ := strconv.Atoi(q.Get("se_code_user_type_id"))
	seAccntSTR := q.Get("se_accnt_id")
	actionButton := q.Get("action_button")
	var actionFunction *string
	if v := q.Get("action_function"); v != "" {
		actionFunction = &v
	}

	var seAccntID *int
	if v, err := strconv.Atoi(seAccntSTR); err == nil {
		seAccntID = &v
	}

	err := database.UpdateRule(database.Rule{
		RULE_ID:         ruleID,
		FROM_STATUS_ID:   fromStatusId,
		TO_STATUS_ID:     toStatusId,
		SE_CODE_USER_TYPE_ID: seCodeUserType,
		SE_ACCNT_ID:      seAccntID,
		ACTION_BUTTON:   actionButton,
		ACTION_FUNCTION: actionFunction,
	}, workflowID, version)

	if err != nil {
		fmt.Println("Error updating rule:", err)
		http.Error(w, "Failed to update rule", http.StatusInternalServerError)
		return
	}

	// Respond JSON success
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
}

func DeleteRuleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteRuleHandler called")
	
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	workflowID, _ := strconv.Atoi(q.Get("workflow_id"))
	version, _:= strconv.Atoi(q.Get("version"))
	ruleID, _ := strconv.Atoi(q.Get("rule_id"))

	err := database.DeleteRule(ruleID, workflowID, version)
	if err != nil {
		fmt.Println("Error deleting rule:", err)
		http.Error(w, "Failed to delete rule", http.StatusInternalServerError)
		return
	}

	// Respond JSON success
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
}



