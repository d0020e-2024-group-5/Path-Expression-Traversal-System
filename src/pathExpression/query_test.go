package pathExpression

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	p "pets/parse"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// This functions traverse the data and sees if the resulting mermaid matches the expected outcome
// TODO, make the companions ignore the order of the mermaid
func TestTraversal(t *testing.T) {
	data := map[string][]p.DataEdge{
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
	header = binary.BigEndian.AppendUint16(header, 100)
	qid, err := uuid.New().MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	header = append(header, qid...)

	payload := strings.NewReader("s/pickaxe/{obtainedBy/hasInput}*")

	stream := io.MultiReader(bytes.NewReader(header), payload)

	q, err := QueryStructFromStream(&stream)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s\n\n", q.DebugToString())

	result := TraverseQuery(&q, data)

	expected := "s-->|pickaxe|pickaxe\npickaxe-->|obtainedBy|Pickaxe_From_Stick_And_Stone_Recipe\nPickaxe_From_Stick_And_Stone_Recipe-->|hasInput|Stick\nPickaxe_From_Stick_And_Stone_Recipe-->|hasInput|Cobblestone\n"
	t.Log(result)
	if expected != result {
		t.Fatal("Result did not match expected value")
	}
}

// This function test if traversal without converting from reader and back
// outputs the same path as acutely doing the conversion.
// The propose of this is to se if sending to other server is viable as
// to send to another server requires it to be converted to a format that could be sent over the network
func TestTraversalWithToReader(t *testing.T) {

	// test data
	data := map[string][]p.DataEdge{
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

	// make a query
	header := make([]byte, 0, 32)
	header = binary.BigEndian.AppendUint16(header, 100)
	qid, _ := uuid.New().MarshalBinary()
	header = append(header, qid...)

	payload := strings.NewReader("s/pickaxe/{obtainedBy/hasInput}*")

	stream := io.MultiReader(bytes.NewReader(header), payload)

	qStart, err := QueryStructFromStream(&stream)
	if err != nil {
		t.Fatal(err)
	}

	// store the steps taken
	internal_steps := make([]string, 0, 10)

	// add starting point
	internal_steps = append(internal_steps, qStart.nextNode)

	// walk once and add step
	q := qStart.next(data)[0]
	internal_steps = append(internal_steps, q.nextNode)

	// walk once and add step
	q = q.next(data)[0]
	internal_steps = append(internal_steps, q.nextNode)

	// walk once and add step
	q = q.next(data)[0]
	internal_steps = append(internal_steps, q.nextNode)

	// ================================
	// Repeat but convert back and forth

	server_steps := make([]string, 0, 10)

	// add starting point
	server_steps = append(server_steps, qStart.nextNode)

	// walk once and add step
	r := qStart.ToReader()
	q, err = QueryStructFromStream(&r)
	if err != nil {
		t.Fatal(err)
	}
	q = q.next(data)[0]
	server_steps = append(server_steps, q.nextNode)

	// walk once and add step
	r = q.ToReader()
	q, err = QueryStructFromStream(&r)
	if err != nil {
		t.Fatal(err)
	}
	q = q.next(data)[0]
	server_steps = append(server_steps, q.nextNode)

	// walk once and add step
	r = q.ToReader()
	q, err = QueryStructFromStream(&r)
	if err != nil {
		t.Fatal(err)
	}
	q = q.next(data)[0]
	server_steps = append(server_steps, q.nextNode)

	// ================================
	// compare result

	for i := 0; i < 4; i++ {
		if internal_steps[i] != server_steps[i] {
			t.Fatalf("step %d, they are not equal %s != %s", i, internal_steps[i], server_steps[i])
		}

	}

}
