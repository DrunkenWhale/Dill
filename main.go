package main

import (
	"Hibiscus/config"
	"Hibiscus/static"
	"log"
	"net/http"
	"strconv"
)

func main() {
	log.Printf("[%v] %-v", "INFO", "Init Server")
	static.RegisterStaticResourceServerFromConfigFile()
	log.Printf("[%v] %-v", "INFO", "Start Server")
	err := http.ListenAndServe(":"+strconv.Itoa(config.Port()), nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
