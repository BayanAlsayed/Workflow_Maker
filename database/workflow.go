package database

type Workflow struct {
	WORKFLOW_ID          int    `json:"workflow_id"`
	WORKFLOW_NAME        string `json:"workflow_name"`
	WORKFLOW_DESCRIPTION string `json:"workflow_description"`
}

func GetAllWorkflows() ([]Workflow, error) {
	rows, err := db.Query(`SELECT WORKFLOW_ID, WORKFLOW_NAME, DESCRIPTION FROM WORKFLOW`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []Workflow
	for rows.Next() {
		var wf Workflow
		if err := rows.Scan(&wf.WORKFLOW_ID, &wf.WORKFLOW_NAME, &wf.WORKFLOW_DESCRIPTION); err != nil {
			return nil, err
		}
		workflows = append(workflows, wf)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return workflows, nil
}

// AddWorkflow inserts a workflow and creates its initial version (v1) atomically.
// If anything fails, nothing is persisted.
func AddWorkflow(wf Workflow) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		// Roll back if the function returns with err set
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1) Insert the workflow
	res, err := tx.Exec(`
		INSERT INTO WORKFLOW (WORKFLOW_NAME, DESCRIPTION)
		VALUES (?, ?)
	`, wf.WORKFLOW_NAME, wf.WORKFLOW_DESCRIPTION)
	if err != nil {
		return 0, err
	}
	workflowID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 2) Insert initial version as latest (v1)
	// For a brand-new workflow, v1 is always correct.
	_, err = tx.Exec(`
		INSERT INTO WF_VERSION (WORKFLOW_ID, VERSION, IS_ACTIVE, IS_APPROVED)
		VALUES (?, 1, 0, 0)
	`, workflowID)
	if err != nil {
		return 0, err
	}

	// 3) Commit
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return workflowID, nil
}

func DeleteWorkflowByID(id string) error {
		tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.Exec(`DELETE FROM WF_VERSION WHERE WORKFLOW_ID = ?`, id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM WORKFLOW WHERE WORKFLOW_ID = ?`, id)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateWorkflow(wf Workflow) error {
	_, err := db.Exec(`
		UPDATE WORKFLOW
		SET WORKFLOW_NAME = ?, DESCRIPTION = ?
		WHERE WORKFLOW_ID = ?
	`, wf.WORKFLOW_NAME, wf.WORKFLOW_DESCRIPTION, wf.WORKFLOW_ID)
	return err
}



