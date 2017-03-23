package main

import (
	// "fmt"
	// "io"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func webServer() {
	defer wg.Done()
	http.HandleFunc("/hello", HelloServer)
	info("Starting web server...")
	var server http.Server
	server.Addr = ":8443"
	err := server.ListenAndServeTLS("server.crt", "server.key")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	/*
		go func() {
			select {
			case <-exiting:
				server.Close()
			}
		}()
	*/
}
