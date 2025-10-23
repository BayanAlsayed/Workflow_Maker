package database

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Status struct {
	WF_STATUS_ID          int     `json:"wf_status_id"`
	STATUS_ID             int     `json:"status_id"`
	STATUS_NAME           string  `json:"status_name"`
	ED_CODE_STATUS_ID     *int    `json:"ed_code_status_id,omitempty"`
	ED_DESCR_EN           *string `json:"ed_descr_en,omitempty"`
	GS_CODE_REQ_STATUS_ID *int    `json:"gs_code_req_status_id,omitempty"`
	GS_DESCR_EN           *string `json:"gs_descr_en,omitempty"`
	IS_TERMINAL           bool    `json:"is_terminal"`
	SUCCESS_PATH          *int    `json:"success_path,omitempty"`
}


//revise this
func ViewStatuses(WF_ID int) ([]Status, error) {
	//select the newest version for the given workflow
	newestVersion, err := getNewestVersion(WF_ID)
	if err != nil {
		fmt.Println("Error getting newest version: ", err)
		return nil, err
	}
	rows, err := db.Query(`
		SELECT
			s.WF_STATUS_ID,
			s.STATUS_ID,
			s.STATUS_NAME,
			s.ED_CODE_STATUS_ID,
			e.DESCR_EN AS ed_descr_en,
			s.GS_CODE_REQ_STATUS_ID,
			g.DESCR_EN AS gs_descr_en,
			s.IS_TERMINAL,
			s.SUCCESS_PATH
		FROM WF_STATUS AS s
		LEFT JOIN ED_CODE_STATUS     AS e ON s.ED_CODE_STATUS_ID      = e.ED_CODE_STATUS_ID
		LEFT JOIN GS_CODE_REQ_STATUS AS g ON s.GS_CODE_REQ_STATUS_ID  = g.GS_CODE_REQ_STATUS_ID
		WHERE s.WF_VERSION_ID = ?
	`, newestVersion)
	if err != nil {
		if err == sql.ErrNoRows {
			return []Status{}, nil // return empty slice if no statuses found
		} else {
			fmt.Println("Error fetching statuses: ", err)
			return nil, err
		}
	}
	defer rows.Close()

	var statuses []Status
	for rows.Next() {
		var s Status
		var (
			edCodeID    sql.NullInt64
			edDescr     sql.NullString
			gsCodeID    sql.NullInt64
			gsDescr     sql.NullString
			successPath sql.NullInt64
			isTermInt   int
		)

		if err := rows.Scan(
			&s.WF_STATUS_ID,
			&s.STATUS_ID,
			&s.STATUS_NAME,
			&edCodeID,
			&edDescr,
			&gsCodeID,
			&gsDescr,
			&isTermInt,
			&successPath,
		); err != nil {
			fmt.Println("Error scanning status row: ", err)
			return nil, err
		}

		// Map nulls to pointers
		if edCodeID.Valid {
			v := int(edCodeID.Int64)
			s.ED_CODE_STATUS_ID = &v
		}
		if edDescr.Valid {
			v := edDescr.String
			s.ED_DESCR_EN = &v
		}
		if gsCodeID.Valid {
			v := int(gsCodeID.Int64)
			s.GS_CODE_REQ_STATUS_ID = &v
		}
		if gsDescr.Valid {
			v := gsDescr.String
			s.GS_DESCR_EN = &v
		}
		if successPath.Valid {
			v := int(successPath.Int64)
			s.SUCCESS_PATH = &v
		}
		s.IS_TERMINAL = isTermInt != 0

		statuses = append(statuses, s)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Row iteration error: ", err)
		return nil, err
	}
	return statuses, nil
}

func AddStatus(WF_ID int, STATUS_NAME string, IS_TERMINAL string, SUCCESS_PATH string, ED_CODE_STATUS_CAT []string, ED_CODE_STATUS []string, GS_CODE_REQ_STATUS []string) error {
	//select the newest version for the given workflow
	newestVersion, err := getNewestVersion(WF_ID)
	if err != nil {
		fmt.Println("Error getting newest version: ", err)
		return err
	}

	// edCodeStatusCatID, err := getIDFromStringSlice(ED_CODE_STATUS_CAT)
	// if err != nil {
	// 	return err
	// }

	edCodeStatusID, err := getIDFromStringSlice(ED_CODE_STATUS)
	if err != nil {
		fmt.Println("Error parsing ED code status ID: ", err)
		return err
	}

	gsCodeReqStatusID, err := getIDFromStringSlice(GS_CODE_REQ_STATUS)
	if err != nil {
		fmt.Println("Error parsing GS code req status ID: ", err)
		return err
	}

	statusID, err := getNextStatusID(newestVersion)
	if err != nil {
		fmt.Println("Error getting next status ID: ", err)
		return err
	}

	var SUCCESS_PATH_Int *int
	if SUCCESS_PATH != "" {
		sp, err := strconv.Atoi(SUCCESS_PATH)
		if err != nil {
			fmt.Println("Error converting SUCCESS_PATH to int: ", err)
			return err
		}
		SUCCESS_PATH_Int = &sp
	} else {
		SUCCESS_PATH_Int = nil
	}

	_, err = db.Exec(`
		INSERT INTO WF_STATUS (
			STATUS_ID,
			WF_VERSION_ID,
			STATUS_NAME,
			ED_CODE_STATUS_ID,
			GS_CODE_REQ_STATUS_ID,
			IS_TERMINAL,
			SUCCESS_PATH
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`, statusID, newestVersion, STATUS_NAME, edCodeStatusID, gsCodeReqStatusID, IS_TERMINAL, SUCCESS_PATH_Int)
	if err != nil {
		fmt.Println("Error inserting new status: ", err)
		return err
	}
	return nil

}

func getIDFromStringSlice(input []string) (*int, error) {
	if len(input) == 0 || input[0] == "" {
		return nil, nil
	}
	idStr := input[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func getNextStatusID(WF_VERSION_ID int) (int, error) {
	var nextID int
	err := db.QueryRow(`
		SELECT IFNULL(MAX(STATUS_ID), 0) + 1 AS next_id
		FROM WF_STATUS
		WHERE WF_VERSION_ID = ?
	`, WF_VERSION_ID).Scan(&nextID)
	if err != nil {
		return 0, err
	}
	return nextID, nil
}
