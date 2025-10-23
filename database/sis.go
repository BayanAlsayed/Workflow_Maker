package database

type UserType struct {
	SE_CODE_USER_TYPE_ID int    `json:"se_code_user_type_id"`
	DESCR_EN            string `json:"descr_en"`
}

type Account struct {
	SE_ACCNT_ID int    `json:"se_accnt_id"`
	DESCR_EN   string `json:"descr_en"`
	SE_CODE_USER_TYPE_ID int `json:"se_code_user_type_id"`
}

type EDCodeStatusCat struct {
	ED_CODE_STATUS_CAT_ID int    `json:"ed_code_status_cat_id"`
	DESCR_EN            string `json:"descr_en"`
}

type EDCodeStatus struct {
	ED_CODE_STATUS_ID int    `json:"ed_code_status_id"`
	DESCR_EN        string `json:"descr_en"`
	CAT_ID  int    `json:"cat_id"`
}

type GSCodeReqStatus struct {
	GS_CODE_REQ_STATUS_ID int    `json:"gs_code_req_status_id"`
	DESCR_EN            string `json:"descr_en"`
}

func GetUserTypes () ([]UserType, error) {
	rows, err := db.Query(`SELECT SE_CODE_USER_TYPE_ID, DESCR_EN FROM SE_CODE_USER_TYPE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userTypes []UserType
	for rows.Next() {
		var ut UserType
		if err := rows.Scan(&ut.SE_CODE_USER_TYPE_ID, &ut.DESCR_EN); err != nil {
			return nil, err
		}
		userTypes = append(userTypes, ut)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userTypes, nil
}

func GetAccounts () ([]Account, error) {
	rows, err := db.Query(`SELECT SE_ACCNT_ID, DESCR_EN, SE_CODE_USER_TYPE_ID FROM SE_ACCNT`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var acc Account
		if err := rows.Scan(&acc.SE_ACCNT_ID, &acc.DESCR_EN, &acc.SE_CODE_USER_TYPE_ID); err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func GetEdCodeStatusesCat() ([]EDCodeStatusCat, error) {
	rows, err := db.Query(`SELECT ED_CODE_STATUS_CAT_ID, DESCR_EN FROM ED_CODE_STATUS_CAT`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var edStatusesCat []EDCodeStatusCat
	for rows.Next() {
		var edCat EDCodeStatusCat
		if err := rows.Scan(&edCat.ED_CODE_STATUS_CAT_ID, &edCat.DESCR_EN); err != nil {
			return nil, err
		}
		edStatusesCat = append(edStatusesCat, edCat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return edStatusesCat, nil
}

func GetEdCodeStatuses() ([]EDCodeStatus, error) {
	rows, err := db.Query(`SELECT ED_CODE_STATUS_ID, DESCR_EN, ED_CODE_STATUS_CAT_ID FROM ED_CODE_STATUS`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var edStatuses []EDCodeStatus
	for rows.Next() {
		var ed EDCodeStatus
		if err := rows.Scan(&ed.ED_CODE_STATUS_ID, &ed.DESCR_EN, &ed.CAT_ID); err != nil {
			return nil, err
		}
		edStatuses = append(edStatuses, ed)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return edStatuses, nil
}

func GetGsCodeReqStatuses() ([]GSCodeReqStatus, error) {
	rows, err := db.Query(`SELECT GS_CODE_REQ_STATUS_ID, DESCR_EN FROM GS_CODE_REQ_STATUS`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gsStatuses []GSCodeReqStatus
	for rows.Next() {
		var gs GSCodeReqStatus
		if err := rows.Scan(&gs.GS_CODE_REQ_STATUS_ID, &gs.DESCR_EN); err != nil {
			return nil, err
		}
		gsStatuses = append(gsStatuses, gs)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return gsStatuses, nil
}