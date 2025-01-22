package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {

	// This request path forwards the request to serve B
	http.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		// get hostname, (in the case its running in docker return the container id)
		hname, err := os.Hostname()

		// send back that the server is sending a request to server B
		fmt.Fprintf(w, "server %s sending request to server B\n", hname)
		// get response from server b which forward to server c

		resp, err := http.Get("http://b/b")

		// if we got an error send back the error
		if err != nil {
			fmt.Fprintf(w, "Got an error %s", err)
		} else {
			// copy the resposne from b and send back
			io.Copy(w, resp.Body)
		}
		// print that we are done
		fmt.Fprintf(w, "\nclosing")
	})

	// This request path forwards the request to serve C, works the same way as /a
	http.HandleFunc("/b", func(w http.ResponseWriter, r *http.Request) {

		hname, err := os.Hostname()
		fmt.Fprintf(w, "server %s sending request to server C\n", hname)
		resp, err := http.Get("http://c/c")
		if err != nil {
			fmt.Fprintf(w, "Got an error %s", err)
		} else {
			io.Copy(w, resp.Body)
		}
		fmt.Fprintf(w, "\nclosing")
	})

	// acting as the final node,  this does not forward the request
	// and only responds with its hostname
	http.HandleFunc("/c", func(w http.ResponseWriter, r *http.Request) {
		hname, _ := os.Hostname()
		fmt.Fprintf(w, "server %s end", hname)
	})

	// create the server and listen to port 80
	http.ListenAndServe(":80", nil)
}
