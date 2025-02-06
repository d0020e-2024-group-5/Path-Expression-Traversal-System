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
	name      string
	target_id int
}

// struct to represent all the info we need in the query
type QueryStruct struct {
	Query       string    // the query as a string
	Rootpointer *RootNode // pointer to the tree
	CurrentLeaf *LeafNode // the path that should be taken to the next node
	NextNode    string    // the name of the node next node
}

// creates and returns QueryStruct from a query string
// if the string matches waht querystruct.Tostring
func bobTheBuilder(input_query string, data map[int]DataNode) (QueryStruct, error) {
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
	q.CurrentLeaf = root.NextNode(nil)[0].NextNode(nil)[0]

	// TODO this need to be changed to being conditione if we have passed in the next node in the input_query
	q.NextNode = root.NextNode(nil)[0].Value
	// find the start node

	// TODO test if nextNode is -1 and return error

	return q, nil
}

// This function converts the queryStruct to an string which could be passed on to another server
func (q *QueryStruct) ToString() string {
	return fmt.Sprintf("%s\n%s\n%d", q.Query, q.NextNode, q.CurrentLeaf.ID)
}

func (q *QueryStruct) DebugToString() string {
	return fmt.Sprintf("%s\nNextNode: %s\nFollowingEdge: %d (%s)", q.Query, q.NextNode, q.CurrentLeaf.ID, q.CurrentLeaf.Value)
}

func TestBob() {
	data := map[int]DataNode{
		1: {
			"start",
			[]DataEdge{{"from_s", 2}}},
		2: {
			"end",
			[]DataEdge{},
		},
	}

	q, _ := bobTheBuilder("start/from_s", data)

	fmt.Print(q.DebugToString())
}
