package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"pets/parse"
	"pets/pathExpression"

	"github.com/TyphonHill/go-mermaid/diagrams/flowchart"
)

type RequestData struct {
	Data string `json:"data"`
}

type ResponseData struct {
	Message string `json:"message"`
}

var nodeLst = map[string][]pathExpression.DataEdge{} // NODE HASHMAP WITH A TUPLE LIST (EDGES) AS VALUE

func main() {
	nodeLst = parse.Parse()

	// EXAMPLE REMOVE ME LATER
	q, _ := pathExpression.BobTheBuilder("Pickaxe_Instance_Henry/{obtainedBy/hasInput}*", nodeLst)
	s := pathExpression.TraverseQuery(&q, nodeLst)
	println(s)

	// END EXAMPLE
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
		fmt.Println("hnadleSubmit invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("handleSubmit error reading request body")
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// parse the request data into requestData struct
	var requestData RequestData
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		fmt.Println("handleSubmit error parsing JSON")
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	if requestData.Data == "" {
		fmt.Println("handleSubmit error empty request data")
		http.Error(w, "empty request data", http.StatusBadRequest)
		return
	}

	res, err := sendQuery(requestData.Data)
	if err != nil {
		fmt.Println("error processing query")
		response := ResponseData{Message: "error: invalid query format"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Println("res: " + res)

	//if requestData.Data == "mermaid" {
	//	mermaid(w, r)
	//	return
	//}

	fmt.Println("Received from client:", res)

	// create a response containging the recieived data
	response := ResponseData{Message: res}

	// set response content type to JSON and send it back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendQuery(queryString string) (string, error) {
	if queryString == "" {
		fmt.Println("error empty query string")
		return "", fmt.Errorf("empty query string")
	}
	q, _ := pathExpression.BobTheBuilder(queryString, nodeLst)

	s := pathExpression.TraverseQuery(&q, nodeLst)
	println("s: ", s)
	return s, nil
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
