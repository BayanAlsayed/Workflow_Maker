package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	database "workflow/database"
)

func AddRuleConditionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AddRuleConditionHandler called")

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

	fmt.Println("rule condition to be added: " + string(bodyBytes))
	//rule condition to be added: {"wf_condition_id":"1","condition_type":"pre","rule_id":1,"workflow_id":2,"version":1}

	type RuleConditionInput struct {
		WF_CONDITION_ID string `json:"wf_condition_id"`
		CONDITION_TYPE  string `json:"condition_type"`
		RULE_ID         int    `json:"rule_id"`
		WORKFLOW_ID     int    `json:"workflow_id"`
		VERSION         int    `json:"version"`
	}

	var input RuleConditionInput
	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	WF_CONDITION_ID, err := strconv.Atoi(input.WF_CONDITION_ID)
	if err != nil {
		fmt.Println("Error converting WF_CONDITION_ID:", err)
		http.Error(w, "Invalid WF_CONDITION_ID", http.StatusBadRequest)
		return
	}

	if WF_CONDITION_ID == 0 || input.CONDITION_TYPE == "" || input.RULE_ID == 0 {
		http.Error(w, "WF_CONDITION_ID, CONDITION_TYPE, and RULE_ID are required", http.StatusBadRequest)
		return
	}

	err = database.AddRuleCondition(WF_CONDITION_ID, input.CONDITION_TYPE, input.RULE_ID, input.WORKFLOW_ID, input.VERSION)
	if err != nil {
		fmt.Println("Error adding rule condition:", err)
		http.Error(w, "Failed to add rule condition", http.StatusInternalServerError)
		return
	}

}

func DeleteRuleConditionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteRuleConditionHandler called")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	wfRuleConditionID, _ := strconv.Atoi(q.Get("id"))

	if err := database.DeleteRuleCondition(wfRuleConditionID); err != nil {
		http.Error(w, "Failed to delete rule condition", http.StatusInternalServerError)
		return
	}

}

func EditRuleConditionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("EditRuleConditionHandler called")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	wf_rule_condition_id, _ := strconv.Atoi(q.Get("wf_rule_condition_id"))
	wf_condition_id, _ := strconv.Atoi(q.Get("wf_condition_id"))
	condition_type := q.Get("condition_type")

	err := database.UpdateRuleCondition(wf_rule_condition_id, wf_condition_id, condition_type)
	if err != nil {
		fmt.Println("Error updating rule condition:", err)
		http.Error(w, "Failed to update rule condition", http.StatusInternalServerError)
		return
	}

	// Respond JSON success
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))

}
