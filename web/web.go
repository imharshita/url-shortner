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

	log.Println(conf.Conf.Http.Listen)
	if err := http.ListenAndServe(conf.Conf.Http.Listen, r); err != nil {
		log.Panicf("web start error. %v", err)
	}
}
