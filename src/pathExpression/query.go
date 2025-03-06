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
	"pets/pathExpression/mermaidError"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// The first 6 bytes of a recursive mermaid query, a 4 byte magic "PETS"
// and an big endian u16 with the query type of 1, representing recursive mermaid
var PetsMermaidQueryHeader = [...]byte{'P', 'E', 'T', 'S', 0x00, 0x01}

// preprocesses the query, removes all white spaces and converts */ to *
func preprocessQuery(inp string) string {
	//remove whitespace
	inp = strings.Replace(inp, "*/", "*", -1)
	inp = strings.Join(strings.Fields(inp), "")
	return inp
}

// Struct to represent all the info we need in the query.
// It Has the methods, ToReader, used for sending the query to other servers.
// Also has the function next, which return zero or more queries where the new queries have traversed one step.
type QueryStruct struct {
	// the query as a string, pointer to allow for faster copy, as its only read and not write, no write conflict can occur.
	// this might not be necessary as its possible to build the string from the tree
	Query *string
	// A time to live counter, if this reaches zero the query will return error and not traverse any further
	TimeToLive uint16
	// The id of the query, used mostly for logging purposes
	QueryID uuid.UUID
	// pointer to the tree
	rootPointer *RootNode
	// the path that should be taken to the next node
	followLeaf *LeafNode
	// the name of the node next node
	nextNode string
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

	// read the TTL from the stream
	err := binary.Read(*stream, binary.BigEndian, &ttl)
	if err != nil {
		return ttl, read_uuid, err
	}

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
	query.rootPointer = &RootNode{}
	tmp := grow_tree(components[0], query.rootPointer, &id_int)
	query.rootPointer.Child = tmp

	if len(components) == 3 {
		i, err := strconv.Atoi(components[2])
		if err != nil {
			return err
		}
		query.followLeaf = query.rootPointer.GetLeaf(i)
		query.nextNode = components[1]

	} else if len(components) == 1 {

		query.followLeaf = query.rootPointer.NextNode(nil)[0]
		// TODO this need to be changed to being conditione if we have passed in the next node in the input_query
		// TODO error handling, we cant be sure that the first "operator" is traverse and therefore might get a multiple return
		query.nextNode = query.followLeaf.Value
	} else {
		return errors.New("amount of 'lines' in query is not 1 or 3")
	}

	return nil
}

// Returns a reader with ttl uuid, and payload.
// type and magic are NOT included
func (q *QueryStruct) ToReader() io.Reader {
	header := make([]byte, 0, 18)
	header = binary.BigEndian.AppendUint16(header, q.TimeToLive)
	// TODO this can error
	qidBin, _ := q.QueryID.MarshalBinary()
	header = append(header, qidBin...)
	payload := q.toString()
	ret := io.MultiReader(bytes.NewReader(header), strings.NewReader(payload))
	return ret
}

// This function converts the queryStruct to an string which could be passed on to another server
func (q *QueryStruct) toString() string {
	return fmt.Sprintf("%s;%s;%d", *q.Query, q.nextNode, q.followLeaf.ID)
}

// this function returns a string which contains information useful for debugging.
// URN format of the id, TTL, query as string, nextnode as string, the follow edge id (and its value)
func (q *QueryStruct) DebugToString() string {
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "UUID: %s\n", q.QueryID.URN())
	fmt.Fprintf(&sb, "TTL: %d\n", q.TimeToLive)
	sb.WriteString(*q.Query)
	fmt.Fprintf(&sb, "\nNextNode: %s\n", q.nextNode)
	fmt.Fprintf(&sb, "FollowingEdge: %d (%s)", q.followLeaf.ID, q.followLeaf.Value)
	return sb.String()
}

// this function evolutes the query and with the help of the data and return new queries which have traversed one step.
// Note that if TTL == 0, then it returns an empty list (will be changed to return error later)
func (q *QueryStruct) next(data map[string][]parse.DataEdge) []QueryStruct {
	nextQ := make([]QueryStruct, 0)

	// for each edge we want to follow
	for _, follow_edge := range q.followLeaf.NextNode(nil) {

		// for each edge that exist from node
		for _, exist_edge := range data[q.nextNode] {

			// if it exist and we want to follow it
			if follow_edge.Value == exist_edge.EdgeName {
				// create a new query with new current leaf
				// and next node
				copy := QueryStruct{
					QueryID:     q.QueryID,
					TimeToLive:  q.TimeToLive - 1,
					Query:       q.Query,
					rootPointer: q.rootPointer,
					followLeaf:  follow_edge,
					nextNode:    exist_edge.TargetName,
				}
				nextQ = append(nextQ, copy)
			}
		}
	}
	return nextQ
}

// this function takes an query struct and traverses the data adn traverses to other servers if necessary.
// Returns the path the query in mermaid format, errors that occurred during evaluation will be converted to valid mermaid
func TraverseQuery(q *QueryStruct, data map[string][]parse.DataEdge) string {
	sBuilder := new(strings.Builder)
	RecursiveTraverse(q, data, sBuilder)
	return sBuilder.String()
}

// This recursively (depth first) traverses the query.
// If the path encounters a "false node" it will traverse to the server that it points to.
// TODO error that are encounter will be written as valid mermaid and have an edge to the queries nextnode
func RecursiveTraverse(q *QueryStruct, data map[string][]parse.DataEdge, res io.Writer) {
	// if Time to live is zero write error
	if q.TimeToLive == 0 {
		mermaidError.MermaidErrorEdge(res, q.nextNode, " ", fmt.Sprintf("Time to live expired\n%s(%s)", strings.ReplaceAll(q.toString(), ";", "\n"), q.followLeaf.Value))
		return
	}

	for _, qRec := range q.next(data) {
		// test if it has en edge that indicates its a false node
		// TODO error, there might exist a scenario when next node dont exists in our data, it should not happen but we need to be able to handle it
		edges := data[qRec.nextNode]

		// see if edge with weight "pointsToServer" exist
		// TODO break this out to own function
		for _, edge := range edges {
			if edge.EdgeName == "pointsToServer" {

				// get the domain of the server
				for _, server_edge := range data[edge.TargetName] {
					// if the edge has contact information
					if server_edge.EdgeName == "hasIP" && false {

						stream := io.MultiReader(bytes.NewReader(PetsMermaidQueryHeader[:]), qRec.ToReader())
						log.Printf("query following querydata to %s \n%s", server_edge.TargetName, qRec.DebugToString())
						resp, err := http.Post("http://"+server_edge.TargetName+"/api/pets", "PETSQ", stream)
						// TODO write error as valid mermaid
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
		fmt.Fprintf(res, "%s-->|%s|%s\n", q.nextNode, qRec.followLeaf.Value, qRec.nextNode)
		RecursiveTraverse(&qRec, data, res)
	}
}
