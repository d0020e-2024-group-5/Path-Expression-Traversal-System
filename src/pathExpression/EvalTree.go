package pathExpression

import (
	"regexp"
	"strings"
)

// interface type that all nodes need to implement
// it add the possibility to query for the next node
// returns the leaf nodes that are next in the query
type Node interface {
	NextNode(Node) []*LeafNode
	GetLeaf(int) *LeafNode
}

// A Traverse Node represent a traversal from right to left
type TraverseNode struct {
	Parent Node
	Left   Node
	Right  Node
}

// A leaf node represent en edge in the query, these are also the leafs in the evaluation tree
type LeafNode struct {
	Parent Node
	Value  string
	ID     int
}

// Struct representing a loop in the query structure
type LoopNode struct {
	Parent Node
	Left   Node
	Right  Node
}

// The root of a query tree, passes next node to child with required info
type RootNode struct {
	Child Node
}

func (r *RootNode) GetLeaf(id int) *LeafNode {
	return r.Child.GetLeaf(id)
}

func (t *TraverseNode) GetLeaf(id int) *LeafNode {
	tmp1 := t.Left.GetLeaf(id)
	tmp2 := t.Right.GetLeaf(id)

	if tmp1 != nil {
		return tmp1
	} else if tmp2 != nil {
		return tmp2
	} else {
		return nil
	}
}

func (l *LoopNode) GetLeaf(id int) *LeafNode {
	tmp1 := l.Left.GetLeaf(id)
	tmp2 := l.Right.GetLeaf(id)

	if tmp1 != nil {
		return tmp1
	} else if tmp2 != nil {
		return tmp2
	} else {
		return nil
	}
}

func (l *LeafNode) GetLeaf(id int) *LeafNode {
	if l.ID == id {
		return l
	}
	return nil
}

// Passes next node to child with required info
// TODO error can occur when a child calls this function
func (r *RootNode) NextNode(caller Node) []*LeafNode {
	return r.Child.NextNode(r)
}

// node that implements the traverse function.
// If left node calls traverse, continue to right tree
// if right tree calls, branch has been evaluated and call next tree on parent
func (t *TraverseNode) NextNode(caller Node) []*LeafNode {
	var leafs []*LeafNode

	if caller == t.Parent { // Pointer comparison to avoid same value struct bug
		// if the caller is parent, we should deced into the left branch,
		// t.Left.NextNode(t)
		leafs = append(leafs, t.Left.NextNode(t)...)

	} else if caller == t.Left {
		// when the left branch has evaluated it will call us again
		// an we than have to evaluate the right branch
		// t.Right.NextNode(t)
		leafs = append(leafs, t.Right.NextNode(t)...)

	} else if caller == t.Right {
		// when the right brach has evaluated it will call us again
		// we then know we have been fully evaluated and can call our parent saying we are done
		t.Parent.NextNode(t)
		leafs = append(leafs, t.Parent.NextNode(t)...)

	} else {
		panic("i dont know what should happen here1?")
	}
	return leafs
}

// we want to have an looping behavior, when left is done evaluating we want to repeat it which means
// when left calls an we evaluate left again, but exiting the loop is also viable si right should also evaluate.
// right is the exit so when right is done pass ask parent for next node
func (l *LoopNode) NextNode(caller Node) []*LeafNode {
	var leafs []*LeafNode
	if caller == l.Parent {
		tmp1 := l.Left.NextNode(l)
		tmp2 := l.Right.NextNode(l)
		leafs = append(leafs, tmp1...)
		leafs = append(leafs, tmp2...)
	} else if caller == l.Left {
		tmp1 := l.Left.NextNode(l)
		tmp2 := l.Right.NextNode(l)
		leafs = append(leafs, tmp1...)
		leafs = append(leafs, tmp2...)
	} else if caller == l.Right {
		tmp1 := l.Parent.NextNode(l)
		leafs = append(leafs, tmp1...)
	} else {
		panic("loopnode nextnode panic")
	}
	return leafs
}

