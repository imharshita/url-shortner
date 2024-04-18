# url-shortner
URL Shortener

### Introduction
----
URL Shortener is a url shortening service written in Golang.  

### Features
----
* If the user again ask for the same URL, itshould give me the same URL as it gave before instead
of generating a new one.

* If the user clicks on the short URL then he should be redirected to the original URL. Write
a Redirection API that implements thisfunctionality. 

* The URL and shortened URL is stored in memory by application.

* Includes a metrics API that returns top 3 domain names that have been shortened the most
number of times. For eg. if the user has shortened 2 YouTube video links and 1 Wikipedia links.
Then the output would be:
youtube: 2
wikipedia: 1

### Api
----
* `/health`
    * `HTTP GET`
    * Health check
    * Example
        * `curl http://127.0.0.1:3030/health`
* `/short`
    * `HTTP POST`
    * Short the long url
    * Example
        * `curl -X POST -H "Content-Type:application/json" -d "{\"longURL\": \"http://www.google.com\"}" http://127.0.0.1:3030/short`
* `/expand`
    * `HTTP POST`
    * Expand the short url
    * Example
        * `curl -X POST -H "Content-Type:application/json" -d "{\"shortURL\": \"http://127.0.0.1:3030/ed646a3\"}" http://127.0.0.1:3030/expand`
* `/metrics`
    * `HTTP GET`
    * Get the metrics
    * Example
        * `curl http://127.0.0.1:3030/metrics`
* Redirect
    * `HTTP GET`
    * Redirect short URL to the long/original URL
    * Click on short URL to get the original one