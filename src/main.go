package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var hName string

type Node struct {
	NodeName string
	NodeID   int
}
type Tuple struct { // FOR EDGES
	key   string // ex. obtainedBy
	value string // ex. recipe
}

var nodeLst = map[Node][]Tuple{} // NODE HASHMAP WITH A TUPLE LIST (EDGES) AS VALUE

func main() {
	parse()
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

	// create the server and listen to port 80
	http.ListenAndServe(":80", nil)
}

func parse() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path) // for example /home/user
	file, err := os.Open("./shared_volume/data.ttl")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	var lines []string
	var newNode Node
	var tempTuple Tuple
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if strings.HasPrefix(line, "@prefix") {
			continue
		}
		if strings.HasPrefix(line, "minecraft:") {
			temp := strings.TrimPrefix(line, "minecraft:")

			wrd := getWrd(temp) // FIRST WORD

			newNode.NodeName = wrd // ASSIGN TO NODENAME
			fmt.Printf(wrd + "\n\n")

		} else if strings.HasPrefix(line, "	nodeOntology:hasID ") {
			temp := strings.TrimPrefix(line, "	nodeOntology:hasID ")

			wrd := getWrd(temp) // FIRST WORD

			i, err := strconv.Atoi(wrd) // STRING TO INT
			if err != nil {
				fmt.Println("Error reading file:", err)
			}
			newNode.NodeID = i // ASSIGN TO NODEID
			nodeLst[newNode] = append(nodeLst[newNode], tempTuple)
			//fmt.Println(nodeLst)
			fmt.Println(i)
		} else if strings.HasPrefix(line, "    minecraft:obtainedBy") || (strings.HasPrefix(line, "    minecraft:hasInput")) || (strings.HasPrefix(line, "    minecraft:hasOutput") || (strings.HasPrefix(line, "    minecraft:usedInStation"))) {

			if strings.HasPrefix(line, "    minecraft:obtainedBy") {
				temp := strings.TrimPrefix(line, "    minecraft:obtainedBy minecraft:")
				wrd := getWrd(temp)
				tempTuple.key = "obtainedBy"
				tempTuple.value = wrd
			} else if strings.HasPrefix(line, "    minecraft:hasInput") {
				temp := strings.TrimPrefix(line, "    minecraft:hasInput minecraft:")
				wrd := getWrd(temp)
				tempTuple.key = "hasInput"
				tempTuple.value = wrd
			} else if strings.HasPrefix(line, "    minecraft:hasOutput") {
				temp := strings.TrimPrefix(line, "    minecraft:hasOutput minecraft:")
				wrd := getWrd(temp)
				tempTuple.key = "hasOutput"
				tempTuple.value = wrd
			} else if strings.HasPrefix(line, "    minecraft:usedInStation") {
				temp := strings.TrimPrefix(line, "    minecraft:usedInStation minecraft:")
				wrd := getWrd(temp)
				tempTuple.key = "usedInStation"
				tempTuple.value = wrd
			}
			nodeLst[newNode] = append(nodeLst[newNode], tempTuple)

			//fmt.Println(nodeLst)
		}
		if strings.HasSuffix(line, ";") {
			continue
		} else if strings.HasSuffix(line, ".") { // END OF NODE
			//var tempTuple Tuple

			// APPEND TO LIST OF NODES
			tempTuple.key = ""
			tempTuple.value = ""
		} else {
			continue // NEWLINE/EMPTY SPACE
		}
		//fmt.Println(line)

	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	fmt.Println(nodeLst)
	fmt.Printf("bomba")
}
func getWrd(w string) string {
	wrd := ""
	for i := range w {
		if w[i] == ' ' {
			wrd = w[0:i]
			break
		}
	}
	return wrd
}
