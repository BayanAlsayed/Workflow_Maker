package server

import (
	"database/sql"
	"fmt"
	"net/http"

	database "workflow/database"
)

func BaseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("base handler")

	defer HandleServerError(w, r)
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		RenderTemplate(w, "404.html", nil)
		return
	}

	workflows, err := database.GetAllWorkflows()
	if err != nil {
		if err == sql.ErrNoRows {
			workflows = []database.Workflow{}
		} else {
			fmt.Println("Error fetching workflows: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			RenderTemplate(w, "500.html", nil)
			return
		}
	}

	RenderTemplate(w, "index.html", workflows)
}