// This node represents an edge in the query
func (l *LeafNode) NextNode(caller Node) []*LeafNode {
	if caller == l.Parent {
		return []*LeafNode{l}
	} else if caller == nil {
		return l.Parent.NextNode(l)
	} else {
		panic("leafnode nextnode panic")
	}
}

// creates a branch where the top node has the given parent.
// return a the top node
func grow_tree(str string, parent Node, id *int) Node {
	// split the string to a left operator and right part
	// this functions takes into account brackets {}
	parts := split_q(str)
	// fmt.Printf("%v\n", parts)
	Left, operator, Right := parts[0], parts[1], parts[2]

	// if the operator is traverse create a traverse node
	if operator == "/" {
		t := TraverseNode{}
		t.Parent = parent
		t.Left = grow_tree(Left, &t, id)
		t.Right = grow_tree(Right, &t, id)
		return &t

		// if the operator is loop (aka match zero or more) create a loop node
	} else if operator == "*" {
		l := LoopNode{}
		l.Parent = parent
		l.Left = grow_tree(Left, &l, id)
		l.Right = grow_tree(Right, &l, id)
		return &l

		// if the operator is "0" this indicates that the parsing resulted in only a left side
		// this means this is an leaf node
	} else if operator == "0" {
		l := LeafNode{Value: Left, ID: *id, Parent: parent}
		tmp := *id
		tmp++
		*id = tmp
		return &l

		// No operator matched
	} else {
		panic("invalid operator")
	}
}

// if the passed string contains a valid operator,
// note that this returns true even for operators that are planed but not implanted
func containsOperators(s string) bool {
	re := regexp.MustCompile(`[/*&|]`)
	return re.MatchString(s)
}

// this splits the query into the first evaluated operator off the string and its left and right sides.
// this string “{recipe/input}*price/currency“ would split into “{recipe/input}“ “*“ “price/currency“
// as the first evaluated operator is *
func split_q(str string) [3]string {

	// maybe this should return error?
	// when is it the case that an empty would be passed?
	// except in these cases "s/pick/recipe/"
	// or "s//error_here/recipe"
	// both of these might be better to check beforehand?
	if str == "" {
		return [3]string{"", "0", ""}

		// if no operator is contained
	} else if !containsOperators(str) {
		// if the string is contained inside brackets remove them
		if str[0] == '{' && str[len(str)-1] == '}' {
			str = str[1:]
			str = str[:len(str)-1]
			// why call split again?
			// we know no operators are in the string
			// this will just end up att the return bellow, ``return [3]string{str, "0", ""}``
			return split_q(str)
		} else {
			// since we have no operators or enclosing brackets this is an edge
			return [3]string{str, "0", ""}
		}
	}
	//is inside brackets counter
	insideCount := 0
	// for each character advance until the operator is found
	for i, char := range str {
		if insideCount > 0 {
			if char == '}' {
				//remove level of inside brackets
				insideCount = insideCount - 1
			}
			if i == len(str)-1 {
				str = str[1:]
				str = str[:len(str)-1]
				// if this was the last and char encountered and its an closing bracket
				// the whole statement is enclosed on one
				// remove brackets and parse again
				// we should never be att the end and not be on a closing bracket
				return split_q(str)
			}
		} else {
			if char == '{' {
				//add level of inside brackets
				insideCount = insideCount + 1

			} else {
				// The first operator found outside of brackets it the one we want to split on
				if strings.Contains("/*&|", string(char)) {
					return [...]string{str[:i], string(str[i]), str[i+1:]}
				}
			}
		}
	}

	// noting should be able to cause this if the query is correctly formed?
	// only two ways i see into this is "{hello/test}hej"
	// and two ways i see into this is "he{llo/test"
	// TODO add test in the beginning if an non operator is after an closing bracket, example "}h" or before and opening bracket "a{"
	panic("something went wrong, no operators to split on")
}
