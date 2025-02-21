package pathExpression

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestBob(t *testing.T) {
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

func TestTraversalSingleStep(t *testing.T) {

	// test data
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
	internal_steps = append(internal_steps, qStart.NextNode)

	// walk once and add step
	q := qStart.next(data)[0]
	internal_steps = append(internal_steps, q.NextNode)

	// walk once and add step
	q = q.next(data)[0]
	internal_steps = append(internal_steps, q.NextNode)

	// walk once and add step
	q = q.next(data)[0]
	internal_steps = append(internal_steps, q.NextNode)

	// ================================
	// Repeat but convert back and forth

	server_steps := make([]string, 0, 10)

	// add starting point
	server_steps = append(server_steps, qStart.NextNode)

	// walk once and add step
	r := qStart.ToReader()
	q, err = QueryStructFromStream(&r)
	if err != nil {
		t.Fatal(err)
	}
	q = q.next(data)[0]
	server_steps = append(server_steps, q.NextNode)

	// walk once and add step
	r = q.ToReader()
	q, err = QueryStructFromStream(&r)
	if err != nil {
		t.Fatal(err)
	}
	q = q.next(data)[0]
	server_steps = append(server_steps, q.NextNode)

	// walk once and add step
	r = q.ToReader()
	q, err = QueryStructFromStream(&r)
	if err != nil {
		t.Fatal(err)
	}
	q = q.next(data)[0]
	server_steps = append(server_steps, q.NextNode)

	// ================================
	// compare result

	for i := 0; i < 4; i++ {
		if internal_steps[i] != server_steps[i] {
			t.Fatalf("step %d, they are not equal %s != %s", i, internal_steps[i], server_steps[i])
		}

	}

}
