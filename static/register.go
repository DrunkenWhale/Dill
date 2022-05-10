package static

import (
	"Hibiscus/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func RegisterStaticResourceServerFromConfigFile() {
	RegisterStaticResourceServer(config.Rules())
}

func RegisterStaticResourceServer(rules []config.Rule) {
	for _, rule := range rules {
		RegisterSingleStaticResourceServer(rule)
	}
}

func RegisterSingleStaticResourceServer(rule config.Rule) {
	http.HandleFunc(rule.Route, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Server", "Hibiscus")
		if len(rule.Methods) == 1 && rule.Methods[0] == "*" {
			staticServerFunctionFactory(rule)(writer, request)
			return
		} else {
			for _, method := range rule.Methods {
				if request.Method == method {
					staticServerFunctionFactory(rule)(writer, request)
					return
				}
			}
		}
		writer.WriteHeader(405)
		_, err := writer.Write([]byte("Method Not Allowed"))
		if err != nil {
			log.Printf("[ %-7v ] %v", "ERROR", err)
			return
		}
	})
}

func staticServerFunctionFactory(rule config.Rule) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("[ %-7v ] request [ %-v ]", "INFO", request.URL)
		path := request.URL.Path
		tempPathSplitArray := strings.SplitN(path, rule.Route, 2)
		if len(tempPathSplitArray) <= 1 {
			writer.WriteHeader(404)
			_, err := writer.Write([]byte("Page Not Found"))
			if err != nil {
				log.Printf("[ %-7v ] %v", "ERROR", err)
				return
			}
			return
		}
		tempPath := tempPathSplitArray[1]
		if tempPath == "" {
			file, err := os.Open(rule.Dest)
			if err != nil {
				log.Printf("[ %-7v ] %v", "ERROR", err)
				writer.WriteHeader(404)
				_, err := writer.Write([]byte("Page Not Found"))
				if err != nil {
					log.Printf("[ %-7v ] %v", "ERROR", err)
				}
				return
			}
			bytes, err := ioutil.ReadAll(file)
			if err != nil {
				log.Printf("[ %-7v ] %v", "ERROR", err)
				return
			}
			_, err = writer.Write(bytes)
			if err != nil {
				log.Printf("[ %-7v ] %v", "ERROR", err)
				return
			}
			return
		}
		var legalAbsolutePath string
		if tempPath[len(tempPath)-1] != '/' {
			legalAbsolutePath = strings.ReplaceAll(rule.Dest+string(os.PathSeparator)+tempPath, "/", string(os.PathSeparator))
		} else {
			legalAbsolutePath = strings.ReplaceAll(rule.Dest+tempPath, "/", string(os.PathSeparator))
		}
		file, err := os.Open(legalAbsolutePath)
		if err != nil {
			log.Printf("[ %-7v ] %v", "ERROR", err)
			writer.WriteHeader(404)
			_, err := writer.Write([]byte("Page Not Found"))
			if err != nil {
				log.Printf("[ %-7v ] %v", "ERROR", err)
			}
			return
		}
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Printf("[ %-7v ] %v", "ERROR", err)
			return
		}
		_, err = writer.Write(bytes)
		if err != nil {
			log.Printf("[ %-7v ] %v", "ERROR", err)
			return
		}
	}
}
