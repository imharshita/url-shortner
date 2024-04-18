package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/imharshita/url-shortner/conf"
	"github.com/imharshita/url-shortner/short"
)

func ShortURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"msg":"short url"}`))
	var err error
	if r.Body == nil {
		log.Printf("request body is empty.")
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(errorResp{Msg: http.StatusText(http.StatusBadRequest)})
		w.Write(errMsg)
		return
	}
	defer r.Body.Close()

	// Attempt to decode the request body into a shortReq struct
	var shortReq shortReq
	if err := json.NewDecoder(r.Body).Decode(&shortReq); err != nil {
		log.Printf("parse short request error. %v", err)
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(errorResp{Msg: http.StatusText(http.StatusBadRequest)})
		w.Write(errMsg)
		return
	}

	// Check if the request body satisfies the structure of shortReq
	if shortReq.LongURL == "" {
		log.Printf("LongURL is required.")
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(errorResp{Msg: http.StatusText(http.StatusBadRequest)})
		w.Write(errMsg)
		return
	}

	var longURL *url.URL
	longURL, err = url.Parse(shortReq.LongURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(errorResp{Msg: "requested url is malformed"})
		w.Write(errMsg)
		return
	} else {
		if longURL.Host == conf.Conf.Common.DomainName {
			w.WriteHeader(http.StatusBadRequest)
			errMsg, _ := json.Marshal(errorResp{Msg: "requested url is already shortened"})
			w.Write(errMsg)
			return
		}
		if strings.ToLower(longURL.Scheme) != "http" && strings.ToLower(longURL.Scheme) != "https" {
			w.WriteHeader(http.StatusBadRequest)
			errMsg, _ := json.Marshal(errorResp{Msg: "requested url is not a http or https url"})
			w.Write(errMsg)
			return
		}
	}

	var shortenedURL string
	shortenedURL, err = short.Shorter.Short(shortReq.LongURL)
	shortenedURL = (&url.URL{
		Scheme: conf.Conf.Common.Schema,
		Host:   conf.Conf.Common.DomainName,
		Path:   shortenedURL,
	}).String()
	if err != nil {
		log.Printf("short url error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(errorResp{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	} else {
		shortResp, _ := json.Marshal(shortResp{ShortURL: shortenedURL})
		w.Write(shortResp)
	}
}
