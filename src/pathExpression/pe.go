package main

import (
	// "bytes"

	"fmt"
	"regexp"
	"strings"
	// "unicode"
	// "testing"
)

// https://stackoverflow.com/questions/4466091/split-string-using-regular-expression-in-go
func RegSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:len(text)]
	return result
}

func removeWhitespace(inp string) string {
	return strings.Join(strings.Fields(inp), "")
}

type Node interface {
	NextNode(Node) []*LeafNode
}	
type TraverseNode struct {
	Parent Node
	Left   Node
	Right  Node
}
type LeafNode struct {
	Parent Node
	Value  string
	ID int
}
type LoopNode struct {
	Parent Node
	Left   Node
	Right  Node
}
type RootNode struct {
	Child Node
}

// NextNode implements Node.
// func (r RootNode) NextNode(Node) []*LeafNode {
// 	panic("unimplemented")
// }

func (r *RootNode) NextNode(caller Node) []*LeafNode {
	return r.Child.NextNode(r)
}
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
func (l *LeafNode) NextNode(caller Node) []*LeafNode {
	if &caller == &(*l).Parent {
		return []*LeafNode{l}
	} else {
		panic("leafnode nextnode panic")
	}
}
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

func grow_tree(str string, parent Node, id *int) Node {
	parts := split_q(str)
    fmt.Printf("%v\n",parts)
	Left, operator, Right := parts[0], parts[1], parts[2]
    
	if operator == "/" {
		t := TraverseNode{}
		t.Parent = parent
		t.Left = grow_tree(Left, &t, id)
		t.Right = grow_tree(Right, &t, id)
		return &t
	} else if operator == "*" {
		l := LoopNode{}
		l.Left = grow_tree(Left, &l, id)
		l.Right = grow_tree(Right, &l, id)
		return &l
	} else if operator == "0" {
		l := LeafNode{Value: Left, ID: *id, Parent: parent}
		tmp := *id
		tmp++
		*id = tmp
		return &l
	} else {
		panic("invalid operator")
	}
}



func containsOperators(s string) bool {
	re := regexp.MustCompile(`[/*&|]`)
	return re.MatchString(s)
}

func split_q(str string) [3]string {
    fmt.Println(str)
	opened := 0
	closed := 0
	if str == "" {
		return [3]string{"", "0", ""}
	} else if !containsOperators(str){
        if str[0] == '{' && str[len(str)-1] == '}' {
            str = str[1:]
            str = str[:len(str)-1]
            return split_q(str)
        } else {
            return [3]string{str, "0", ""}
        }
    }

	for i, char := range str {
        // fmt.Println(string(char))

		if char == '{' {
			opened += 1
		} else if char == '}' {
			closed += 1
			if i == len(str)-1 {
				str = str[1:]
				str = str[:len(str)-1]
                fmt.Println(char)
				return split_q(str)
			}
		}

		if opened == closed {
			// if no paranthasees are opned
			if strings.Contains("/*&|", string(char)) {
				return [...]string{str[:i], string(str[i]), str[i+1:]}
			} 
		} else if opened > closed {
			// do nothing
		}
	}

	panic("something went wrong")
}

func main() {
	txt2 := "Sasdas/Pick/{Made_of}*"
	txt2 = removeWhitespace(txt2)

	id_int := 0
	root := RootNode{}
	tmp := grow_tree(txt2, &root, &id_int)
    root.Child = tmp

    fmt.Println(tmp)

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
