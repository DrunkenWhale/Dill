package register

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func RegisterReverseProxy(proxyData map[string]string) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		client := http.Client{}
		forwardPath, ok := proxyData[request.URL.Path]
		if !ok {
			log.Printf("unregister path: [%v]", request.URL.Path)
			return
		}
		req, err := http.NewRequest(request.Method, forwardPath, request.Body)
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
		_, err = writer.Write(bytes)
		if err != nil {
			log.Print(err)
			return
		}
		fmt.Printf("forward from %v to %v", request.URL, req.URL)
	})
}
