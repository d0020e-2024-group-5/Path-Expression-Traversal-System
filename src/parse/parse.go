package parse

import (
	"bufio"
	"fmt"
	"log"
	"os"
	p "pets/pathExpression"
	"strings"
)

func Parse() map[string][]p.DataEdge { // FUNCTION READS DATA FILE LINE BY LINE, CHECKING PREFIX KEYWORDS BEFORE APPENDING RELEVANT NODES AND EDGES TO THE NODE HASHMAP "nodeLst"

	file, err := os.Open("./shared_volume/data.ttl") // READ DATA FILE
	if err != nil {
		log.Printf("can't open data.ttl, fallback to Example Data_c.ttl: %s", err.Error())
		file, err = os.Open("./../Example Data/Server C/Example Data_C.ttl")
		if err != nil {
			log.Fatalf("cant open fallback data: %s", err.Error())
		}
	}

	defer file.Close()
	nodeLst := make(map[string][]p.DataEdge)
	var tempTuple p.DataEdge
	var firstWord string
	scanner := bufio.NewScanner(file) // READ ONTOLOGY LINE BY LINE
	for scanner.Scan() {

		line := scanner.Text()
		if strings.HasPrefix(line, "@prefix") { // SKIP LINES WITH "@PREFIX"
			continue
		} // minecraft:Server_a a nodeOntology:Server ;

		if strings.HasPrefix(line, "minecraft:") {
			temp := strings.TrimPrefix(line, "minecraft:")

			wrd := getWrd(temp) // FIRST WORD IN LINE
			firstWord = wrd

			//newNode.NodeName = wrd // ASSIGN TO NODENAME

			// } else if strings.HasPrefix(line, "	nodeOntology:hasID ") { // CHECK ID

			temp = strings.TrimPrefix(line, "	nodeOntology:hasID ")
			_, ok := nodeLst[firstWord]
			if !ok {
				nodeLst[firstWord] = make([]p.DataEdge, 0)
			}
			wrd = getWrd(temp) // FIRST WORD IN LINE

			if entry, ok := nodeLst[wrd]; ok {
				entry = append(entry, tempTuple)
				nodeLst[wrd] = entry
			} // APPEND THE EDGES TO TUPLE SLICE // APPEND KEY NODE TO MAP OF NODES
			// CHECK FOR EDGES IN FOLLOWING ELSE IF STATEMENT
		} else if strings.HasPrefix(line, "	minecraft:obtainedBy") || (strings.HasPrefix(line, "	minecraft:hasInput")) || (strings.HasPrefix(line, "	minecraft:hasOutput") || (strings.HasPrefix(line, "	minecraft:usedInStation")) || (strings.HasPrefix(line, "	nodeOntology:pointsToServer")) || (strings.HasPrefix(line, "	nodeOntology:hasIP"))) {

			if strings.HasPrefix(line, "	minecraft:obtainedBy") {
				temp := strings.TrimPrefix(line, "	minecraft:obtainedBy minecraft:")
				wrd := getWrd(temp)
				tempTuple.EdgeName = "obtainedBy"
				tempTuple.TargetName = wrd
			} else if strings.HasPrefix(line, "	minecraft:hasInput") {
				temp := strings.TrimPrefix(line, "	minecraft:hasInput minecraft:")
				wrd := getWrd(temp)
				tempTuple.EdgeName = "hasInput"
				tempTuple.TargetName = wrd
			} else if strings.HasPrefix(line, "	minecraft:hasOutput") {
				temp := strings.TrimPrefix(line, "	minecraft:hasOutput minecraft:")
				wrd := getWrd(temp)
				tempTuple.EdgeName = "hasOutput"
				tempTuple.TargetName = wrd
			} else if strings.HasPrefix(line, "	minecraft:usedInStation") {
				temp := strings.TrimPrefix(line, "	minecraft:usedInStation minecraft:")
				wrd := getWrd(temp)
				tempTuple.EdgeName = "usedInStation"
				tempTuple.TargetName = wrd
			} else if strings.HasPrefix(line, "	nodeOntology:pointsToServer") {
				temp := strings.TrimPrefix(line, "	nodeOntology:pointsToServer minecraft:")
				wrd := getWrd(temp)
				tempTuple.EdgeName = "pointsToServer"
				tempTuple.TargetName = wrd
			} else if strings.HasPrefix(line, "	nodeOntology:hasIP") {
				temp := strings.TrimPrefix(line, "	nodeOntology:hasIP ")
				wrd := getWrd(temp)
				tempTuple.EdgeName = "hasIP"
				tempTuple.TargetName = wrd
			}

			if entry, ok := nodeLst[firstWord]; ok {
				entry = append(entry, tempTuple)
				nodeLst[firstWord] = entry
			} // APPEND THE EDGES TO TUPLE SLICE
		}
		if strings.HasSuffix(line, ";") {
			continue // NEXT LINE IN SAME NODE
		} else if strings.HasSuffix(line, ".") { // END OF NODE
			// APPEND TO LIST OF NODES
			firstWord = "" // EMPTY NODE (NEW NODE)
		} else {
			continue // NEWLINE/EMPTY SPACE
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
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
