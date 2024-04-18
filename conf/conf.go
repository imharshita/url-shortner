package conf

import (
	"bytes"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type http struct {
	Listen string `toml:"listen"`
}

var Conf config

type common struct {
	DomainName string `toml:"domain_name"`
	Schema     string `toml:"schema"`
}

type config struct {
	Http   http   `toml:"http"`
	Common common `toml:"common"`
}

func MustParseConfig(configFile string) {
	if fileInfo, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			log.Panicf("configuration file %v does not exist.", configFile)
		} else {
			log.Panicf("configuration file %v can not be stated. %v", configFile, err)
		}
	} else {
		if fileInfo.IsDir() {
			log.Panicf("%v is a directory name", configFile)
		}
	}

	content, err := os.ReadFile(configFile)
	if err != nil {
		log.Panicf("read configuration file error. %v", err)
	}
	content = bytes.TrimSpace(content)

	err = toml.Unmarshal(content, &Conf)
	if err != nil {
		log.Panicf("unmarshal toml object error. %v", err)
	}
}
