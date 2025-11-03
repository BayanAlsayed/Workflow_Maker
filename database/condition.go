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
			DESCRIPTION,
			
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
		conditions = append(conditions, cond)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return conditions, nil
}