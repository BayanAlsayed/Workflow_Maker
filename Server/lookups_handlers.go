package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	database "workflow/database"
)

func LookupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		RenderTemplate(w, "405.html", nil)
		return
	}

	userTypes, err := database.GetUserTypes()
	if err != nil {
		http.Error(w, "Failed to retrieve user types", http.StatusInternalServerError)
		return
	}

	accounts, err := database.GetAccounts()
	if err != nil {
		http.Error(w, "Failed to retrieve accounts", http.StatusInternalServerError)
		return
	}

	edStatusCodesCat, err := database.GetEdCodeStatusesCat()
	if err != nil {
		http.Error(w, "Failed to retrieve ED status codes", http.StatusInternalServerError)
		return
	}

	edStatusCodes, err := database.GetEdCodeStatuses()
	if err != nil {
		http.Error(w, "Failed to retrieve ED status codes", http.StatusInternalServerError)
		return
	}

	gsStatusCodes, err := database.GetGsCodeReqStatuses()
	if err != nil {
		http.Error(w, "Failed to retrieve GS status codes", http.StatusInternalServerError)
		return
	}

	workflowConditions, err := database.GetAllConditions()
	if err != nil {
		http.Error(w, "Failed to retrieve workflow conditions", http.StatusInternalServerError)
		return
	}

	lookups := struct {
		UserTypes        []database.UserType        `json:"user_types"`
		Accounts         []database.Account         `json:"accounts"`
		EDStatusCodesCat []database.EDCodeStatusCat `json:"ed_status_codes_cat"`
		EDStatusCodes    []database.EDCodeStatus    `json:"ed_status_codes"`
		GSStatusCodes    []database.GSCodeReqStatus `json:"gs_status_codes"`
		WFConditions     []database.Condition       `json:"wf_conditions"`
	}{
		UserTypes:        userTypes,
		Accounts:         accounts,
		EDStatusCodesCat: edStatusCodesCat,
		EDStatusCodes:    edStatusCodes,
		GSStatusCodes:    gsStatusCodes,
		WFConditions:     workflowConditions,

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lookups)
}

func WFLookupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		RenderTemplate(w, "405.html", nil)
		return
	}

	Str := strings.Split(strings.TrimPrefix(r.URL.Path, "/wf_lookups/"), "/")
	WF_ID, err := strconv.Atoi(Str[0])
	if err != nil {
		http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
		return
	}
	version, err := strconv.Atoi(Str[1])
	if err != nil {
		http.Error(w, "Invalid workflow version", http.StatusBadRequest)
		return
	}

	workflowStatuses, err := database.ViewStatuses(WF_ID, version)
	if err != nil {
		http.Error(w, "Failed to retrieve workflow statuses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflowStatuses)
}


