package mermaidError

import (
	"fmt"
	"io"
	"math/rand"
)

// this functions returns the string repenting a comment node with the given message,
// it also has an edge from "from" to its self with the label "with_edge"
func MermaidErrorEdge(writer io.Writer, from string, with_edge string, message string) {
	err_id := rand.Int()
	fmt.Fprintf(writer, "%s -->|%s| %d@{ shape: braces, label: \"%s\" }\n", from, with_edge, err_id, message)
}
