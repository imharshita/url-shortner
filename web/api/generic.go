package api

import (
	"html/template"
	"log"
	"net/http"
)

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	tpl := template.New("health.html")
	var err error
	tpl, err = tpl.ParseFiles("template/health.html")
	if err != nil {
		log.Printf("parse template error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		log.Printf("execute template error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	}
}
