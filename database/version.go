package database

import (
	"database/sql"
	"fmt"
)

type Version struct {
	VERSION     int  `json:"version"`
	IS_ACTIVE   bool `json:"is_active"`
	IS_APPROVED bool `json:"is_approved"`
}

func getVersionID(WF_ID int, version int) (int, error) {
	//select the newest version for the given workflow
	fmt.Println("getting WF_VERSION_ID for workflow: ", WF_ID, " and Version: ", version)
	var versionID int
	err := db.QueryRow(`
		SELECT WF_VERSION_ID
		FROM WF_VERSION
		WHERE WORKFLOW_ID = ? AND VERSION = ?
	`, WF_ID, version).Scan(&versionID)
	if err != nil {
		fmt.Println("Error fetching newest version: ", err)
		return 0, err
	}
	return versionID, nil
}

func GetWorkflowVersions(WF_ID int) ([]Version, error) {
	var versions []Version
	rows, err := db.Query(`
		SELECT VERSION, IS_ACTIVE, IS_APPROVED
		FROM WF_VERSION
		WHERE WORKFLOW_ID = ?
	`, WF_ID)
	if err != nil {
		fmt.Println("Error fetching versions: ", err)
		return nil, err
	}

	for rows.Next() {
		var v Version
		if err := rows.Scan(&v.VERSION, &v.IS_ACTIVE, &v.IS_APPROVED); err != nil {
			fmt.Println("Error scanning version: ", err)
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, nil
}

func ActivateWorkflowVersion(WF_ID int, VERSION int) error {
	// First, set all versions of the workflow to inactive

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.Exec(`
		UPDATE WF_VERSION
		SET IS_ACTIVE = 0
		WHERE WORKFLOW_ID = ?
	`, WF_ID)
	if err != nil {
		fmt.Println("Error activating version: ", err)
		return err
	}

	_, err = tx.Exec(`
		UPDATE WF_VERSION
		SET IS_ACTIVE = 1
		WHERE WORKFLOW_ID = ? AND VERSION = ?
	`, WF_ID, VERSION)
	if err != nil {
		fmt.Println("Error activating version: ", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func ApproveWorkflowVersion(WF_ID int, VERSION int) error {
	_, err := db.Exec(`
		UPDATE WF_VERSION
		SET IS_APPROVED = 1
		WHERE WORKFLOW_ID = ? AND VERSION = ?
	`, WF_ID, VERSION)
	if err != nil {
		fmt.Println("Error approving version: ", err)
		return err
	}
	return nil
}

func DeleteWorkflowVersion(WORKFLOW_ID int, VERSION int) error {
	_, err := db.Exec(`
		DELETE FROM WF_VERSION
		WHERE WORKFLOW_ID = ? AND VERSION = ?
	`, WORKFLOW_ID, VERSION)
	if err != nil {
		fmt.Println("Error deleting version: ", err)
		return err
	}
	return nil
} 

func AddVersion(workflowID int) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Compute next version (1 if none exist)
	var maxVersion sql.NullInt64
	if err = tx.QueryRow(`SELECT MAX(VERSION) FROM WF_VERSION WHERE WORKFLOW_ID = ?`, workflowID).Scan(&maxVersion); err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	next := 1
	if maxVersion.Valid {
		next = int(maxVersion.Int64) + 1
	}

	res, err := tx.Exec(`
		INSERT INTO WF_VERSION (WORKFLOW_ID, VERSION, IS_ACTIVE, IS_APPROVED)
		VALUES (?, ?, 0, 0)
	`, workflowID, next)
	if err != nil {
		return 0, err
	}

	newID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return newID, nil
}

func DuplicateWorkflowVersion(WF_ID int, VERSION int) error {
	NEW_WF_VERSION_ID, err := AddVersion(WF_ID)
	if err != nil {
		fmt.Println("error adding version", err)
		return err
	}

	OLD_WF_VERSION_ID, err := getVersionID(WF_ID, VERSION)
	if err != nil {
		fmt.Println("error getting WF version", err)
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		// ensure rollback on panic
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	_, err = tx.Exec(`
		INSERT INTO WF_STATUS (STATUS_ID, WF_VERSION_ID, STATUS_NAME, ED_CODE_STATUS_ID, GS_CODE_REQ_STATUS_ID, IS_TERMINAL, SUCCESS_PATH)
		SELECT STATUS_ID, ?, STATUS_NAME, ED_CODE_STATUS_ID, GS_CODE_REQ_STATUS_ID, IS_TERMINAL, SUCCESS_PATH
		FROM WF_STATUS
		WHERE WF_VERSION_ID = ?
	`, NEW_WF_VERSION_ID, OLD_WF_VERSION_ID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("copy WF_STATUS failed: %w", err)
	}

	// 5) copy WF_RULE rows (do NOT copy WF_RULE_ID autoinc PK)
	// Make sure column list matches your schema (RULE_ID, WF_VERSION_ID, FROM_STATUS_ID, TO_STATUS_ID, ...)
	_, err = tx.Exec(`
		INSERT INTO WF_RULE (RULE_ID, WF_VERSION_ID, FROM_STATUS_ID, TO_STATUS_ID, SE_CODE_USER_TYPE_ID, SE_ACCNT_ID, CALENDAR_ID, ACTION_BUTTON, ACTION_FUNCTION)
		SELECT RULE_ID, ?, FROM_STATUS_ID, TO_STATUS_ID, SE_CODE_USER_TYPE_ID, SE_ACCNT_ID, CALENDAR_ID, ACTION_BUTTON, ACTION_FUNCTION
		FROM WF_RULE
		WHERE WF_VERSION_ID = ?
	`, NEW_WF_VERSION_ID, OLD_WF_VERSION_ID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("copy WF_RULE failed: %w", err)
	}

	// 6) copy WF_RULE_CONDITION rows (do NOT copy WF_RULE_CONDITION_ID autoinc PK)
	_, err = tx.Exec(`
		INSERT INTO WF_RULE_CONDITION (RULE_CONDITION_ID, WF_VERSION_ID, RULE_ID, WF_CONDITION_ID, TYPE, EXPECTED_VALUE)
		SELECT RULE_CONDITION_ID, ?, RULE_ID, WF_CONDITION_ID, TYPE, EXPECTED_VALUE
		FROM WF_RULE_CONDITION
		WHERE WF_VERSION_ID = ?
	`, NEW_WF_VERSION_ID, OLD_WF_VERSION_ID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("copy WF_RULE_CONDITION failed: %w", err)
	}

	// commit
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("tx commit failed: %w", err)
	}

	return nil
}
