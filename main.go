package main

import (
	"Hibiscus/config"
	"Hibiscus/static"
	"log"
	"net/http"
	"strconv"
)

func main() {
	static.RegisterStaticResourceServerFromConfigFile()
	err := http.ListenAndServe(":"+strconv.Itoa(config.Port()), nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
