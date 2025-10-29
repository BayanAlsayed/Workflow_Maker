package database

import (
	"database/sql"
	"fmt"
)

type Rule struct {
	WF_RULE_ID           int     `json:"wf_rule_id"`
	RULE_ID              int     `json:"rule_id"`
	FROM_STATUS_ID       int     `json:"from_status_id"`
	FROM_STATUS_NAME     string  `json:"from_status_name"`
	TO_STATUS_ID         int     `json:"to_status_id"`
	TO_STATUS_NAME       string  `json:"to_status_name"`
	SE_CODE_USER_TYPE_ID int     `json:"se_code_user_type_id"`
	USER_TYPE_EN         *string `json:"user_type_en"`
	SE_ACCNT_ID          *int    `json:"se_accnt_id"`
	ACCNT_EN             *string `json:"accnt_en"`
	ACTION_BUTTON        string  `json:"action_button"`
	ACTION_FUNCTION      *string `json:"action_function"`
}

func ViewRules(WF_ID int) ([]Rule, error) {
	//select the newest version for the given workflow
	newestVersion, err := getNewestVersion(WF_ID)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
			SELECT 
				R.WF_RULE_ID,
				R.RULE_ID,
				R.FROM_STATUS_ID,
				S.STATUS_NAME AS from_status_name,
				R.TO_STATUS_ID,
				T.STATUS_NAME AS to_status_name,
				R.SE_CODE_USER_TYPE_ID,
				UT.DESCR_EN    AS user_type_name,
				R.SE_ACCNT_ID,
				ACC.DESCR_EN   AS account_name,
				R.ACTION_BUTTON,
				R.ACTION_FUNCTION
			FROM WF_RULE AS R
			JOIN WF_STATUS AS S ON R.FROM_STATUS_ID = S.STATUS_ID AND R.WF_VERSION_ID = S.WF_VERSION_ID
			JOIN WF_STATUS AS T ON R.TO_STATUS_ID   = T.STATUS_ID AND R.WF_VERSION_ID = T.WF_VERSION_ID
			LEFT JOIN SE_CODE_USER_TYPE AS UT ON R.SE_CODE_USER_TYPE_ID = UT.SE_CODE_USER_TYPE_ID
			LEFT JOIN SE_ACCNT AS ACC ON R.SE_ACCNT_ID = ACC.SE_ACCNT_ID
			WHERE R.WF_VERSION_ID = ?
	`, newestVersion)

	if err != nil {
		fmt.Println("Error fetching rules: ", err)

		if err == sql.ErrNoRows {
			return []Rule{}, nil // return empty slice if no rules found
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	var rules []Rule
	for rows.Next() {
		var r Rule

		var (
			userTypeEn     sql.NullString
			seAccntID      sql.NullInt64
			accntEn        sql.NullString
			actionFunction sql.NullString
		)

		if err := rows.Scan(
			&r.WF_RULE_ID,
			&r.RULE_ID,
			&r.FROM_STATUS_ID,
			&r.FROM_STATUS_NAME,
			&r.TO_STATUS_ID,
			&r.TO_STATUS_NAME,
			&r.SE_CODE_USER_TYPE_ID,
			&r.USER_TYPE_EN,
			&seAccntID,
			&accntEn,
			&r.ACTION_BUTTON,
			&actionFunction,
		); err != nil {
			return nil, err
		}

		// Handle nullable fields
		if userTypeEn.Valid {
			v := userTypeEn.String
			r.USER_TYPE_EN = &v
		}
		if seAccntID.Valid {
			v := int(seAccntID.Int64)
			r.SE_ACCNT_ID = &v
		}
		if accntEn.Valid {
			v := accntEn.String
			r.ACCNT_EN = &v
		}
		if actionFunction.Valid {
			v := actionFunction.String
			r.ACTION_FUNCTION = &v
		}

		rules = append(rules, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

func AddRule(WF_ID int, FROM_STATUS_SLICE []string, TO_STATUS_SLICE []string, USER_TYPE_SLICE []string, ACCOUNT_SLICE []string, ACTION_BUTTON string, ACTION_FUNCTION string) error {
	//select the newest version for the given workflow
	newestVersion, err := getNewestVersion(WF_ID)
	if err != nil {
		fmt.Println("Error getting newest version: ", err)
		return err
	}

	FROM_STATUS_ID, err := getIDFromStringSlice(FROM_STATUS_SLICE)
	if err != nil {
		fmt.Println("Error parsing FROM_WF_STATUS_SLICE : ", err)
		return err
	}
	TO_STATUS_ID, err := getIDFromStringSlice(TO_STATUS_SLICE)
	if err != nil {
		fmt.Println("Error parsing TO_WF_STATUS_SLICE : ", err)
		return err
	}
	USER_TYPE_ID, err := getIDFromStringSlice(USER_TYPE_SLICE)
	if err != nil {
		fmt.Println("Error parsing USER_TYPE_SLICE : ", err)
		return err
	}
	ACCOUNT_ID, err := getIDFromStringSlice(ACCOUNT_SLICE)
	if err != nil {
		fmt.Println("Error parsing ACCOUNT_SLICE : ", err)
		return err
	}

	ruleID, err := getNextRuleID(newestVersion)
	if err != nil {
		fmt.Println("Error getting next status ID: ", err)
		return err
	}

	_, err = db.Exec(`
		INSERT INTO WF_RULE (
			RULE_ID,
			WF_VERSION_ID,
			FROM_STATUS_ID,
			TO_STATUS_ID,
			SE_CODE_USER_TYPE_ID,
			SE_ACCNT_ID,
			ACTION_BUTTON,
			ACTION_FUNCTION
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, ruleID, newestVersion, FROM_STATUS_ID, TO_STATUS_ID, USER_TYPE_ID, ACCOUNT_ID, ACTION_BUTTON, ACTION_FUNCTION)
	if err != nil {
		fmt.Println("Error inserting new rule: ", err)
		return err
	}
	return nil
}
func strPtr(s string) *string {
	return &s
}

