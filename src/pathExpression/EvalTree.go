package pathExpression

import (
	"errors"
	"log"
	"regexp"
	"slices"
	"strings"
)

// interface type that all nodes need to implement
// it add the possibility to query for the next node
// returns the leaf nodes that are next in the query
type Node interface {
	NextNode(Node, []string) []*LeafNode
	GetLeaf(int) *LeafNode
}

// The root of a query tree, passes next node to child with required info
type RootNode struct {
	Child Node
}

// A Traverse Node represent a traversal from right to left
type TraverseNode struct {
	Parent   Node
	Children []Node
}

// Struct representing a loop in the query structure
type LoopNode struct {
	Parent   Node
	Children []Node
}

// Struct representing the OR node operator
type ORNode struct {
	Parent   Node
	Children []Node
}

type ANDNode struct {
	Parent   Node
	Children []Node
}

type XORNode struct {
	Parent   Node
	Children []Node
}

// type PNode struct {
// 	Parent Node
// 	Children []Node
// }

// A leaf node represent en edge in the query, these are also the leafs in the evaluation tree
type LeafNode struct {
	Parent Node
	Value  string
	ID     int
}

// getleaf does DFS and returns *leafnode with matching id
func (r *RootNode) GetLeaf(id int) *LeafNode {
	return r.Child.GetLeaf(id)
}

func (t *TraverseNode) GetLeaf(id int) *LeafNode {
	// returns the first node where id matches or nil
	for i, _ := range t.Children {
		tmp := t.Children[i].GetLeaf(id)
		if tmp != nil {
			return tmp
		}
	}
	return nil
}

func (l *LoopNode) GetLeaf(id int) *LeafNode {
	// returns the first node where id matches or nil
	for i, _ := range l.Children {
		tmp := l.Children[i].GetLeaf(id)
		if tmp != nil {
			return tmp
		}
	}
	return nil
}

func (o *ORNode) GetLeaf(id int) *LeafNode {
	// returns the first node where id matches or nil
	for i, _ := range o.Children {
		tmp := o.Children[i].GetLeaf(id)
		if tmp != nil {
			return tmp
		}
	}
	return nil
}

func (a *ANDNode) GetLeaf(id int) *LeafNode {
	// returns the first node where id matches or nil
	for i, _ := range a.Children {
		tmp := a.Children[i].GetLeaf(id)
		if tmp != nil {
			return tmp
		}
	}
	return nil
}

func (x *XORNode) GetLeaf(id int) *LeafNode {
	// returns the first node where id matches or nil
	for i, _ := range x.Children {
		tmp := x.Children[i].GetLeaf(id)
		if tmp != nil {
			return tmp
		}
	}
	return nil
}

func (l *LeafNode) GetLeaf(id int) *LeafNode {
	if l.ID == id {
		return l
	}
	return nil
}

// Passes next node to child with required info
// TODO error can occur when a child calls this function
func (r *RootNode) NextNode(caller Node, availablePaths []string) []*LeafNode {
	return r.Child.NextNode(r, availablePaths)
}

// node that implements the traverse function.
// If left node calls traverse, continue to right tree
// if right tree calls, branch has been evaluated and call next tree on parent
func (t *TraverseNode) NextNode(caller Node, availablePaths []string) []*LeafNode {
	var leafs []*LeafNode
	// if caller is parent we check the "first" node
	if caller == t.Parent {
		leafs = append(leafs, t.Children[0].NextNode(t, availablePaths)...)
		// then we check all the following Children
	} else if caller != t.Children[len(t.Children)-1] {
		for i, n := range t.Children {
			if caller == n {
				leafs = append(leafs, t.Children[i+1].NextNode(t, availablePaths)...)
				break
			}
		}
		// untill we reach the last chil where we call the parent
	} else if caller == t.Children[len(t.Children)-1] {
		leafs = append(leafs, t.Parent.NextNode(t, availablePaths)...)
	} else {
		panic("Should not happen!")
	}
	return leafs
}

// we want to have an looping behavior, when left is done evaluating we want to repeat it which means
// when left calls an we evaluate left again, but exiting the loop is also viable si right should also evaluate.
// right is the exit so when right is done pass ask parent for next node
func (l *LoopNode) NextNode(caller Node, availablePaths []string) []*LeafNode {
	var leafs []*LeafNode

	// if caller is parent we return all Children paths
	if caller == l.Parent {
		for i, _ := range l.Children {
			leafs = append(leafs, l.Children[i].NextNode(l, availablePaths)...)
		}
		// if caller is not last child we return all Childrens paths
	} else if caller != l.Children[len(l.Children)-1] {
		for i, _ := range l.Children {
			leafs = append(leafs, l.Children[i].NextNode(l, availablePaths)...)
		}
		// if child is last child we call parents nextnode
	} else if caller == l.Children[len(l.Children)-1] {
		leafs = append(leafs, l.Parent.NextNode(l, availablePaths)...)
	} else {
		panic("Should not happen!")
	}
	return leafs
}

