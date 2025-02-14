package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"pets/parse"
	"pets/pathExpression"
	"strings"

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

	http.HandleFunc("/", handler) // servers the main HTML file

	http.HandleFunc("/api/submit", handleSubmit) // API endpoint to handle form submission

	http.HandleFunc("/api/recq", queryHandler)

	// create the server and listen to port 80
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
	}
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	b := strings.Builder{}
	full, _ := io.ReadAll(r.Body)
	b.Write(full)
	log.Printf("Received request: \n%s", b.String())

	q, _ := pathExpression.BobTheBuilder(b.String())
	log.Printf("parsed request: \n%s", q.DebugToString())
	pathExpression.RecursiveTraverse(&q, nodeLst, w)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// return 404 for other paths
	//if r.URL.Path != "/" {
	//	http.NotFound(w, r)
	//	return
	//}

	// reads the html file
	html, err := os.ReadFile("./src/index.html")
	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
		tmp, _ := os.Getwd()
		fmt.Fprintf(w, "PWD %s", tmp)
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

	// create a response containging the recieived data
	response := ResponseData{Message: res}

	// set response content type to JSON and send it back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendQuery(queryString string) string {

	q, _ := pathExpression.BobTheBuilder(queryString)
	s := pathExpression.TraverseQuery(&q, nodeLst)
	// println(s)
	return s
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
