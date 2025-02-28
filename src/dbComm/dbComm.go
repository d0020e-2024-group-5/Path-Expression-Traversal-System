package dbComm

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
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

var cashMap = map[string][]DataEdge{}

// DBGetNodeEdgesString returns the edges of a node in the form of a list of DataEdge
// The node is specified by the string node
// The prefix is a list of strings that are used as prefixes in the query (can be empty)
func DBGetNodeEdgesString(node string, prefix []string) ([]DataEdge, error) {
	hashing := os.Getenv("CASHING")
	if hashing == "TRUE" {
		i, ok := cashMap[node]
		if ok {
			return i, nil
		}
	}
	// TODO hårdkådad
	prefixMap := map[string]string{
		"rdf:":          "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
		"rdfs:":         "http://www.w3.org/2000/01/rdf-schema#",
		"owl:":          "http://www.w3.org/2002/07/owl#",
		"nodeOntology:": "http://example.org/NodeOntology#",
		"minecraft:":    "http://example.org/minecraft#",
	}

	//sql injection protection
	re := regexp.MustCompile(`\s|{|}|\s\.|\n|,|;`)
	//if re.MatchString(node) {
	//	return nil, fmt.Errorf("sql injection protection node. ")
	//}

	//predetermined prefix
	prefixStr := "PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>\nPREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>\nPREFIX owl: <http://www.w3.org/2002/07/owl#>\nPREFIX nodeOntology: <http://example.org/NodeOntology#>\n"
	//loading the prefix
	// sql injection might be possible here
	for _, str := range prefix {
		//sql injection protection
		if !strings.ContainsRune(str, rune('ö')) {
			str = strings.Replace(str, ": <", "ö", 1)
			if re.MatchString(str) {
				return nil, fmt.Errorf("sql injection protection Prefix1. ")
			}
			str = strings.Replace(str, "ö", ": <", 1)
		} else {
			return nil, fmt.Errorf("sql injection protection Prefix2. ")
		}
		prefixStr += "PREFIX " + str + "\n"
	}
	fmt.Println(" Hoasname: ", os.Getenv("GRAPHDB_HOSTNAME"), " Repository: ", os.Getenv("GRAPHDB_REPOSITORY"))
	hostname := "http://" + os.Getenv("GRAPHDB_HOSTNAME") + ":7200" + "/repositories/" + os.Getenv("GRAPHDB_REPOSITORY")
	// sql injection might be possible here
	query := prefixStr + "SELECT ?p ?o WHERE { " + node + " ?p ?o } limit 100"
	fmt.Println(query)
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

	// parse the response
	for _, byte := range body {
		//fmt.Println(v)
		if byte == '\n' || byte == ',' {
			if everyOther {
				if strings.TrimSpace(replacePrefix(temp, prefixMap)) != replacePrefix(temp, prefixMap) {
					fmt.Println("inte lika: ", strings.TrimSpace(replacePrefix(temp, prefixMap)), "|", replacePrefix(temp, prefixMap))
				}
				list = append(list, DataEdge{EdgeName: strings.TrimSpace(replacePrefix(temp, prefixMap)), TargetName: "nil"})
				everyOther = false

			} else {
				if strings.TrimSpace(replacePrefix(temp, prefixMap)) != replacePrefix(temp, prefixMap) {
					fmt.Println("inte lika: ", strings.TrimSpace(replacePrefix(temp, prefixMap)), "|", replacePrefix(temp, prefixMap))
				}
				list[len(list)-1].TargetName = strings.TrimSpace(replacePrefix(temp, prefixMap))
				everyOther = true
			}

			temp = ""

		} else {
			temp = temp + string(byte)
		}
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("error: query wrong / wrong address1 ")
	}
	if list[0].EdgeName != "p" {
		return nil, fmt.Errorf("error: query wrong / wrong address2 ")
	}
	// remove the first element
	var ret []DataEdge = nil
	for i := 1; i < len(list); i++ {
		ret = append(ret, list[i])
	}
	if hashing == "TRUE" {
		cashMap[node] = ret
	}
	return ret, nil
}

func main() {
	DBGetNodeEdgesString("http://example.org/minecraft#Stick_Bamboo_made_Instance", nil)
}

func replacePrefix(in string, prefixMap map[string]string) string {
	for key, value := range prefixMap {
		if strings.HasPrefix(in, key) {
			return strings.Replace(in, key, value, 1)
		}
		if strings.HasPrefix(in, value) {
			return strings.Replace(in, value, key, 1)
		}
	}
	return in
}
