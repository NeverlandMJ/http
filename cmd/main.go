package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	host := "0.0.0.0"
	port := "9999"

	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}

func execute(host string, port string) (err error) {
	srv := &http.Server{
		Addr: net.JoinHostPort(host, port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Print(r.RequestURI) // full URI
			log.Print(r.Method)     // method

			body, err := ioutil.ReadAll(r.Body) // request body
			if err != nil {
				log.Print(err)
			}
			log.Printf("%s", body)

		}),
	}

	return srv.ListenAndServe()
}

// type handler struct {
// 	mu       *sync.RWMutex
// 	handlers map[string]http.HandlerFunc
// }

// func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
// 	h.mu.RLock()
// 	handler, ok := h.handlers[request.URL.Path]
// 	h.mu.RUnlock()

// 	if !ok {
// 		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	handler(writer, request)

// }

// func execute(host string, port string) (err error) {
// 	mux := http.NewServeMux()
// 	bannersSvc := banners.NewService()

// 	server := app.NewServer(mux, bannersSvc)
// 	server.Init()

// 	srv := &http.Server{
// 		Addr:    net.JoinHostPort(host, port),
// 		Handler: server,
// 	}
// 	log.Print("server start" + host + ":" + port)
// 	return srv.ListenAndServe()
// }
