package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//	type Node struct {
//		NodeName string
//		NodeID   int
//	}
//
// type Tuple struct { // FOR EDGES
//
//		key   string // ex. obtainedBy
//		value string // ex. recipe
//	}
func parse(nodeLst map[Node][]Tuple) map[Node][]Tuple {

	file, err := os.Open("./shared_volume/data.ttl") // READ DATA FILE
	if err != nil {
		fmt.Println(err)
		return nodeLst
	}

	defer file.Close()

	var newNode Node
	var tempTuple Tuple
	scanner := bufio.NewScanner(file) // READ ONTOLOGY LINE BY LINE
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "@prefix") { // SKIP LINES WITH "@PREFIX"
			continue
		}
		if strings.HasPrefix(line, "minecraft:") {
			temp := strings.TrimPrefix(line, "minecraft:")

			wrd := getWrd(temp) // FIRST WORD IN LINE

			newNode.NodeName = wrd // ASSIGN TO NODENAME
			fmt.Printf("%s", wrd+"\n\n")

		} else if strings.HasPrefix(line, "	nodeOntology:hasID ") { // CHECK ID
			temp := strings.TrimPrefix(line, "	nodeOntology:hasID ")

			wrd := getWrd(temp) // FIRST WORD IN LINE

			i, err := strconv.Atoi(wrd) // STRING TO INT
			if err != nil {
				fmt.Println("Error reading file:", err)
			}
			newNode.NodeID = i                                     // ASSIGN TO NODEID
			nodeLst[newNode] = append(nodeLst[newNode], tempTuple) // APPEND KEY NODE TO MAP OF NODES
			fmt.Println(i)
			// CHECK FOR EDGES IN FOLLOWING ELSE IF STATEMENT
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
			nodeLst[newNode] = append(nodeLst[newNode], tempTuple) // APPEND THE EDGES TO TUPLE SLICE
		}
		if strings.HasSuffix(line, ";") {
			continue // NEXT LINE IN SAME NODE
		} else if strings.HasSuffix(line, ".") { // END OF NODE
			// APPEND TO LIST OF NODES
			tempTuple.key = ""
			tempTuple.value = ""
		} else {
			continue // NEWLINE/EMPTY SPACE
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	fmt.Println(nodeLst)
	return nodeLst
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
