package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	// "unicode"
	// "testing"
)

func copied_code() {

	match, _ := regexp.MatchString("p([a-z]+)ch", "peach")
	fmt.Println(match)

	r, _ := regexp.Compile("p([a-z]+)ch")

	fmt.Println(r.MatchString("peach"))

	fmt.Println(r.FindString("peach punch"))

	fmt.Println("idx:", r.FindStringIndex("peach punch"))

	fmt.Println(r.FindStringSubmatch("peach punch"))

	fmt.Println(r.FindStringSubmatchIndex("peach punch"))

	fmt.Println(r.FindAllString("peach punch pinch", -1))

	fmt.Println("all:", r.FindAllStringSubmatchIndex(
		"peach punch pinch", -1))

	fmt.Println(r.FindAllString("peach punch pinch", 2))

	fmt.Println(r.Match([]byte("peach")))

	r = regexp.MustCompile("p([a-z]+)ch")
	fmt.Println("regexp:", r)

	fmt.Println(r.ReplaceAllString("a peach", "<fruit>"))

	in := []byte("a peach")
	out := r.ReplaceAllFunc(in, bytes.ToUpper)
	fmt.Println(string(out))
}

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
	NextNode()
}
type TraverseNode struct {
	Parent Node
	Left   Node
	Right  Node
}
type LeafNode struct {
	Parent Node
	Value  string
}
type LoopNode struct {
    Parent Node
    Left Node
    Right Node
}
// NextNode implements Node.
func (t TraverseNode) NextNode() {
	panic("unimplemented")
}
func (t LeafNode) NextNode() {
	panic("unimplemented")
}
func (t LoopNode) NextNode() {
	panic("unimplemented")
}


func grow_tree(str string, parent Node) Node {
    parts := split_q(str)
    Left, operator, Right := parts[0], []rune(parts[1])[0], parts[2] 

    if operator == '/' {
        t := TraverseNode{}
        t.Parent = parent
        t.Left = grow_tree(Left, t)
        t.Right = grow_tree(Right, t)
        return t
    } else if operator == '*' {
        l := LoopNode{}
        l.Left = grow_tree(Left, l)
        l.Right = grow_tree(Right, l)
        return l
    } else if operator == '0' {
        l := LeafNode{Value: Left, Parent: parent}
        return l
    } else {
        panic("invalid operator")
    }
}


func split_q(str string) [3]string {
    opened := 0
    closed := 0
    
	for i, char := range str {
        if char == '{' {
            opened += 1
            continue
        }
        if char == '}' {
            closed += 1
            continue
        } 

        if opened == closed{
            // if no paranthasees are opned
            if strings.Contains("/*&|", string(char)) {
                return [...]string{str[:i], string(str[i]), str[i:]}
            }
        } else if opened > closed {
            // do nothing
        } 
	}
    panic("something went wrong")
}

func main() {
	// re := regexp.MustCompile("\\/")
	// txt1 := "S/Pickaxe/obtainedBy/crafting_recipe/hasInput"
	txt2 := "S/Pick/{Made_of | Crafting_recipie}/rarity"
	txt2 = removeWhitespace(txt2)
	// index := 2
	// last := ???
	// tmp := RegSplit(txt2, "\\/") // split expression on / into parts
	// fmt.Println(tmp)
	// start := tmp[0]
	// tmp = append(tmp[:0], tmp[1:]...)
	// root := &BinaryTreeNode{Edge: start}
	// fmt.Print(tmp)
	re := regexp.MustCompile("^(.*?)\\/")

	root := TraverseNode{}
	// root.Left
	GrowTree(root, txt2)

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