func UpdateRule(rule Rule, WorkflowID int) error {
	//select the newest version for the given workflow
	newestVersion, err := getNewestVersion(WorkflowID)
	if err != nil {
		fmt.Println("Error getting newest version: ", err)
		return err
	}

	fmt.Println("updated rule:", rule)

	if rule.ACTION_FUNCTION == nil {
		rule.ACTION_FUNCTION = strPtr("")
	}
	_, err = db.Exec(`
		UPDATE WF_RULE
		SET
			FROM_STATUS_ID = ?,
			TO_STATUS_ID = ?,
			SE_CODE_USER_TYPE_ID = ?,
			SE_ACCNT_ID = ?,
			ACTION_BUTTON = ?,
			ACTION_FUNCTION = ?
		WHERE WF_VERSION_ID = ? AND RULE_ID = ?
	`, rule.FROM_STATUS_ID, rule.TO_STATUS_ID, rule.SE_CODE_USER_TYPE_ID, rule.SE_ACCNT_ID, rule.ACTION_BUTTON, rule.ACTION_FUNCTION, newestVersion, rule.RULE_ID)
	if err != nil {
		fmt.Println("Error updating rule: ", err)
		return err
	}
	return nil
}

func DeleteRule(ruleID int, workflowID int) error {
	newestVersion, err := getNewestVersion(workflowID)
	if err != nil {
		fmt.Println("Error getting newest version: ", err)
		return err
	}
	_, err = db.Exec(`
		DELETE FROM WF_RULE
		WHERE RULE_ID = ? AND WF_VERSION_ID = ?
	`, ruleID, newestVersion)
	if err != nil {
		fmt.Println("Error deleting rule: ", err)
		return err
	}
	return nil
}

func getNextRuleID(WF_VERSION_ID int) (int, error) {
	var nextID int
	err := db.QueryRow(`
		SELECT IFNULL(MAX(RULE_ID), 0) + 1 AS next_id
		FROM WF_RULE
		WHERE WF_VERSION_ID = ?
	`, WF_VERSION_ID).Scan(&nextID)
	if err != nil {
		return 0, err
	}
	return nextID, nil
}
