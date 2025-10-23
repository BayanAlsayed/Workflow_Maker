package database

import "fmt"

func getNewestVersion(WF_ID int) (int, error) {
	//select the newest version for the given workflow
	var newestVersion int
	err := db.QueryRow(`
		SELECT WF_VERSION_ID
		FROM WF_VERSION
		WHERE WORKFLOW_ID = ? AND IS_LATEST = 1
	`, WF_ID).Scan(&newestVersion)
	if err != nil {
		fmt.Println("Error fetching newest version: ", err)
		return 0, err
	}
	return newestVersion, nil
}
