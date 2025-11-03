package database

import "fmt"

type Version struct {
	VERSION    int  `json:"version"`
	IS_ACTIVE  bool `json:"is_active"`
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
