package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	"github.com/imharshita/url-shortner/conf"
	"github.com/imharshita/url-shortner/short"
)

func ShortURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	err = json.NewDecoder(r.Body).Decode(&shortReq)
	if err != nil {
		log.Printf("parse short request error. %v", err)
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(errorResp{Msg: err.Error() + ": " + http.StatusText(http.StatusBadRequest)})
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
		errMsg, _ := json.Marshal(errorResp{Msg: err.Error() + ": " + http.StatusText(http.StatusBadRequest)})
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
		errMsg, _ := json.Marshal(errorResp{Msg: err.Error() + ": " + http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	} else {
		shortResp, _ := json.Marshal(shortResp{ShortURL: shortenedURL})
		w.Write(shortResp)
	}
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortededURL := vars["shortenedURL"]

	longURL, err := short.Shorter.Expand(shortededURL)
	if err != nil {
		log.Printf("redirect short url error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	} else {
		if len(longURL) != 0 {
			w.Header().Set("Location", longURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func ExpandURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
	var expandReq expandReq
	err = json.NewDecoder(r.Body).Decode(&expandReq)
	if err != nil {
		log.Printf("parse short request error. %v", err)
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(errorResp{Msg: err.Error() + ": " + http.StatusText(http.StatusBadRequest)})
		w.Write(errMsg)
		return
	}

	// Check if the request body satisfies the structure of shortReq
	if expandReq.ShortURL == "" {
		log.Printf("ShortURL is required.")
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(errorResp{Msg: http.StatusText(http.StatusBadRequest)})
		w.Write(errMsg)
		return
	}

	var shortURL *url.URL
	shortURL, err = url.Parse(expandReq.ShortURL)
	if err != nil {
		log.Printf(`"%v" is not a valid url`, expandReq.ShortURL)
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(errorResp{Msg: err.Error() + ": " + http.StatusText(http.StatusBadRequest)})
		w.Write(errMsg)
		return
	} else {
		var expandedURL string
		expandedURL, err = short.Shorter.Expand(strings.TrimLeft(shortURL.Path, "/"))
		if err != nil {
			log.Printf("expand url error. %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			errMsg, _ := json.Marshal(errorResp{Msg: err.Error() + ": " + http.StatusText(http.StatusInternalServerError)})
			w.Write(errMsg)
			return
		} else {
			expandResp, _ := json.Marshal(expandResp{LongURL: expandedURL})
			w.Write(expandResp)
		}
	}
}

func Metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Get the count of each domain
	domainCounts := make(map[string]int)
	for _, url := range short.Shorter.GetShortenedURLs() {
		domain := getDomainFromURL(url)
		domainCounts[domain]++
	}

	// Sort the domain counts in descending order
	sortedCounts := make([]domainCount, 0, len(domainCounts))
	for domain, count := range domainCounts {
		sortedCounts = append(sortedCounts, domainCount{Domain: domain, Count: count})
	}
	sort.Slice(sortedCounts, func(i, j int) bool {
		return sortedCounts[i].Count > sortedCounts[j].Count
	})

	// Get the top 3 domain names
	var top3Domains []domainCount
	if len(sortedCounts) >= 3 {
		top3Domains = sortedCounts[:3]
	} else {
		top3Domains = sortedCounts
	}

	for _, count := range top3Domains {
		log.Printf("%s: %v", count.Domain, count.Count)
	}

	// Write the response
	jsonResp, err := json.Marshal(top3Domains)
	if err != nil {
		log.Printf("error marshaling metrics response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func getDomainFromURL(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("error parsing URL: %v", err)
		return ""
	}

	hostParts := strings.Split(parsedURL.Hostname(), ".")
	// Extract the subdomain (e.g., "www") and the domain (e.g., "youtube")
	var domain string
	if len(hostParts) >= 2 {
		domain = hostParts[len(hostParts)-2]
	} else {
		domain = parsedURL.Hostname()
	}

	return domain
}
