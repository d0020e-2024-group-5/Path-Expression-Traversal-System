package mermaidError

import (
	"fmt"
	"io"
	"math/rand"
)

func MermaidErrorEdge(writer io.Writer, from string, with_edge string, message string) {
	err_id := rand.Int()
	fmt.Fprintf(writer, "%s -->|%s| %d@{ shape: braces, label: \"%s\" }\n", from, with_edge, err_id, message)
}
