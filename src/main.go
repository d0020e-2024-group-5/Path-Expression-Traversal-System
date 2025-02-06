package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type RequestData struct {
	Data string `json:"data"`
}

type ResponseData struct {
	Message string `json:"message"`
}

func main() {

	// This request path forwards the request to serve B
	http.HandleFunc("/contact_b", func(w http.ResponseWriter, r *http.Request) {
		// get hostname, (in the case its running in docker return the container id)
		hname, err := os.Hostname()

		// send back that the server is sending a request to server B
		fmt.Fprintf(w, "server %s sending request to server B\n", hname)
		// get response from server b which forward to server c

		resp, err := http.Get("http://b/contact_c")

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

	http.HandleFunc("/", handler) // servers the main HTML file

	http.HandleFunc("/api/submit", handleSubmit) // API endpoint to handle form submission

	// create the server and listen to port 80
	http.ListenAndServe(":80", nil)

	fmt.Printf("Server is running on http://localhost:80")
}

func handler(w http.ResponseWriter, r *http.Request) {
	// return 404 for other paths
	//if r.URL.Path != "/" {
	//	http.NotFound(w, r)
	//	return
	//}

	// reads the html file
	html, err := os.ReadFile("../website/index.html")
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	// set the response content type to HTML and write file content
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "%s", html)
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {

	// checking request method is POST
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// parse the request data into requestData struct
	var requestData RequestData
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("Received from client:", requestData.Data)

	// create a response containging the recieived data
	response := ResponseData{Message: "You entered: " + requestData.Data}

	// set response content type to JSON and send it back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