// split query
func (o *ORNode) NextNode(caller Node, availablePaths []string) []*LeafNode {
	var leafs []*LeafNode

	// if caller is parent return all Childrens paths
	if caller == o.Parent {
		for i, _ := range o.Children {
			leafs = append(leafs, o.Children[i].NextNode(o, availablePaths)...)
		}
		// if caller is any of the Children return parents nextnode
	} else {
		leafs = append(leafs, o.Parent.NextNode(o, availablePaths)...)
	}

	return leafs
}

// split if both
func (a *ANDNode) NextNode(caller Node, availablePaths []string) []*LeafNode {
	var leafs []*LeafNode

	// if caller return all paths, if path does not exist return empty slice
	if caller == a.Parent {
		for i, _ := range a.Children {
			leafs = append(leafs, a.Children[i].NextNode(a, availablePaths)...)
		}

		// check if leafs.value exist as an available path
		for _, leaf := range leafs {
			isPath := slices.Contains(availablePaths, leaf.Value)
			// if path is not an available path return empty array
			if !isPath {
				var tmp []*LeafNode
				return tmp
			}
		}
	} else {
		leafs = append(leafs, a.Parent.NextNode(a, availablePaths)...)
	}
	return leafs
}

// split if only path
func (x *XORNode) NextNode(caller Node, availablePaths []string) []*LeafNode {
	var leafs []*LeafNode

	// if parent called returns the the one path that exists if that is the only path
	if caller == x.Parent {
		for i, _ := range x.Children {
			leafs = append(leafs, x.Children[i].NextNode(x, availablePaths)...)
		}

		numOfPaths := 0
		for _, leaf := range leafs {
			isPath := slices.Contains(availablePaths, leaf.Value)
			// if path is not an available path return empty array
			if isPath {
				numOfPaths += 1
			}
		}
		if numOfPaths == 1 {
			for _, leaf := range leafs {
				isPath := slices.Contains(availablePaths, leaf.Value)
				if isPath {
					var tmp []*LeafNode
					tmp = append(tmp, leaf)
					return tmp
				}
			}
		}
		// else call parents nextnode
	} else {
		leafs = append(leafs, x.Parent.NextNode(x, availablePaths)...)
	}
	return leafs
}

// This node represents an edge in the query
func (l *LeafNode) NextNode(caller Node, availablePaths []string) []*LeafNode {
	if caller == l.Parent {
		return []*LeafNode{l}
	} else if caller == nil {
		return l.Parent.NextNode(l, availablePaths)
	} else {
		panic("leafnode nextnode panic")
	}
}

// creates a branch where the top node has the given parent.
// return a the top node
func grow_tree(str string, parent Node, id *int) (Node, error) {
	//pre process

	// splitq split the string to parts that are children of the operator
	operator, parts := split_q(str)
	// fmt.Printf("%v\n", parts)

	// if the operator is traverse create a traverse node
	if operator == "/" {
		t := TraverseNode{}
		t.Parent = parent
		for _, part := range parts {
			child, _ := grow_tree(part, &t, id)
			t.Children = append(t.Children, child)
		}
		return &t, nil

		// if the operator is loop (aka match zero or more) create a loop node
	} else if operator == "*" {
		l := LoopNode{}
		l.Parent = parent
		for _, part := range parts {
			child, _ := grow_tree(part, &l, id)
			l.Children = append(l.Children, child)
		}
		return &l, nil

		// if the operator is a OR, create a ORNode
	} else if operator == "|" {
		o := ORNode{}
		o.Parent = parent
		for _, part := range parts {
			child, _ := grow_tree(part, &o, id)
			o.Children = append(o.Children, child)
		}
		return &o, nil

		// if the operator is a AND, create a ANDNode
	} else if operator == "&" {
		a := ANDNode{}
		a.Parent = parent
		for _, part := range parts {
			child, _ := grow_tree(part, &a, id)
			a.Children = append(a.Children, child)
		}
		return &a, nil

		// if the operator is a XOR, create a XORNode
	} else if operator == "^" {
		x := XORNode{}
		x.Parent = parent
		for _, part := range parts {
			child, _ := grow_tree(part, &x, id)
			x.Children = append(x.Children, child)
		}
		return &x, nil

		// this means this is an leaf node
	} else if operator == "0" {
		l := LeafNode{}
		if parts == nil {
			l = LeafNode{Value: "", ID: *id, Parent: parent}
		} else {
			l = LeafNode{Value: parts[0], ID: *id, Parent: parent}
		}
		tmp := *id
		tmp++
		*id = tmp
		return &l, nil

		// No operator matched
	} else {
		panic("invalid operator")
	}
}

