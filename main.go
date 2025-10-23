package main

import (
	"database/sql"
	"fmt"
	"net/http"

	server "workflow/Server"
	database "workflow/database"
)

func main() {
	db, err := sql.Open("sqlite3", "./database/workflow.db")
	if err != nil {

		fmt.Println("Failed to open database connection: ", err)
		return
	}
	defer db.Close()

	database.SetDB(db)
	err = database.CreateTables()
	if err != nil {
		fmt.Println("error in creating tables: ", err)
		return
	}

	// err = database.SeedUserTypes()
	// if err != nil {
	// 	fmt.Println("error in seeding user types: ", err)
	// 	return
	// }

	// if err := database.SeedAccounts(); err != nil {
	// 	fmt.Println("error in seeding accounts: ", err)
	// 	return
	// }

	// if err := database.SeedEdCodeStatusCat(); err != nil {
	// 	fmt.Println("error in seeding ED_CODE_STATUS_CAT: ", err)
	// 	return
	// }

	// if err := database.SeedEdCodeStatus(); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := database.SeedGsCodeReqStatus(); err != nil {
	// 	log.Fatal(err)
	// }

	mux := http.NewServeMux()

	mux.HandleFunc("/", server.BaseHandler)
	mux.HandleFunc("/add_workflow", server.AddWorkflowHandler)
	mux.HandleFunc("/delete_workflow", server.DeleteWorkflowHandler)
	mux.HandleFunc("/edit_workflow", server.EditWorkflowHandler)

	mux.HandleFunc("/view_workflow/", server.ViewDetailsHandler)
	mux.HandleFunc("/lookups/", server.LookupsHandler)

	mux.HandleFunc("/add_status", server.AddStatusHandler)

	mux.HandleFunc("/add_rule", server.AddRuleHandler)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))

	fmt.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", corsMiddleware(mux)); err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Modify to specific origins as needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Enable if you need to send cookies

		// Handle preflight checks
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
