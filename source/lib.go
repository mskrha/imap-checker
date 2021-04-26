package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type configGlobal struct {
	Delimiter string    `json:"delimiter"`
	Servers   []Server  `json:"servers"`
	Accounts  []Account `json:"accounts"`

	servers map[string]Server
}

/*
	Parse the configuration file and preprocess configuration
*/
func parseConfig() error {
	var dcf, cf string
	if h, ok := os.LookupEnv("HOME"); ok {
		dcf = h + "/.imap-checker/config.json"
	} else {
		dcf = "/etc/imap-checker/config.json"
	}

	flag.StringVar(&cf, "config", dcf, "Path to the configuration file")
	flag.Parse()

	/*
		Try to read the configuration file
	*/
	file, err := ioutil.ReadFile(cf)
	if err != nil {
		return err
	}

	/*
		Try to parse the configuration file
	*/
	if err := json.Unmarshal(file, &config); err != nil {
		return err
	}

	/*
		Prepare asociative array with servers
	*/
	config.servers = make(map[string]Server)
	for _, s := range config.Servers {
		if len(s.Name) == 0 {
			return fmt.Errorf("No server name specified")
		}
		if len(s.Host) == 0 {
			return fmt.Errorf("No host specified for server %s", s.Name)
		}
		config.servers[s.Name] = s
	}

	return nil
}