// if the passed string contains a valid operator,
func containsOperators(s string) bool {
	re := regexp.MustCompile(`[/*&|^]`)
	return re.MatchString(s)
}

// this splits the query into the first evaluated operator off the string and its left and right sides.
// this string “{recipe/input}*price/currency“ would split into “{recipe/input}“ “*“ “price/currency“
// as the first evaluated operator is *
func split_q(str string) (string, []string) {
	var previusOperators []string
	var parts []string

	// if empty return empty
	if str == "" {
		return "0", nil
		// if no operator is contained
	} else if !containsOperators(str) {
		// if the string is contained inside brackets remove them
		if str[0] == '{' && str[len(str)-1] == '}' {
			return split_q(str[1 : len(str)-1])
		} else {
			// since we have no operators or enclosing brackets this is an edge
			return "0", []string{str}
		}
	}

	insideCount := 0
	for i, char := range str {

		if insideCount > 0 {
			// decrement insidecount if we find closing bracket
			if char == '}' {
				insideCount -= 1
			}
			// since last element and inside brackets remove them
			if i == len(str)-1 {
				return split_q(str[1 : len(str)-1])
			}
		} else if insideCount == 0 {
			// increment open brackets
			if char == '{' {
				insideCount += 1
			} else {
				// if thing contains operator take left of operator into parts
				// and save the operator in the operator list, all operators should be the same according to pre proccesing
				// this will solit it up into all parts,
				// "A&B&C"  => ["&", "&"] and ["A", "B", "C"]
				if containsOperators(string(char)) {
					parts = append(parts, str[:i])
					previusOperators = append(previusOperators, string(str[i]))
				}
			}
		} else {
			panic("insideCount should never be negative")
		}
	}
	return previusOperators[0], parts
}

// Checks for invalid operator combinations and returns nil if no invalid combination is found.
func IsValid(str string) error {
	operands := "/^*&|" // current available operands
	right := 0
	left := 0
	index := strings.IndexAny(str, operands)
	if string(str[index]) != "/" { // if first operand isn't a traverse (/)
		return errors.New("Error; First operator is " + string(str[index]) + " , not traverse (/)") // return error
	}
	for i := 0; i < (len(str) - 1); i++ {
		// checks if last character is an operand (with the exception of / or *)
		if strings.Contains(operands, string(str[len(str)-1])) && !(strings.Contains("*", string(str[len(str)-1])) || strings.Contains("/", string(str[len(str)-1]))) {
			return errors.New("Error; Invalid operand as last character"+string(str[len(str)-1]))
		}
		char := str[i]
		log.Print(string(str[i]))
		if i == len(str)-2 {
			log.Print(string(str[i+1]))
			if str[i+1] == '}' {
				log.Print("right")
				right += 1
			}
			if str[i+1] == '{' {
				log.Print("left")
				left += 1
			}
			if str[i] == '}' {
				log.Print("right")
				right += 1
			}
			if str[i] == '{' {
				log.Print("left")
				left += 1
			}
		} else {
			if str[i] == '}' {
				log.Print("right")
				right += 1
			}
			if str[i] == '{' {
				log.Print("left")
				left += 1
			}
		}

		if string(char)+string(str[i+1]) == "*/" { // exception is */ which is equal to *
			continue
		}
		if string(char) == "{" { // invalid combination { and op
			if strings.Contains(operands, string(str[i+1])) {
				return errors.New("Error; Invalid group operand combination: " + (string(char) + string(str[i+1])))
			}
		}
		if strings.Contains(operands, string(char)) { // if current char is operand
			if string(str[i+1]) == "}" { // invalid combination op and }
				return errors.New("Error; Invalid group operand combination: " + (string(char) + string(str[i+1])))
			}

			if strings.Contains(operands, string(str[i+1])) { // check right char for invalid operand
				return errors.New("Error; Invalid operand combination: " + (string(char) + string(str[i+1])))
			}
		}

	}
	if right != left {
		return errors.New("Error; Unequal amount of left (" + string(left) + ")and right (" + string(right) + ") brackets")
	}
	return nil
}
