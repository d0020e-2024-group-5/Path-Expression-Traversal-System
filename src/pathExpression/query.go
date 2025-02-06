package pathExpression

import (
	"fmt"
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
	Rootpointer *RootNode // pointer to the tree
	FollowLeaf  *LeafNode // the path that should be taken to the next node
	NextNode    string    // the name of the node next node
}

// creates and returns QueryStruct from a query string
// if the string matches waht querystruct.Tostring
func bobTheBuilder(input_query string, data map[string][]DataEdge) (QueryStruct, error) {
	// pre process query, remove spaces and change
	input_query = preprocessQuery(input_query)

	id_int := 0
	root := RootNode{}
	tmp := grow_tree(input_query, &root, &id_int)
	root.Child = tmp

	// construct the tree
	q := QueryStruct{}
	q.Query = input_query
	q.Rootpointer = &root

	// TODO this need to be changed to being conditione if we have passed in the leaf node in the input_query
	// TODO these might also not be a single return, assume single for the moment
	q.FollowLeaf = root.NextNode(nil)[0].NextNode(nil)[0]

	// TODO this need to be changed to being conditione if we have passed in the next node in the input_query
	// TODO these might also not be a single return, assume single for the moment
	// TODO throw error if it cant find first node, spent way to long wondering why i got nil pointer deref
	for _, edge := range data[root.NextNode(nil)[0].Value] {
		if edge.EdgeName == q.FollowLeaf.Value {
			q.NextNode = edge.TargetName
			break
		}
	}

	return q, nil
}

// This function converts the queryStruct to an string which could be passed on to another server
func (q *QueryStruct) ToString() string {
	return fmt.Sprintf("%s\n%s\n%d", q.Query, q.NextNode, q.FollowLeaf.ID)
}

func (q *QueryStruct) DebugToString() string {
	return fmt.Sprintf("%s\nNextNode: %s\nFollowingEdge: %d (%s)", q.Query, q.NextNode, q.FollowLeaf.ID, q.FollowLeaf.Value)
}

func (q *QueryStruct) next(data map[string][]DataEdge) []QueryStruct {
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
					Rootpointer: q.Rootpointer,
					FollowLeaf:  follow_edge,
					NextNode:    exist_edge.TargetName,
				}
				nextQ = append(nextQ, copy)
			}
		}
	}
	return nextQ
}

func RecursiveTraverse(q *QueryStruct, data map[string][]DataEdge) {
	for _, qRec := range q.next(data) {
		fmt.Printf("%s\n\n", qRec.DebugToString())
		RecursiveTraverse(&qRec, data)
	}
}

func TestBob() {
	data := map[string][]DataEdge{
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

	q, _ := bobTheBuilder("s/pickaxe/{obtainedBy/hasInput}*", data)
	fmt.Printf("%s\n\n", q.DebugToString())

	RecursiveTraverse(&q, data)
}
