package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"pets/parse"
	"pets/pathExpression"
)

var nodeLst = map[string][]pathExpression.DataEdge{} // NODE HASHMAP WITH A TUPLE LIST (EDGES) AS VALUE

func main() {
	nodeLst = parse.Parse()

	// EXAMPLE REMOVE ME LATER
	q, _ := pathExpression.BobTheBuilder("Pickaxe_Instance_Henry/{obtainedBy/hasInput}*")
	s := pathExpression.TraverseQuery(&q, nodeLst)
	println(s)
	// END EXAMPLE

	// This request path forwards the request to serve B
	http.HandleFunc("/contact_b", func(w http.ResponseWriter, r *http.Request) {
		// get hostname, (in the case its running in docker return the container id)
		hname, err := os.Hostname()
		fmt.Println(hname)
		fmt.Println("test")
		// send back that the server is sending a request to server B
		fmt.Fprintf(w, "server %s sending request to server B\n", hname)
		// get response from server b which forward to server c

		resp, err := http.Get("http://b/contact_c") //Pickaxe/obtainedBy

		// if we got an error send back the error
		if err != nil {
			fmt.Fprintf(w, "Got an error %s", err)
		} else {
			// copy the response from b and send back
			io.Copy(w, resp.Body)
		}
		// print that we are done
		fmt.Fprintf(w, "\nclosing")
	})

	// This request path forwards the request to serve C, works the same way as /a
	http.HandleFunc("/contact_c", func(w http.ResponseWriter, r *http.Request) {

		hname, err := os.Hostname()
		fmt.Fprintf(w, "server %s sending request to server C\n", hname)
		resp, err := http.Get("http://c/return")
		if err != nil {
			fmt.Fprintf(w, "Got an error %s", err)
		} else {
			io.Copy(w, resp.Body)
		}
		fmt.Fprintf(w, "\nclosing")
	})

	// acting as the final node,  this does not forward the request
	// and only responds with its hostname
	http.HandleFunc("/return", func(w http.ResponseWriter, r *http.Request) {
		hname, _ := os.Hostname()
		fmt.Fprintf(w, "server %s return", hname)
	})

	// create the server and listen to port 80
	http.ListenAndServe(":80", nil)
}
