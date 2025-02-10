package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)


func parse(nodeLst map[string]DataNode) map[string]DataNode { // FUNCTION READS DATA FILE LINE BY LINE, CHECKING PREFIX KEYWORDS BEFORE APPENDING RELEVANT NODES AND EDGES TO THE NODE HASHMAP "nodeLst"

	file, err := os.Open("./shared_volume/data.ttl") // READ DATA FILE
	if err != nil {
		fmt.Println(err)
		return nodeLst
	}

	defer file.Close()
	nodeLst = make(map[string]DataNode)
	var tempN DataNode
	var tempTuple DataEdge
	var firstWord string
	scanner := bufio.NewScanner(file) // READ ONTOLOGY LINE BY LINE
	for scanner.Scan() {
		
		line := scanner.Text()
		if strings.HasPrefix(line, "@prefix") { // SKIP LINES WITH "@PREFIX"
			continue
		}
		if strings.HasPrefix(line, "minecraft:") {
			temp := strings.TrimPrefix(line, "minecraft:")

			wrd := getWrd(temp) // FIRST WORD IN LINE
			firstWord = wrd			

		} else if strings.HasPrefix(line, "	nodeOntology:hasID ") { // CHECK ID
			temp := strings.TrimPrefix(line, "	nodeOntology:hasID ")
			nodeLst[firstWord] = tempN
			wrd := getWrd(temp) // FIRST WORD IN LINE
			
			i, err := strconv.Atoi(wrd) // STRING TO INT
			if err != nil {
				fmt.Println("Error reading file:", err)
			}
			
			if entry, ok := nodeLst[wrd]; ok {
				entry.Edges = append(entry.Edges, tempTuple)
				nodeLst[wrd] = entry
			} // APPEND THE EDGES TO TUPLE SLICE // APPEND KEY NODE TO MAP OF NODES
			fmt.Println(i)
			// CHECK FOR EDGES IN FOLLOWING ELSE IF STATEMENT
		} else if strings.HasPrefix(line, "    minecraft:obtainedBy") || (strings.HasPrefix(line, "    minecraft:hasInput")) || (strings.HasPrefix(line, "    minecraft:hasOutput") || (strings.HasPrefix(line, "    minecraft:usedInStation")) || (strings.HasPrefix(line, "	nodeOntology:pointsToServer"))) {

			if strings.HasPrefix(line, "    minecraft:obtainedBy") {
				temp := strings.TrimPrefix(line, "    minecraft:obtainedBy minecraft:")
				wrd := getWrd(temp)
				tempTuple.EdgeName = "obtainedBy"
				tempTuple.TargetName = wrd
			} else if strings.HasPrefix(line, "    minecraft:hasInput") {
				temp := strings.TrimPrefix(line, "    minecraft:hasInput minecraft:")
				wrd := getWrd(temp)
				tempTuple.EdgeName = "hasInput"
				tempTuple.TargetName = wrd
			} else if strings.HasPrefix(line, "    minecraft:hasOutput") {
				temp := strings.TrimPrefix(line, "    minecraft:hasOutput minecraft:")
				wrd := getWrd(temp)
				tempTuple.EdgeName = "hasOutput"
				tempTuple.TargetName = wrd
			} else if strings.HasPrefix(line, "    minecraft:usedInStation") {
				temp := strings.TrimPrefix(line, "    minecraft:usedInStation minecraft:")
				wrd := getWrd(temp)
				tempTuple.EdgeName = "usedInStation"
				tempTuple.TargetName = wrd
			} else if strings.HasPrefix(line, "	nodeOntology:pointsToServer"){
				temp := strings.TrimPrefix(line, "	nodeOntology:pointsToServer minecraft:")
				wrd := getWrd(temp)
				tempTuple.EdgeName = "pointsToServer"
				tempTuple.TargetName = wrd	
			}
			
			if entry, ok := nodeLst[firstWord]; ok {
				entry.Edges = append(entry.Edges, tempTuple)
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
	fmt.Println(nodeLst)
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
