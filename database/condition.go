package database

type Condition struct {
	WF_CONDITION_ID int    `json:"wf_condition_id"`
	FUNC_NAME       string `json:"func_name"`
	DESCRIPTION     string `json:"description"`
	ACTIVE_WORKFLOWS_USING  string `json:"active_workflows_using"`
	ACTIVE_RULES_COUNT       int    `json:"active_rules_count"`
}
func GetAllConditions() ([]Condition, error) {
	rows, err := db.Query(`
		SELECT 
			WF_CONDITION_ID,
			FUNC_NAME,
			DESCRIPTION
		FROM WF_CONDITION
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conditions []Condition
	for rows.Next() {
		var cond Condition
		if err := rows.Scan(&cond.WF_CONDITION_ID, &cond.FUNC_NAME, &cond.DESCRIPTION); err != nil {
			return nil, err
		}
		db.QueryRow(`
			SELECT GROUP_CONCAT(DISTINCT W.WORKFLOW_NAME, ' (v', V.VERSION, ')' SEPARATOR ', ') AS ACTIVE_WORKFLOWS_USING,
			COUNT(DISTINCT CR.RULE_ID) AS ACTIVE_RULES_COUNT
			FROM WF_CONDITION_RULE CR
			JOIN WF_VERSION V ON CR.WF_VERSION_ID = V.WF_VERSION_ID
			JOIN WF_WORKFLOW W ON V.WORKFLOW_ID = W.WORKFLOW_ID
			WHERE CR.WF_CONDITION_ID = ? AND V.IS_ACTIVE = 1
		`, cond.WF_CONDITION_ID).Scan(&cond.ACTIVE_WORKFLOWS_USING, &cond.ACTIVE_RULES_COUNT)
		conditions = append(conditions, cond)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return conditions, nil
}

func AddCondition(funcName string, description string) error {
	_, err := db.Exec(`
		INSERT INTO WF_CONDITION (FUNC_NAME, DESCRIPTION)
		VALUES (?, ?)
	`, funcName, description)
	return err
}