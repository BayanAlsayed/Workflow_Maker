package database

import "fmt"

type Rule_Condition struct {
	WF_RULE_CONDITION_ID     int    `json:"wf_rule_condition_id"`
	RULE_CONDITION_ID        int    `json:"rule_condition_id"`
	WF_CONDITION_ID          int    `json:"wf_condition_id"`
	WF_CONDITION_FUNC_NAME   string `json:"wf_condition_func_name"`
	WF_CONDITION_DESCRIPTION string `json:"wf_condition_description"`
	RULE_CONDITION_TYPE      string `json:"rule_condition_type"`
}

func AddRuleCondition(WF_CONDITION_ID int, CONDITION_TYPE string, RULE_ID int, WORKFLOW_ID int, VERSION int) error {
	//get the WF_VERSION_ID for the given WORKFLOW_ID and VERSION
	WF_VERSION_ID, err := getVersionID(WORKFLOW_ID, VERSION)
	if err != nil {
		return err
	}

	ruleConditionID, err := getNextRuleConditionID(WF_VERSION_ID)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO WF_RULE_CONDITION (
			RULE_CONDITION_ID,
			WF_VERSION_ID,
			RULE_ID,
			WF_CONDITION_ID,
			TYPE,
			EXPECTED_VALUE
		) VALUES (?, ?, ?, ?, ?, ?)
	`, ruleConditionID, WF_VERSION_ID, RULE_ID, WF_CONDITION_ID, CONDITION_TYPE, true)
	if err != nil {
		return err
	}
	return nil
}

func getNextRuleConditionID(WF_VERSION_ID int) (int, error) {
	var nextID int
	err := db.QueryRow(`
		SELECT IFNULL(MAX(RULE_CONDITION_ID), 0) + 1 AS next_id
		FROM WF_RULE_CONDITION
		WHERE WF_VERSION_ID = ?
	`, WF_VERSION_ID).Scan(&nextID)
	if err != nil {
		return 0, err
	}
	return nextID, nil
}

func DeleteRuleCondition(WF_RULE_CONDITION_ID int) error {
	_, err := db.Exec(`
		DELETE FROM WF_RULE_CONDITION
		WHERE WF_RULE_CONDITION_ID = ?
	`, WF_RULE_CONDITION_ID)
	if err != nil {
		fmt.Println("Error deleting rule condition: ", err)
		return err
	}
	return nil
}

func UpdateRuleCondition(WF_RULE_CONDITION_ID int, WF_CONDITION_ID int, CONDITION_TYPE string) error {
	_, err := db.Exec(`
		UPDATE WF_RULE_CONDITION
		SET 
			WF_CONDITION_ID = ?,
			TYPE = ?
		WHERE WF_RULE_CONDITION_ID = ?
	`, WF_CONDITION_ID, CONDITION_TYPE, WF_RULE_CONDITION_ID)
	if err != nil {
		fmt.Println("Error updating rule condition: ", err)
		return err
	}

	return nil
}
