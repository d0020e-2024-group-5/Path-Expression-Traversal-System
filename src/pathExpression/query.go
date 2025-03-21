package pathExpression

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"pets/dbComm"
	"pets/pathExpression/mermaidError"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// The first 6 bytes of a recursive mermaid query, a 4 byte magic "PETS"
// and an big endian u16 with the query type of 1, representing recursive mermaid
var PetsMermaidQueryHeader = [...]byte{'P', 'E', 'T', 'S', 0x00, 0x01}
var prefixList = []string{"minecraft: <http://example.org/minecraft#>"}

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
	root := RootNode{}
	tmp, err := grow_tree(components[0], &root, &id_int)
	root.Child = tmp
	query.rootPointer = &root

	if err != nil {
		log.Print(err)
	}

	if len(components) == 3 {
		i, err := strconv.Atoi(components[2])
		if err != nil {
			return err
		}
		query.followLeaf = query.rootPointer.GetLeaf(i)
		query.nextNode = components[1]

	} else if len(components) == 1 {

		query.followLeaf = root.NextNode(nil, []string{})[0]
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

// this function evolutes the query and with the help of the data and return new queries which have traversed one step
func (q *QueryStruct) next() ([]QueryStruct, []error) {
	nextQ := make([]QueryStruct, 0)
	errLst := make([]error, 0)

	// get the edges from this node, used for eval tree
	// TODO error handling
	// Why dont go have a map with collet
	//```rust
	// like rust list.map(|edge| edge.EdgeName).collet()
	//```
	list, _ := dbComm.DBGetNodeEdgesString(q.nextNode, prefixList)
	edges := make([]string, 0)
	for _, edge := range list {
		edges = append(edges, edge.EdgeName)
	}
	// for each edge we want to follow
	log.Println("Following \n", q.DebugToString(), "\n")
	for _, follow_edge := range q.followLeaf.NextNode(nil, edges) {
		log.Printf("Result edges from NextNode: %s", follow_edge.Value)

		// for each edge that exist from node
		nodeList, err := dbComm.DBGetNodeEdgesString(q.nextNode, prefixList)
		if err != nil {
			errLst = append(errLst, err)
			continue
		}

		for _, exist_edge := range nodeList {
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
	return nextQ, errLst
}

// this function takes an query struct and traverses the data.
// Returns the path the query in mermaid format
func TraverseQuery(q *QueryStruct) string {
	sBuilder := new(strings.Builder)
	RecursiveTraverse(q, sBuilder)
	return sBuilder.String()
}

// This recursively (depth first) traverses the query.
// If the path encounters a "false node" it will traverse to the server that it points to.
// TODO error that are encounter will be written as valid mermaid and have an edge to the queries nextnode
func RecursiveTraverse(q *QueryStruct, res io.Writer) {
	// if Time to live is zero write error
	if q.TimeToLive == 0 {
		mermaidError.MermaidErrorEdge(res, q.nextNode, "TTL err", fmt.Sprintf("Time to live expired\n%s(%s)", strings.ReplaceAll(q.toString(), ";", "\n"), q.followLeaf.Value))
		return
	}

	// get the queries traversed one step forward
	qNext, qErr := q.next()

	// next might multiple errors
	for _, err := range qErr {
		mermaidError.MermaidErrorEdge(res, q.nextNode, "Next err", err.Error())
	}

	// for each of the next queys
	for _, qRec := range qNext {
		fmt.Println(q.DebugToString(), "\n")

		domains, err := pointsToServer(qRec)
		if err != nil {
			mermaidError.MermaidErrorEdge(res, qRec.nextNode, "FalseNode err",
				fmt.Sprintf("Error testing for falseNode %s:\n%s", qRec.nextNode, err.Error()))
			continue
		}

		// this node does not point to another server
		if len(domains) == 0 {
			fmt.Fprintf(res, "%s-->|%s|%s\n", q.nextNode, qRec.followLeaf.Value, qRec.nextNode)
			RecursiveTraverse(&qRec, res)

		} else {
			// this does point to other servers
			fmt.Fprintf(res, "%s-.->|%s|%s\n", q.nextNode, qRec.followLeaf.Value, qRec.nextNode)

			// for each domain the false node points to
			for _, domain := range domains {

				// convert to an qRec to an stream of bytes
				stream := io.MultiReader(bytes.NewReader(PetsMermaidQueryHeader[:]), qRec.ToReader())
				log.Printf("query following querydata to %s \n%s", domain, qRec.DebugToString())

				// send that stream to the other server
				resp, err := http.Post("http://"+domain+"/api/pets", "PETSQ", stream)

				// if error sending to server
				if err != nil {
					mermaidError.MermaidErrorEdge(res, qRec.nextNode, "Network err",
						fmt.Sprintf("Could not send this this to server %s\nErr: %s\n%s", domain, err.Error(), qRec.DebugToString()))
					continue // try next domain
				}
				body := resp.Body
				defer body.Close()

				// write result from other server to my result
				io.Copy(res, body)
			}
		}
	}
}

// gets the domains a false node points to return no domains if the node points to nothing
func pointsToServer(qRec QueryStruct) ([]string, error) {
	domains := make([]string, 0)

	// test if it has en edge that indicates its a false node
	edges, err := dbComm.DBGetNodeEdgesString(qRec.nextNode, prefixList)
	if err != nil {
		return domains, err
	}

	// see if edge with weight "pointsToServer" exist
	for _, edge := range edges {
		if edge.EdgeName == "nodeOntology:pointsToServer" {
			edgesList, err := dbComm.DBGetNodeEdgesString(edge.TargetName, prefixList)
			if err != nil {
				return domains, err
			}
			// get the domain of the server
			for _, server_edge := range edgesList {
				// if the edge has contact information
				if server_edge.EdgeName == "nodeOntology:hasIP" {
					domains = append(domains, server_edge.TargetName)
				}
			}
		}
	}
	return domains, nil
}
