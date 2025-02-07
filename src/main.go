package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/TyphonHill/go-mermaid/diagrams/flowchart"
)

type DataEdge struct {
	EdgeName   string
	TargetName string
}

type RequestData struct {
	Data string `json:"data"`
}

type ResponseData struct {
	Message string `json:"message"`
}

func main() {

	//http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprint(w, sendQuery("hello"))
	//})

	http.HandleFunc("/", handler) // servers the main HTML file

	http.HandleFunc("/api/submit", handleSubmit) // API endpoint to handle form submission

	fmt.Printf("Server is running on http://localhost:80")

	// create the server and listen to port 80
	http.ListenAndServe(":80", nil)
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

	//if requestData.Data == "mermaid" {
	//	mermaid(w, r)
	//	return
	//}
	res := sendQuery(requestData.Data)
	fmt.Println("Received from client:", res)

	// create a response containging the recieived data
	response := ResponseData{Message: res}

	// set response content type to JSON and send it back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func TraverseQuery(data map[string][]DataEdge) string {
	return "pickaxe-->|obtainedBy|Pickaxe_From_Stick_And_Stone_Recipe\nPickaxe_From_Stick_And_Stone_Recipe-->|hasInput|Stick\nPickaxe_From_Stick_And_Stone_Recipe-->|hasInput|Cobblestone"
}

func sendQuery(queryString string) string {
	//q, _ := bobTheBuilder(queryString, data)

	ret := TraverseQuery(nil)
	return ret
}

func mermaid(w http.ResponseWriter, r *http.Request) {
	fc := flowchart.NewFlowchart()
	fc.Title = "Mermaid"

	node1 := fc.AddNode("Start")
	node2 := fc.AddNode("End")

	link := fc.AddLink(node1, node2)
	link.Shape = flowchart.LinkShapeDotted

	diagram := fc.String()

	response := ResponseData{Message: diagram}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
