package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"

	"io"
	"log"
	"net/http"
	"os"
	"pets/dbComm"
	"pets/pathExpression"

	"strconv"
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

var prefixList = []string{"minecraft: <http://example.org/minecraft#>"}

//var nodeLst = map[string][]dbComm.DataEdge{} // NODE HASHMAP WITH A TUPLE LIST (EDGES) AS VALUE

func main() {

	fmt.Println(dbComm.DBGetNodeEdgesString("minecraft:obtainedBy", prefixList))
	http.HandleFunc("/", handler) // servers the main HTML file

	http.HandleFunc("/api/submit", handleSubmit) // API endpoint to handle form submission

	http.HandleFunc("/api/pets", queryHandler)

	// create the server and listen to port 80
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
	}
}

func queryHandler(w http.ResponseWriter, r *http.Request) {

	// read the first 4 bytes
	var redMagic [4]byte
	_, err := io.ReadFull(r.Body, redMagic[:])
	if err != nil {
		fmt.Fprintf(w, "%% error reading: %s", err.Error())
		return
	}

	// if this magic doesn't match
	if !reflect.DeepEqual(redMagic[:], pathExpression.PetsMermaidQueryHeader[:4]) {
		fmt.Fprint(w, "%% Bad magic")
		return
	}

	// read the type
	var petsType uint16
	err = binary.Read(r.Body, binary.BigEndian, &petsType)
	if err != nil {
		fmt.Fprintf(w, "%% error reading type: %s", err.Error())
		return
	}

	// if the type is not recursive mermaid, then write the error
	if petsType != 1 {
		fmt.Fprint(w, "%% types other than recursive mermaid is not implemented")
		return
	}

	stream := io.Reader(r.Body)
	q, _ := pathExpression.QueryStructFromStream(&stream)
	log.Printf("parsed request: \n%s", q.DebugToString())
	pathExpression.RecursiveTraverse(&q, w)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// return 404 for other paths
	//if r.URL.Path != "/" {
	//	http.NotFound(w, r)
	//	return
	//}

	// reads the html file
	html, err := os.ReadFile("index.html")
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

	fmt.Println(json.Unmarshal(body, &requestData))
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

	queryList := strings.Split(requestData.Data, "#")

	newQuery := ""
	ttl := uint16(100)
	for i := 0; i < len(queryList); i++ {
		newQuery = queryList[0]
		ttlTemp, err := strconv.Atoi(queryList[1])
		if err != nil {
			log.Fatal("failed to convert string to integer", err)
		}
		ttl = uint16(ttlTemp)
	}
	fmt.Println("ttl =", ttl)
	fmt.Println("query =", newQuery)
	res := sendQuery(newQuery, ttl)

	// create a response containing the received data
	response := ResponseData{Message: res}

	// set response content type to JSON and send it back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// this functions creates an header with an ttl = 100 and an uuid
func sendQuery(queryString string, ttl uint16) string {

	// test if syntax is valid
	err := pathExpression.IsValid(queryString)
	if err != nil {
		return "%%" + err.Error()
	}

	// create a header of ttl = 100 and an uuid
	header := make([]byte, 0, 32)
	header = binary.BigEndian.AppendUint16(header, ttl)
	tmp := uuid.New()
	log.Printf("Created query with id %s", tmp.URN())
	// uuid marshal binary can not error, but have error return to match interface
	qid, _ := tmp.MarshalBinary()
	header = append(header, qid...)

	payload := strings.NewReader(queryString)

	stream := io.MultiReader(bytes.NewReader(header), payload)

	q, _ := pathExpression.QueryStructFromStream(&stream)

	s := pathExpression.TraverseQuery(&q)
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
