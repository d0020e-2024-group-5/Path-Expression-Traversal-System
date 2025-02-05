package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Removes whitespace in a string
func removeWhitespace(inp string) string {
	return strings.Join(strings.Fields(inp), "")
}

// interface type that all nodes need to implement
// it add the possibility to query for the next node
// returns the leaf nodes that are next in the query
type Node interface {
	NextNode(Node) []*LeafNode
<<<<<<< HEAD
}	
=======
}

// A Traverse Node represent a traversal from right to left
>>>>>>> 881939739dd8826072e46e3a833dabe2e707b98f
type TraverseNode struct {
	Parent Node
	Left   Node
	Right  Node
}

// A leaf node represent en edge in the query, these are also the leafs in the evaluation tree
type LeafNode struct {
	Parent Node
	Value  string
	ID int
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

// Passes next node to child with required info
func (r *RootNode) NextNode(caller Node) []*LeafNode {
	return r.Child.NextNode(r)
}
<<<<<<< HEAD
=======

// node that implements the traverse function.
// If left node calls traverse, continue to right tree
// if right tree calls, branch has been evaluated and call next tree on parent
>>>>>>> 881939739dd8826072e46e3a833dabe2e707b98f
func (t *TraverseNode) NextNode(caller Node) []*LeafNode {
	var leafs []*LeafNode

	if &caller == &(*t).Parent { // Pointer comparison to avoid same value struct bug
		// if the caller is parent, we should deced into the left branch,
		// t.Left.NextNode(t)
		leafs = append(leafs, t.Left.NextNode(t)...)

	} else if &caller == &(*t).Left {
		// when the left branch has evaluated it will call us again
		// an we than have to evaluate the right branch
		// t.Right.NextNode(t)
		leafs = append(leafs, t.Right.NextNode(t)...)

	} else if &caller == &(*t).Right {
		// when the right brach has evaluated it will call us again
		// we then know we have been fully evaluated and can call our parent saying we are done
		t.Parent.NextNode(t)
		leafs = append(leafs, t.Parent.NextNode(t)...)

	} else {
		panic("i dont know what should happen here1?")
	}
	return leafs
}

// TODO nil should be able to call next node, this is because we need to call next node from outside of the tree
// This node represents an edge in the query
func (l *LeafNode) NextNode(caller Node) []*LeafNode {
	if &caller == &(*l).Parent {
		return []*LeafNode{l}
	} else {
		panic("leafnode nextnode panic")
	}
}

// we want to have an looping behavior, when left is done evaluating we want to repeat it which means
// when left calls an we evaluate left again, but exiting the loop is also viable si right should also evaluate.
// right is the exit so when right is done pass ask parent for next node
func (l *LoopNode) NextNode(caller Node) []*LeafNode {
	var leafs []*LeafNode
	if &caller == &(*l).Parent {
		tmp1 := l.NextNode(l)
		tmp2 := l.NextNode(l)
		leafs = append(leafs, tmp1...)
		leafs = append(leafs, tmp2...)
	} else if &caller == &(*l).Left {
		tmp1 := l.NextNode(l)
		tmp2 := l.NextNode(l)
		leafs = append(leafs, tmp1...)
		leafs = append(leafs, tmp2...)
	} else if &caller == &(*l).Right {
		tmp1 := l.Parent.NextNode(l)
		leafs = append(leafs, tmp1...)
	} else {
		panic("loopnode nextnode panic")
	}
	return leafs
}

<<<<<<< HEAD
func grow_tree(str string, parent Node, id *int) Node {
=======
// creates a branch where the top node has the given parent.
// return a the top node
func grow_tree(str string, parent Node) Node {
	// split the string to a left operator and right part
	// this functions takes into account brackets {}
>>>>>>> 881939739dd8826072e46e3a833dabe2e707b98f
	parts := split_q(str)
	fmt.Printf("%v\n", parts)
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

<<<<<<< HEAD


=======
// if the passed string contains a valid operator,
// note that this returns true even for operators that are planed but not implanted
>>>>>>> 881939739dd8826072e46e3a833dabe2e707b98f
func containsOperators(s string) bool {
	re := regexp.MustCompile(`[/*&|]`)
	return re.MatchString(s)
}

// this splits the query into the first evaluated operator off the string and its left and right sides.
// this string “{recipe/input}*price/currency“ would split into “{recipe/input}“ “*“ “price/currency“
// as the first evaluated operator is *
func split_q(str string) [3]string {
	fmt.Println(str)
	opened := 0
	closed := 0

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
	// for each character advance until the operator is found
	for i, char := range str {
		// fmt.Println(string(char))

		// count opening brackets
		if char == '{' {
			opened += 1

			// count closing brackets
		} else if char == '}' {
			closed += 1

			// if this was the last and char encountered and its an closing bracket
			// the whole statement is enclosed on one
			// remove brackets and parse again
			// TODO should we not check if closed == opened, if they are not equal we have an error in the query
			if i == len(str)-1 {
				str = str[1:]
				str = str[:len(str)-1]
				fmt.Println(char)
				return split_q(str)
			}
		}

		// if we are not inside a bracket
		if opened == closed {
			// if its an valid opertor att hte position split there
			if strings.Contains("/*&|", string(char)) {
				return [...]string{str[:i], string(str[i]), str[i+1:]}
			}

			// TODO why this empty else if?
		} else if opened > closed {
			// do nothing
		}
	}

	// noting should be able to cause this if the query is correctly formed?
	// only two ways i see into this is "{hello/test}hej"
	// and two ways i see into this is "he{llo/test"
	// TODO add test in the beginning if an non operator is after an closing bracket, example "}h" or before and opening bracket "a{"
	panic("something went wrong, no operators to split on")
}

func main() {
	txt2 := "Sasdas/Pick/{Made_of}*"
	txt2 = removeWhitespace(txt2)

	id_int := 0
	root := RootNode{}
<<<<<<< HEAD
	tmp := grow_tree(txt2, &root, &id_int)
    root.Child = tmp

    fmt.Println(tmp)
=======
	tmp := grow_tree(txt2, &root)
	root.Child = tmp

	fmt.Println(tmp)
>>>>>>> 881939739dd8826072e46e3a833dabe2e707b98f

	// re := regexp.MustCompile("^(.*?)\\/(.*)")
	// match := re.FindStringSubmatch(txt2)
	// fmt.Println(match)
	// root := TraverseNode{}
	// root.Left = &LeafNode{Value: match[1]}
	// tmp := grow_tree(txt2, &root)
	// fmt.Println(tmp)

	// tmp5 := split_q("")
	// fmt.Println(tmp5)
}

// `S/Pickaxe/obtainedBy/crafting_recipe/hasInput`
// The example will start att pickaxe and follow edge `obtainedBy` to `Pickaxe_From_Stick_And_Stone_Recipe`
// where the query will split and go to both `Cobblestone` and `stick`.
// Since this is the end of the query they are returned

// ### Example 2, Loop
// S/Pick/{Made_of}*/pog
// Will see what pick is made of recursivly down to its minimal component

// ### Example 3, Or
// S/Pick/{Made_of | Crafting_recipie}
// Will return either what it was made of or its crafting recipie

// ### Example 4, AND
// S/Pick/{Made_of & Crafting_recipie}
// Will return both what it is made of and its crafting recipie

// ### Example 5, ()
// S/Pick/{(made_of & Crafting_recipie)/made_of}
