package pathExpression

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"pets/parse"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// preprocesses the query
func preprocessQuery(inp string) string {
	//remove whitespace
	inp = strings.Replace(inp, "*/", "*", -1)
	inp = strings.Join(strings.Fields(inp), "")
	return inp
}

// struct to represent all the info we need in the query
type QueryStruct struct {
	// the query as a string, pointer to allow for faster copy, as its only read and not write, no write conflict can occur.
	// this might not be necessary as its possible to build the string from the tree
	Query *string
	// A time to live counter, if this reaches zero the query will return error and not traverse any further
	TimeToLive uint16
	// The id of the query, used mostly for logging purposes
	QueryID uuid.UUID
	// pointer to the tree
	RootPointer *RootNode
	// the path that should be taken to the next node
	FollowLeaf *LeafNode
	// the name of the node next node
	NextNode string
}

// This function reads the required data from the stream to create a Query struct object
// The structure of the API is ass follows
// "PETS" // 32bits (4bytes) used as magic to confirm this is a PETS query
// Type of query 16 bits (2bytes) (network order)
// TTL 16 bits (2 bytes) (network order)
// QueryID 128 bits (16bytes), an uuid
// NOTE that the 6 first bytes needs to have been read for this to work
func QueryStructFromStream(stream *io.Reader) (QueryStruct, error) {

	var q QueryStruct

	// read the ttl and uuid
	ttl, id, err := getTTLandUUID(stream)
	q.TimeToLive = ttl
	q.QueryID = id
	if err != nil {
		return q, err
	}
	// this is cursed syntax
	return q, readPayloadToQuery(stream, &q)

}

// this functions reads 2 bytes as network order uint16 as the ttl,
// it then reads 16 bytes of data and converts that to an valid uuid
func getTTLandUUID(stream *io.Reader) (uint16, uuid.UUID, error) {
	var read_uuid uuid.UUID
	var ttl uint16

	var rawuint16 [2]byte
	_, err := io.ReadFull(*stream, rawuint16[:])
	if err != nil {
		return ttl, read_uuid, err
	}
	ttl = binary.BigEndian.Uint16(rawuint16[:])

	// read the query id
	var raw_uuid [16]byte
	_, err = io.ReadFull(*stream, raw_uuid[:])
	if err != nil {
		return ttl, read_uuid, err
	}

	// convert to uuid
	read_uuid, err = uuid.FromBytes(raw_uuid[:])
	if err != nil {
		return ttl, read_uuid, err
	}

	return ttl, read_uuid, nil
}

// read the payload and update the query structs fields
func readPayloadToQuery(stream *io.Reader, query *QueryStruct) error {

	// read the payload
	full, err := io.ReadAll(*stream)
	if err != nil {
		return err
	}

	// convert payload to a string
	stringPayloadBuilder := strings.Builder{}
	_, err = stringPayloadBuilder.Write(full)
	if err != nil {
		return err
	}

	// pre process query, remove spaces and change
	input_query := preprocessQuery(stringPayloadBuilder.String())

	// split query on separators ";"
	components := strings.Split(input_query, ";")

	// add the query
	query.Query = &components[0]

	// create the evaluation tree
	id_int := 0
	query.RootPointer = &RootNode{}
	tmp := grow_tree(components[0], query.RootPointer, &id_int)
	query.RootPointer.Child = tmp

	if len(components) == 3 {
		i, err := strconv.Atoi(components[2])
		if err != nil {
			return err
		}
		query.FollowLeaf = query.RootPointer.GetLeaf(i)
		query.NextNode = components[1]

	} else if len(components) == 1 {

		query.FollowLeaf = query.RootPointer.NextNode(nil)[0]
		// TODO this need to be changed to being conditione if we have passed in the next node in the input_query
		// TODO error handling, we cant be sure that the first "operator" is traverse and therefore might get a multiple return
		query.NextNode = query.FollowLeaf.Value
	} else {
		return errors.New("amount of 'lines' in query is not 1 or 3")
	}

	return nil
}

// This function converts the queryStruct to an string which could be passed on to another server
func (q *QueryStruct) ToString() string {
	return fmt.Sprintf("%s;%s;%d", *q.Query, q.NextNode, q.FollowLeaf.ID)
}

func (q *QueryStruct) DebugToString() string {
	return fmt.Sprintf("%s\nNextNode: %s\nFollowingEdge: %d (%s)", *q.Query, q.NextNode, q.FollowLeaf.ID, q.FollowLeaf.Value)
}

// this function evolutes the query and with the help of the data and return new queries which have traversed one step
func (q *QueryStruct) next(data map[string][]parse.DataEdge) []QueryStruct {
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
func TraverseQuery(q *QueryStruct, data map[string][]parse.DataEdge) string {
	sBuilder := new(strings.Builder)
	RecursiveTraverse(q, data, sBuilder)
	return sBuilder.String()
}

func RecursiveTraverse(q *QueryStruct, data map[string][]parse.DataEdge, res io.Writer) {
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
	data := map[string][]parse.DataEdge{
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

	// create a header of ttl = 100 and an uuid
	header := make([]byte, 0, 32)
	binary.BigEndian.AppendUint16(header, 100)
	qid, _ := uuid.New().MarshalBinary()
	header = append(header, qid...)

	payload := strings.NewReader("s/pickaxe/{obtainedBy/hasInput}*")

	stream := io.MultiReader(bytes.NewReader(header), payload)

	q, _ := QueryStructFromStream(&stream)
	fmt.Printf("%s\n\n", q.DebugToString())

	fmt.Println(TraverseQuery(&q, data))
}

// func TestBob2() {
// 	data := map[string][]DataEdge{
// 		"s": {
// 			{"pickaxe", "pickaxe"},
// 		},
// 		"pickaxe": {
// 			{"obtainedBy", "Pickaxe_From_Stick_And_Stone_Recipe"},
// 		},
// 		"Pickaxe_From_Stick_And_Stone_Recipe": {
// 			{"hasInput", "Stick"},
// 			{"hasInput", "Cobblestone"},
// 		},
// 	}

// 	q, _ := BobTheBuilder("s/pickaxe/{obtainedBy/hasInput}*")
// 	// fmt.Printf("%s\n\n", q.DebugToString())

// 	q = q.next(data)[0]
// 	fmt.Printf("%s\n\n", q.DebugToString())

// 	q = q.next(data)[0]
// 	fmt.Printf("%s\n\n", q.DebugToString())

// 	q = q.next(data)[0]
// 	fmt.Printf("%s\n\n", q.DebugToString())

// 	fmt.Println("=================================")

// 	q, _ = BobTheBuilder("s/pickaxe/{obtainedBy/hasInput}*")
// 	// fmt.Printf("%s\n\n", q.DebugToString())

// 	q, _ = BobTheBuilder(q.next(data)[0].ToString())
// 	fmt.Printf("%s\n\n", q.DebugToString())

// 	q, _ = BobTheBuilder(q.next(data)[0].ToString())
// 	fmt.Printf("%s\n\n", q.DebugToString())

// 	q, _ = BobTheBuilder(q.next(data)[0].ToString())
// 	fmt.Printf("%s\n\n", q.DebugToString())

// }
