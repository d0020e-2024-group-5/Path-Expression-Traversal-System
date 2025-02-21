package pathExpression

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"pets/dbComm"
	"strconv"
	"strings"
)

// preprocesses the query
func preprocessQuery(inp string) string {
	//remove whitespace
	inp = strings.Replace(inp, "*/", "*", -1)
	inp = strings.Join(strings.Fields(inp), "")
	return inp
}

type DataNode struct {
	Name string
	// Id int id is not required since this is kept track off in the map
	edges []DataEdge
}
type DataEdge struct {
	EdgeName   string
	TargetName string
}

// struct to represent all the info we need in the query
type QueryStruct struct {
	Query       string    // the query as a string
	RootPointer *RootNode // pointer to the tree
	FollowLeaf  *LeafNode // the path that should be taken to the next node
	NextNode    string    // the name of the node next node
}

// creates and returns QueryStruct from a query string
// TODO when input query as multiple lines update NextNode and FollowLeaf accordingly
func BobTheBuilder(input_query string) (QueryStruct, error) {
	// pre process query, remove spaces and change
	input_query = preprocessQuery(input_query)
	// split query on separators ";"
	components := strings.Split(input_query, ";")

	id_int := 0
	root := RootNode{}
	tmp := grow_tree(components[0], &root, &id_int)
	root.Child = tmp

	// construct the tree
	q := QueryStruct{}
	q.Query = components[0]
	q.RootPointer = &root

	// TODO this need to be changed to being conditione if we have passed in the leaf node in the input_query
	if len(components) == 3 {
		i, err := strconv.Atoi(components[2])
		if err != nil {
			return q, err
		}
		q.FollowLeaf = root.GetLeaf(i)
		q.NextNode = components[1]

	} else if len(components) == 1 {

		q.FollowLeaf = root.NextNode(nil)[0]
		// TODO this need to be changed to being conditione if we have passed in the next node in the input_query
		// TODO error handling, we cant be sure that the first "operator" is traverse and therefore might get a multiple return
		q.NextNode = q.FollowLeaf.Value
	} else {
		return q, errors.New("amount of 'lines' in query is not 1 or 3")
	}

	return q, nil
}

// This function converts the queryStruct to an string which could be passed on to another server
func (q *QueryStruct) ToString() string {
	return fmt.Sprintf("%s;%s;%d", q.Query, q.NextNode, q.FollowLeaf.ID)
}

func (q *QueryStruct) DebugToString() string {
	return fmt.Sprintf("%s\nNextNode: %s\nFollowingEdge: %d (%s)", q.Query, q.NextNode, q.FollowLeaf.ID, q.FollowLeaf.Value)
}

// this function evolutes the query and with the help of the data and return new queries which have traversed one step
func (q *QueryStruct) next(data map[string][]dbComm.DataEdge) []QueryStruct {
	nextQ := make([]QueryStruct, 0)

	// for each edge we want to follow
	for _, follow_edge := range q.FollowLeaf.NextNode(nil) {

		// for each edge that exist from node
		for _, exist_edge := range data[q.NextNode] {

			// if it exist and we want to follow it
			if follow_edge.Value == exist_edge.EdgeName {
				// create a new query with new current leaf
				// and next node
				copy := QueryStruct{
					Query:       q.Query,
					RootPointer: q.RootPointer,
					FollowLeaf:  follow_edge,
					NextNode:    exist_edge.TargetName,
				}
				nextQ = append(nextQ, copy)
			}
		}
	}
	return nextQ
}

// this function takes an query struct and traverses the data.
// Returns the path the query in mermaid format
func TraverseQuery(q *QueryStruct, data map[string][]dbComm.DataEdge) string {
	sBuilder := new(strings.Builder)
	RecursiveTraverse(q, data, sBuilder)
	return sBuilder.String()
}

func RecursiveTraverse(q *QueryStruct, data map[string][]dbComm.DataEdge, res io.Writer) {
	for _, qRec := range q.next(data) {
		// test if it has en edge that indicates its a false node
		// TODO error, there might exist a scenario when next node dont exists in our data, it should not happen but we need to be able to handle it
		edges := data[qRec.NextNode]

		// see if edge with weight "pointsToServer" exist
		for _, edge := range edges {
			if edge.EdgeName == "pointsToServer" {

				// get the domain of the server
				for _, server_edge := range data[edge.TargetName] {
					// if the edge has contact information
					if server_edge.EdgeName == "hasIP" {
						q_string := qRec.ToString()

						// TODO error handling
						resp, err := http.Post("http://"+server_edge.TargetName+"/api/recq", "PETSQ", strings.NewReader(q_string))
						if err != nil {
							log.Fatalf("error on passing to server: %s", err.Error())
						}
						body := resp.Body
						defer body.Close()
						io.Copy(res, body)
					}
				}
			}
		}
		// TODO, change arrow type if its a false node
		fmt.Fprintf(res, "%s-->|%s|%s\n", q.NextNode, qRec.FollowLeaf.Value, qRec.NextNode)
		RecursiveTraverse(&qRec, data, res)

	}
}

func TestBob() {
	data := map[string][]dbComm.DataEdge{
		"s": {
			{"pickaxe", "pickaxe"},
		},
		"pickaxe": {
			{"obtainedBy", "Pickaxe_From_Stick_And_Stone_Recipe"},
		},
		"Pickaxe_From_Stick_And_Stone_Recipe": {
			{"hasInput", "Stick"},
			{"hasInput", "Cobblestone"},
		},
	}
	fmt.Printf("%v\n\n", data)

	q, _ := BobTheBuilder("s/pickaxe/{obtainedBy/hasInput}*")
	fmt.Printf("%s\n\n", q.DebugToString())

	fmt.Println(TraverseQuery(&q, data))
}

func TestBob2() {
	data := map[string][]dbComm.DataEdge{
		"s": {
			{"pickaxe", "pickaxe"},
		},
		"pickaxe": {
			{"obtainedBy", "Pickaxe_From_Stick_And_Stone_Recipe"},
		},
		"Pickaxe_From_Stick_And_Stone_Recipe": {
			{"hasInput", "Stick"},
			{"hasInput", "Cobblestone"},
		},
	}

	q, _ := BobTheBuilder("s/pickaxe/{obtainedBy/hasInput}*")
	// fmt.Printf("%s\n\n", q.DebugToString())

	q = q.next(data)[0]
	fmt.Printf("%s\n\n", q.DebugToString())

	q = q.next(data)[0]
	fmt.Printf("%s\n\n", q.DebugToString())

	q = q.next(data)[0]
	fmt.Printf("%s\n\n", q.DebugToString())

	fmt.Println("=================================")

	q, _ = BobTheBuilder("s/pickaxe/{obtainedBy/hasInput}*")
	// fmt.Printf("%s\n\n", q.DebugToString())

	q, _ = BobTheBuilder(q.next(data)[0].ToString())
	fmt.Printf("%s\n\n", q.DebugToString())

	q, _ = BobTheBuilder(q.next(data)[0].ToString())
	fmt.Printf("%s\n\n", q.DebugToString())

	q, _ = BobTheBuilder(q.next(data)[0].ToString())
	fmt.Printf("%s\n\n", q.DebugToString())

}
