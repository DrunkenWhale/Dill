package main

import (
	"Hibiscus/data"
	"Hibiscus/register"
	"log"
	"net/http"
)

const port string = "80"

func main() {
	register.RegisterReverseProxy(data.Data)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
