package main

import (
	"bytes"
	"encoding/binary"
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
	"github.com/google/uuid"
)

type RequestData struct {
	Data string `json:"data"`
}

type ResponseData struct {
	Message string `json:"message"`
}

var nodeLst = map[string][]parse.DataEdge{} // NODE HASHMAP WITH A TUPLE LIST (EDGES) AS VALUE

func main() {
	nodeLst = parse.Parse()
	for k, v := range nodeLst {
		fmt.Printf("%s: %v\n", k, v)
		fmt.Println("-----------------")
	}
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

	header := make([]byte, 0, 32)
	header = binary.BigEndian.AppendUint16(header, 100)
	qid, _ := uuid.New().MarshalBinary()
	header = append(header, qid...)

	stream := io.MultiReader(bytes.NewReader(header), r.Body)

	q, _ := pathExpression.QueryStructFromStream(&stream)
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

	// create a response containing the received data
	response := ResponseData{Message: res}

	// set response content type to JSON and send it back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// this functions creates an header with an ttl = 100 and an uuid
func sendQuery(queryString string) string {

	// create a header of ttl = 100 and an uuid
	header := make([]byte, 0, 32)
	header = binary.BigEndian.AppendUint16(header, 100)
	qid, _ := uuid.New().MarshalBinary()
	header = append(header, qid...)

	payload := strings.NewReader(queryString)

	stream := io.MultiReader(bytes.NewReader(header), payload)

	q, _ := pathExpression.QueryStructFromStream(&stream)

	s := pathExpression.TraverseQuery(&q, nodeLst)
	println(s)
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
