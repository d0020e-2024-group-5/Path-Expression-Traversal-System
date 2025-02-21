package dbComm

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type DataNode struct {
	Name string
	// Id int id is not required since this is kept track off in the map
	edges []DataEdge
}
type DataEdge struct {
	EdgeName   string
	TargetName string
}

func DBGetNodeEdgesString(str string) ([]DataEdge, error) {

	// feel free to change the hostname
	hostname := "http://localhost:7200" + "/GraphDB/repositories/Data"
	query := "SELECT ?p ?o WHERE { <" + str + "> ?p ?o } LIMIT 100"
	reqBody := []byte("query=" + query)
	req, err := http.NewRequest("POST", hostname, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, fmt.Errorf("Error creating request: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, fmt.Errorf("Error sending request: %s", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, fmt.Errorf("Error reading response: %s", err)
	}
	fmt.Println("Response:")
	fmt.Println(string(body))
	list := []DataEdge{}
	var everyOther = true
	var temp string = ""
	for _, byte := range body {
		//fmt.Println(v)
		if byte == '\n' || byte == ',' {
			if everyOther {
				list = append(list, DataEdge{EdgeName: temp, TargetName: "nil"})
				everyOther = false

			} else {
				list[len(list)-1].TargetName = temp
				everyOther = true
			}

			temp = ""

		} else {
			temp = temp + string(byte)
		}
	}
	fmt.Println("List:")
	for _, edge := range list {
		fmt.Println("edge: ", edge.EdgeName, " | ", edge.TargetName)
	}
	var ret []DataEdge = nil
	for i := 1; i < len(list); i++ {
		ret = append(ret, list[i])
	}
	return ret, nil
}

func main() {
	DBGetNodeEdgesString("http://example.org/minecraft#Stick_Bamboo_made_Instance")
}
