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
    result := make([]string, len(indexes) + 1)
    for i, element := range indexes {
            result[i] = text[laststart:element[0]]
            laststart = element[1]
    }
    result[len(indexes)] = text[laststart:len(text)]
    return result
}

func removeWhitespace(inp string) string{
    return strings.Join(strings.Fields(inp), "")
}

func main() {
	// re := regexp.MustCompile("\\/")
    // txt1 := "S/Pickaxe/obtainedBy/crafting_recipe/hasInput"
	txt2 := "S/Pick/{Made_of | Crafting_recipie}/rarity"
	txt2 = removeWhitespace(txt2)

    index := 2
    // last := ???
    
    tmp := RegSplit(txt2, "\\/") // split expression on / into parts
	fmt.Println(tmp[2])

    // nästade listor 
    // lista med indexer beroende på djup av {}









    

    // fmt.Println(set) // ["Have", "a", "great", "day!"]

}

// `S/Pickaxe/obtainedBy/crafting_recipe/hasInput`
// The example will start att pickaxe and follow edge `obtainedBy` to `Pickaxe_From_Stick_And_Stone_Recipe`
// where the query will split and go to both `Cobblestone` and `stick`.
// Since this is the end of the query they are returned

// ### Example 2, Loop
// S/Pick/{Made_of}*
// Will see what pick is made of recursivly down to its minimal component

// ### Example 3, Or
// S/Pick/{Made_of | Crafting_recipie}
// Will return either what it was made of or its crafting recipie

// ### Example 4, AND
// S/Pick/{Made_of & Crafting_recipie}
// Will return both what it is made of and its crafting recipie

// ### Example 5, ()
// S/Pick/{(made_of & Crafting_recipie)/made_of}