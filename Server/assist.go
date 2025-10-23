package server

import (
	"fmt"
	"net/http"
	"text/template"
)

func HandleServerError(w http.ResponseWriter, r *http.Request) {
	if rec := recover(); rec != nil {
		w.WriteHeader(http.StatusInternalServerError)
		RenderTemplate(w, "500.html", nil)
		return
	}
}

func RenderTemplate(w http.ResponseWriter, templatePath string, data interface{}) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println(err, "line 20, utils/handler_helper.go")
		w.WriteHeader(http.StatusInternalServerError)
		RenderTemplate(w, "500.html", nil)
		return
	}
	if err = tmpl.Execute(w, data); err != nil {
		fmt.Println("Error executing template: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		RenderTemplate(w, "500.html", nil)
		return
	}
}
