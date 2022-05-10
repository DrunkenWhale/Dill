package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type HibiscusConfig struct {
	Port int    `json:"port"`
	Rule []Rule `json:"rule"`
}

type Rule struct {
	Route   string   `json:"route"`
	Type    string   `json:"type"`
	Methods []string `json:"methods"`
	Dest    string   `json:"dest"`
}

var (
	hibiscus       HibiscusConfig
	defaultMethods = []string{"*"}
)

const (
	defaultPort = 80
	defaultType = "static"
)

func Rules() []Rule {
	return hibiscus.Rule
}

func Port() int {
	return hibiscus.Port
}

func init() {
	file, err := os.OpenFile("hibiscus.json", os.O_RDONLY, 0777)
	if err != nil {
		log.Fatalln(file)
		return
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = json.Unmarshal(bytes, &hibiscus)
	if err != nil {
		log.Fatalln(err)
		return
	}

	if hibiscus.Port == 0 {
		hibiscus.Port = defaultPort
	}

	for _, rule := range hibiscus.Rule {
		if rule.Type != "dynamic" {
			rule.Type = defaultType
		}
		if len(rule.Methods) == 0 {
			rule.Methods = defaultMethods
		}
		for _, method := range rule.Methods {
			if method == "*" {
				rule.Methods = defaultMethods
				break
			}
		}
	}

}
