package register

import (
	"io/ioutil"
	"log"
	"net/http"
)

var client = http.Client{
	Transport: &http.Transport{
		MaxIdleConns: 20,
	},
}

func ReverseProxyRegister(proxyData map[string]string) {
	for k, v := range proxyData {
		registerSingleReverseProxy(k, v)
	}
}

func registerSingleReverseProxy(src string, dst string) {
	http.HandleFunc(src, func(writer http.ResponseWriter, request *http.Request) {
		req, err := http.NewRequest(request.Method, dst, request.Body)
		if err != nil {
			log.Print(err)
			return
		}
		res, err := client.Do(req)
		if err != nil {
			log.Print(err)
			return
		}
		for k, v := range res.Header {
			writer.Header().Set(k, v[0])
		}
		bytes, err := ioutil.ReadAll(res.Body)
		err = res.Body.Close()
		if err != nil {
			log.Fatalln(err)
			return
		}
		_, err = writer.Write(bytes)
		if err != nil {
			log.Print(err)
			log.Printf("can't forward from %v to %v", src, dst)
			return
		}
		log.Printf("forward from %v to %v", src, dst)
	})
}
