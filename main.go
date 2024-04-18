package main

import (
	"flag"

	"github.com/imharshita/url-shortner/conf"
	"github.com/imharshita/url-shortner/short"
	"github.com/imharshita/url-shortner/web"
)

func main() {
	cfgFile := flag.String("c", "config.conf", "configuration file")

	flag.Parse()

	// parse config
	conf.MustParseConfig(*cfgFile)

	// short service
	short.Start()

	// api
	web.Start()
}
