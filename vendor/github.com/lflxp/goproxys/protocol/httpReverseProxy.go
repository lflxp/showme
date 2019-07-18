package protocol

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func RunHttpProxy() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	log.Println("Started Proxy")
	log.Println("Http Reverse Proxy Listening port: 8081")
	demoURL, err := url.Parse("http://127.0.0.1:9528")
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(demoURL)
	// http2.0
	// proxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	// 	req.Host = demoURL.Host
	// 	req.URL.Host = demoURL.Host
	// 	req.URL.Scheme = demoURL.Scheme
	// 	req.RequestURI = ""

	// 	http2.ConfigureTransport(http.DefaultTransport.(*http.Transport))
	// 	response, err := http.DefaultClient.Do(req)
	// 	if err != nil {
	// 		rw.WriteHeader(http.StatusInternalServerError)
	// 		fmt.Fprint(rw, err)
	// 		return
	// 	}
	// 	rw.WriteHeader(response.StatusCode)
	// 	io.Copy(rw, response.Body)
	// })
	err = http.ListenAndServeTLS(":8081", "tls/server.crt", "tls/server.key", proxy)
	if err != nil {
		log.Fatal(err)
	}
}
