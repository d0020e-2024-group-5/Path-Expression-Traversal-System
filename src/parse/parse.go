package parse

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

func GetNodeEdgesID(node string) []DataEdge {
	var nodeLst = Parse()
	for _, edges := range nodeLst {
		for _, edge := range edges {
			if edge.EdgeName == "hasID" && edge.TargetName == node {
				return edges
			}
		}
	}
	fmt.Println("Node not found")
	return nil
}

func GetNodeEdgesString(node string) []DataEdge {
	return Parse()[node]
}

func Parse() map[string][]DataEdge { // FUNCTION READS DATA FILE LINE BY LINE, CHECKING PREFIX KEYWORDS BEFORE APPENDING RELEVANT NODES AND EDGES TO THE NODE HASHMAP "nodeLst"

	file, err := os.Open("./shared_volume/data.ttl") // READ DATA FILE
	if err != nil {
		log.Printf("can't open data.ttl, fallback to Example Data_C.ttl: %s", err.Error())
		file, err = os.Open("./../Example Data/Server C/Example Data_C.ttl")
		if err != nil {
			log.Fatalf("cant open fallback data: %s", err.Error())
		}
	}

	defer file.Close()
	nodeLst := make(map[string][]DataEdge)
	var tempTuple DataEdge
	var firstWord string
	var first bool = true
	scanner := bufio.NewScanner(file) // READ ONTOLOGY LINE BY LINE
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "@prefix") || strings.HasPrefix(line, "#") {
			first = true
			continue
		}
		if len(strings.TrimSpace(line)) == 0 { // checks empty lines
			continue
		}
		if first { // initialize node
			//fmt.Println(line)
			i := strings.Index(line, ":") + 1
			temp := line[i:]

			wrd := getWrd(temp) // FIRST WORD IN LINE
			firstWord = wrd

			_, ok := nodeLst[firstWord]
			if !ok {
				nodeLst[firstWord] = make([]DataEdge, 0)
			}
			//wrd = getWrd(temp) // FIRST WORD IN LINE
			first = false
			//fmt.Println(firstWord)
		} else { // set node attributes/edges
			i := strings.Index(line, ":") + 1
			temp := line[i:]
			wrd := getWrd(temp)
			tempTuple.EdgeName = wrd
			//fmt.Println(wrd)
			if !strings.Contains(temp, ":") {
				temp = strings.TrimPrefix(temp, (wrd + " "))
				wrd = getWrd(temp)
				tempTuple.TargetName = wrd
				//fmt.Println("id: " + wrd)
			} else {
				i = strings.Index(temp, ":") + 1
				temp = temp[i:]
				wrd = getWrd(temp)
				tempTuple.TargetName = wrd
				//fmt.Println("not id: " + wrd)
			}

			if entry, ok := nodeLst[firstWord]; ok {
				entry = append(entry, tempTuple)
				nodeLst[firstWord] = entry
				//fmt.Println("NODE LIST: ")
			} // APPEND THE EDGES TO TUPLE SLICE
		}
		if strings.HasSuffix(line, ";") {
			continue // NEXT LINE IN SAME NODE
		} else if strings.HasSuffix(line, ".") { // END OF NODE
			// APPEND TO LIST OF NODES
			first = true
			firstWord = "" // EMPTY NODE (NEW NODE)
		} else {
			continue // NEWLINE/EMPTY SPACE
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	//for k, v := range nodeLst {
	//fmt.Printf("%s: %v\n", k, v)
	//fmt.Println("-----------------")
	//}
	return nodeLst
}
func getWrd(w string) string { // GETS THE FIRST WORD SEPARETED BY A SPACE
	wrd := ""
	for i := range w {
		if w[i] == ' ' {
			wrd = w[0:i]
			break
		}
	}
	return wrd
}
