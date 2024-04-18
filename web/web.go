package web

import (
	"log"
	"net/http"

	"github.com/imharshita/url-shortner/conf"
	"github.com/imharshita/url-shortner/web/api"

	"github.com/gorilla/mux"
)

func Start() {
	log.Println("web starts")
	r := mux.NewRouter()

	r.HandleFunc("/health", api.CheckHealth).Methods(http.MethodGet)
	r.HandleFunc("/short", api.ShortURL).Methods(http.MethodPost).HeadersRegexp("Content-Type", "application/json")
	r.HandleFunc("/expand", api.ExpandURL).Methods(http.MethodPost).HeadersRegexp("Content-Type", "application/json")
	r.HandleFunc("/{shortenedURL:[a-zA-Z0-9]{1,11}}", api.Redirect).Methods(http.MethodGet)

	log.Println(conf.Conf.Http.Listen)
	if err := http.ListenAndServe(conf.Conf.Http.Listen, r); err != nil {
		log.Panicf("web start error. %v", err)
	}
}
