# Path-Expression-Traversal-System

![The logo, depicting a cat slowly becoming a mesh network](./readme_images/PETS.png)

PETS, a system to store linked distributed data with traversal functions

## Introduction

Our system is designed to navigate and retrieve information from linked ontologies using path expressions. It enables users to traverse decentralized data structures by following relationships defined in path expressions, allowing for multilevel hierarchical exploration.
The system helps users efficiently access and analyze linked data. For example in a supply chain context, it can model the entire distribution network of a specific item by retrieving manufacturer data at each stage. This provides transparency, traceability, and deeper insights into complex data relationships.
Our solution is valuable for businesses, researchers, and data analysts who needs to explore and make sense of decentralized, linked data in a clear ond structured way.

## Table of contents

- [Introduction](#introduction)
- [Table of contents](#table-of-contents)
- [Ontologies](#ontologies)
- [Node ontologies](#node-ontologies)
  - [Node ontologies distributed](#node-ontologies-distributed)
  - [Truly distributed data](#truly-distributed-data)
- [From Query to Result](#from-qurery-to-result)
- [Parsing the ontologies into GoLang](#parsing-the-ontologies-into-golang)
- [Architecture](#architecture)
- [Query structure](#query-structure)
  - [Example 1, Simple traversal](#example-1-simple-traversal)
  - [Example 2, groups {}](#example-2-groups-)
  - [Example 3, Loop](#example-3-loop)
  - [Example 4, Or](#example-4-or)
  - [Example 5, AND](#example-5-and)
  - [Example 6, XOR](#example-6-xor)
- [Parsing and constructing the evaluationTree](#parsing-and-constructing-the-evaluationtree)
- [Traversing the tree](#treversing-the-tree)
- [go style pseudo code](#go-style-pseudo-code)
- [Example of internal structure of a query](#example-of-internal-structure-of-a-query)
  - [An example of evaluation](#an-example-of-evaluation)
- [Current limitations and future development of the query structure](#current-limitations-and-future-development-of-the-query-structure)
- [The query wrapper](#the-query-wrapper)
  - [Syntax Validation](#syntax-validation)
  - [Passing the query to dirent servers](#passing-the-query-to-dirent-servers)
    - [The Common Header](#the-common-header)
    - [Payload for recursive mermaid query (type 0x1)](#payload-for-recursive-mermaid-query-type-0x1)
- [Webserver](#webserver)

## Ontologies

An Ontology is a way to describe a relationship with a structure of subject, predicate and object. And our data is therefore a list of these structures which can be describe as following:

```mermaid
    graph LR;
        Subject -->|Predicate| Object;
```

We call all subjects and objects nodes and predicates edges.

```mermaid
    graph LR;
        Node1 -->|Edge| Node2;
```

What we want to do is to search such an ontology structure using a query where this structure is spread over several servers. When we split an edge over different server we do that by making the first node point to a false node where that node contains all the information to navigate to the true node on the other server.

## Node ontologies

```mermaid
graph TD;

        minecraft:Stick_Plank_made_Instance -->|minecraft:obtainedBy| minecraft:Stick_Planks_recipe_Instance

        minecraft:Stick_Bamboo_made_Instance -->|minecraft:obtainedBy| minecraft:Stick_bamboo_recipe_Instance



        minecraft:Pickaxe_Instance_Henry -->|minecraft:obtainedBy| minecraft:PickaxeRecipe_Instance


        minecraft:Plank_Instance -->|minecraft:obtainedBy| minecraft:Plannks_From_Logs_Recipe_Instance

        Plank_from_Bamboo_Instance -->|obtainedBy| Plank_from_bamboo_recipe_Instance

        minecraft:PickaxeRecipe_Instance -->|minecraft:hasInput| Stick_Plank_made_Instance
        minecraft:PickaxeRecipe_Instance -->|minecraft:hasInput| Stick_Bamboo_made_Instance
        PickaxeRecipe_Instance -->|hasInput| Cobblestone_Bob
        PickaxeRecipe_Instance -->|hasOutput| Pickaxe_Instance_Henry
        PickaxeRecipe_Instance -->|usedInStation| CraftingTable_Instance

        Stick_bamboo_recipe_Instance -->|hasInput| Bamboo_Instance
        Stick_bamboo_recipe_Instance -->|hasOutput| Stick_Bamboo_made_Instance
        Stick_bamboo_recipe_Instance -->|usedInStation| CraftingTable_Instance

        Stick_Planks_recipe_Instance -->|hasInput| Plank_Instance
        Stick_Planks_recipe_Instance -->|hasInput| Plank_from_Bamboo_Instance
        Stick_Planks_recipe_Instance -->|hasOutput| Stick_Plank_made_Instance
        Stick_Planks_recipe_Instance -->|usedInStation| CraftingTable_Instance

        Plannks_From_Logs_Recipe_Instance -->|hasInput| Log_Instance
        Plannks_From_Logs_Recipe_Instance -->|hasOutput| Plank_Instance
        Plannks_From_Logs_Recipe_Instance -->|usedInStation| CraftingTable_Instance

        Plank_from_bamboo_recipe_Instance-->|hasInput|Bamboo_Instance
        Plank_from_bamboo_recipe_Instance-->|hasOutput|Plank_from_Bamboo_Instance
        Plank_from_bamboo_recipe_Instance-->|usedInStation|CraftingTable_Instance
```

### Node ontologies distributed

```mermaid
graph TD;


        subgraph Server A
            Plank_from_bamboo_recipe_Instance
            Stick_Bamboo_made_Instance
            Stick_bamboo_recipe_Instance
            Bamboo_Instance
            CraftingTable_Instance
        end

        subgraph Server B
            Plank_from_Bamboo_Instance
            Plank_Instance
            Stick_Plank_made_Instance
            Stick_Planks_recipe_Instance
        end

        subgraph Server C
            Pickaxe_Instance_Henry
            PickaxeRecipe_Instance
            Plannks_From_Logs_Recipe_Instance
            Cobblestone_Bob
            Log_Instance
        end

        Stick_Plank_made_Instance -->|obtainedBy| Stick_Planks_recipe_Instance

        Stick_Bamboo_made_Instance -->|obtainedBy| Stick_bamboo_recipe_Instance

        Pickaxe_Instance_Henry -->|obtainedBy| PickaxeRecipe_Instance


        Plank_Instance -->|obtainedBy| Plannks_From_Logs_Recipe_Instance

        Plank_from_Bamboo_Instance -->|obtainedBy| Plank_from_bamboo_recipe_Instance

        PickaxeRecipe_Instance -->|hasInput| Stick_Plank_made_Instance
        PickaxeRecipe_Instance -->|hasInput| Stick_Bamboo_made_Instance
        PickaxeRecipe_Instance -->|hasInput| Cobblestone_Bob
        PickaxeRecipe_Instance -->|hasOutput| Pickaxe_Instance_Henry
        PickaxeRecipe_Instance -->|usedInStation| CraftingTable_Instance

        Stick_bamboo_recipe_Instance -->|hasInput| Bamboo_Instance
        Stick_bamboo_recipe_Instance -->|hasOutput| Stick_Bamboo_made_Instance
        Stick_bamboo_recipe_Instance -->|usedInStation| CraftingTable_Instance

        Stick_Planks_recipe_Instance -->|hasInput| Plank_Instance
        Stick_Planks_recipe_Instance -->|hasInput| Plank_from_Bamboo_Instance
        Stick_Planks_recipe_Instance -->|hasOutput| Stick_Plank_made_Instance
        Stick_Planks_recipe_Instance -->|usedInStation| CraftingTable_Instance

        Plannks_From_Logs_Recipe_Instance -->|hasInput| Log_Instance
        Plannks_From_Logs_Recipe_Instance -->|hasOutput| Plank_Instance
        Plannks_From_Logs_Recipe_Instance -->|usedInStation| CraftingTable_Instance

        Plank_from_bamboo_recipe_Instance-->|hasInput|Bamboo_Instance
        Plank_from_bamboo_recipe_Instance-->|hasOutput|Plank_from_Bamboo_Instance
        Plank_from_bamboo_recipe_Instance-->|usedInStation|CraftingTable_Instance
```

### Truly distributed data

```mermaid
graph TD;
subgraph Server_a
Server_a_minecraft:Server_b[(Server_b)]
end

subgraph Server_a
Server_a_minecraft:Stick_Bamboo_made_Instance([Stick_Bamboo_made_Instance])
end
Server_a_minecraft:Stick_Bamboo_made_Instance-->|obtainedBy|Server_a_minecraft:Stick_bamboo_recipe_Instance

subgraph Server_a
Server_a_minecraft:Bamboo_Instance([Bamboo_Instance])
end

subgraph Server_a
Server_a_minecraft:CraftingTable_Instance([CraftingTable_Instance])
end

subgraph Server_a
Server_a_minecraft:Stick_bamboo_recipe_Instance([Stick_bamboo_recipe_Instance])
end
Server_a_minecraft:Stick_bamboo_recipe_Instance-->|hasInput|Server_a_minecraft:Bamboo_Instance
Server_a_minecraft:Stick_bamboo_recipe_Instance-->|hasOutput|Server_a_minecraft:Stick_Bamboo_made_Instance
Server_a_minecraft:Stick_bamboo_recipe_Instance-->|usedInStation|Server_a_minecraft:CraftingTable_Instance

subgraph Server_a
Server_a_minecraft:Plank_from_Bamboo_Instance[Plank_from_Bamboo_Instance]
end
Server_a_minecraft:Plank_from_Bamboo_Instance-->|inter_server|Server_b_minecraft:Plank_from_Bamboo_Instance
Server_a_minecraft:Plank_from_Bamboo_Instance-->|pointsToServer|Server_a_minecraft:Server_b

subgraph Server_a
Server_a_minecraft:Plank_from_bamboo_recipe_Instance([Plank_from_bamboo_recipe_Instance])
end
Server_a_minecraft:Plank_from_bamboo_recipe_Instance-->|hasInput|Server_a_minecraft:Bamboo_Instance
Server_a_minecraft:Plank_from_bamboo_recipe_Instance-->|hasOutput|Server_a_minecraft:Plank_from_Bamboo_Instance
Server_a_minecraft:Plank_from_bamboo_recipe_Instance-->|usedInStation|Server_a_minecraft:CraftingTable_Instance

subgraph Server_b
Server_b_minecraft:Stick_Plank_made_Instance([Stick_Plank_made_Instance])
end
Server_b_minecraft:Stick_Plank_made_Instance-->|obtainedBy|Server_b_minecraft:Stick_Planks_recipe_Instance

subgraph Server_b
Server_b_minecraft:Server_a[(Server_a)]
end

subgraph Server_b
Server_b_minecraft:Server_c[(Server_c)]
end

subgraph Server_b
Server_b_minecraft:Plank_Instance([Plank_Instance])
end
Server_b_minecraft:Plank_Instance-->|obtainedBy|Server_b_minecraft:Plannks_From_Logs_Recipe_Instance

subgraph Server_b
Server_b_minecraft:CraftingTable_Instance[CraftingTable_Instance]
end
Server_b_minecraft:CraftingTable_Instance-->|inter_server|Server_a_minecraft:CraftingTable_Instance
Server_b_minecraft:CraftingTable_Instance-->|pointsToServer|Server_b_minecraft:Server_a

subgraph Server_b
Server_b_minecraft:Stick_Planks_recipe_Instance([Stick_Planks_recipe_Instance])
end
Server_b_minecraft:Stick_Planks_recipe_Instance-->|hasInput|Server_b_minecraft:Plank_Instance
Server_b_minecraft:Stick_Planks_recipe_Instance-->|hasInput|Server_b_minecraft:Plank_from_Bamboo_Instance
Server_b_minecraft:Stick_Planks_recipe_Instance-->|hasOutput|Server_b_minecraft:Stick_Plank_made_Instance
Server_b_minecraft:Stick_Planks_recipe_Instance-->|usedInStation|Server_b_minecraft:CraftingTable_Instance

subgraph Server_b
Server_b_minecraft:Plannks_From_Logs_Recipe_Instance[Plannks_From_Logs_Recipe_Instance]
end
Server_b_minecraft:Plannks_From_Logs_Recipe_Instance-->|inter_server|Server_c_minecraft:Plannks_From_Logs_Recipe_Instance
Server_b_minecraft:Plannks_From_Logs_Recipe_Instance-->|pointsToServer|Server_b_minecraft:Server_c

subgraph Server_b
Server_b_minecraft:Plank_from_Bamboo_Instance([Plank_from_Bamboo_Instance])
end
Server_b_minecraft:Plank_from_Bamboo_Instance-->|obtainedBy|Server_b_minecraft:Plank_from_bamboo_recipe_Instance

subgraph Server_b
Server_b_minecraft:Plank_from_bamboo_recipe_Instance[Plank_from_bamboo_recipe_Instance]
end
Server_b_minecraft:Plank_from_bamboo_recipe_Instance-->|inter_server|Server_a_minecraft:Plank_from_bamboo_recipe_Instance
Server_b_minecraft:Plank_from_bamboo_recipe_Instance-->|pointsToServer|Server_b_minecraft:Server_a

subgraph Server_c
Server_c_minecraft:Stick_Plank_made_Instance[Stick_Plank_made_Instance]
end
Server_c_minecraft:Stick_Plank_made_Instance-->|inter_server|Server_b_minecraft:Stick_Plank_made_Instance
Server_c_minecraft:Stick_Plank_made_Instance-->|pointsToServer|Server_c_minecraft:Server_b

subgraph Server_c
Server_c_minecraft:Stick_Bamboo_made_Instance[Stick_Bamboo_made_Instance]
end
Server_c_minecraft:Stick_Bamboo_made_Instance-->|inter_server|Server_a_minecraft:Stick_Bamboo_made_Instance
Server_c_minecraft:Stick_Bamboo_made_Instance-->|pointsToServer|Server_c_minecraft:Server_a

subgraph Server_c
Server_c_minecraft:Cobblestone_Bob([Cobblestone_Bob])
end

subgraph Server_c
Server_c_minecraft:Log_Instance([Log_Instance])
end

subgraph Server_c
Server_c_minecraft:Pickaxe_Instance_Henry([Pickaxe_Instance_Henry])
end
Server_c_minecraft:Pickaxe_Instance_Henry-->|obtainedBy|Server_c_minecraft:PickaxeRecipe_Instance

subgraph Server_c
Server_c_minecraft:Pickaxe_Instance_Gustav([Pickaxe_Instance_Gustav])
end

subgraph Server_c
Server_c_minecraft:Server_a[(Server_a)]
end

subgraph Server_c
Server_c_minecraft:Server_b[(Server_b)]
end

subgraph Server_c
Server_c_minecraft:Plank_Instance[Plank_Instance]
end
Server_c_minecraft:Plank_Instance-->|inter_server|Server_b_minecraft:Plank_Instance
Server_c_minecraft:Plank_Instance-->|pointsToServer|Server_c_minecraft:Server_b

subgraph Server_c
Server_c_minecraft:CraftingTable_Instance[CraftingTable_Instance]
end
Server_c_minecraft:CraftingTable_Instance-->|inter_server|Server_a_minecraft:CraftingTable_Instance
Server_c_minecraft:CraftingTable_Instance-->|pointsToServer|Server_c_minecraft:Server_a

subgraph Server_c
Server_c_minecraft:PickaxeRecipe_Instance([PickaxeRecipe_Instance])
end
Server_c_minecraft:PickaxeRecipe_Instance-->|hasInput|Server_c_minecraft:Stick_Plank_made_Instance
Server_c_minecraft:PickaxeRecipe_Instance-->|hasInput|Server_c_minecraft:Stick_Bamboo_made_Instance
Server_c_minecraft:PickaxeRecipe_Instance-->|hasInput|Server_c_minecraft:Cobblestone_Bob
Server_c_minecraft:PickaxeRecipe_Instance-->|hasOutput|Server_c_minecraft:Pickaxe_Instance_Henry
Server_c_minecraft:PickaxeRecipe_Instance-->|usedInStation|Server_c_minecraft:CraftingTable_Instance

subgraph Server_c
Server_c_minecraft:Plannks_From_Logs_Recipe_Instance([Plannks_From_Logs_Recipe_Instance])
end
Server_c_minecraft:Plannks_From_Logs_Recipe_Instance-->|hasInput|Server_c_minecraft:Log_Instance
Server_c_minecraft:Plannks_From_Logs_Recipe_Instance-->|hasOutput|Server_c_minecraft:Plank_Instance
Server_c_minecraft:Plannks_From_Logs_Recipe_Instance-->|usedInStation|Server_c_minecraft:CraftingTable_Instance
```

## From Qurery to Result

The user enter in a query, for example
`S/Pickaxe/obtainedBy/crafting_recipe/hasInput` \
that is being sent with a JSON request {"data":"S/Pickaxe/obtainedBy/crafting_recipe/hasInput#100"} to the webserver to be traversed. The number after # represents the ttl thats being inputted at the website, with the default value of 100.
...
...
...

The Program will then attempt to build an evaluation tree.
First it gets passed into a preprocessing step where it removes whitespaces newlines and so on. It then gets checked so that it follows the requirements of the query structure, for example no following operators, always closing parentheses, etc.
It will then call a function to build the tree structure entering in the query and the available edges it can take in the program from the current server. Current limitations limits the program to the current server, in future development it would be ideal to be able to see edges stored in different servers. It will then pass it into the tree-building function
that takes the query and separates the operators into operator-nodes and stores the edges in leaf-nodes. This three is then used to query the DB.
The tree building process will repeat whenever the query is passed to a new server. After it has been passed to a new server it will also call the nextLeaf function that finds the next leaf that has to be visited by traversing the three with an in order walk and returning a pointer to the node.  

The server sends back the traversed query and renders it as a mermaid diagram through the Mermaid.js library.

## Parsing the ontologies into GoLang

For this parsing function, nodes have been defined as struct containing an array (or slices in GoLang) with edges to the node. We also create a struct for edges with the properties "EdgeName" and "TargetName" with each property denoting how an item is obtained respectively what the edge is pointing to. One server can then save all these nodes in a hashmap (dictionary) with the key being the node name and the value, the DataNode struct.

```go
type DataNode struct {
    Edges []DataEdge
}
type DataEdge struct { // minecraft:obtainedBy minecraft:Stick_bamboo_recipe_Instance
    EdgeName   string  // obtainedBy
    TargetName string  // Stick_bamboo_recipe_Instance
}

var nodeLst map[string]DataNode    // Hashmap (or dictionary) with pairs of nodeNames and DataNode
```

Reading the ontologies into Go is very simple. Since the ontologies follow a certain standard (Subject, Predicate, Object) we utilize this to read the subject prefix in order to infer the type of object and similarly what attributes it may have to apply them to our hashmap of nodes. Here is some rough pseudo-code on how the parsing works;

```go
    if strings.HasPrefix(line, "minecraft:") { // minecraft:Stick_Bamboo_made_Instance a minecraft:Stick ;

        temp := strings.TrimPrefix(line, "minecraft:") // Stick_Bamboo_made_Instance a minecraft:Stick ;

        wrd := getFirstWord(temp) // Get the first word in trimmed line; "Stick_Bamboo_made_Instance"

        nodeLst[wrd] = ""  // Declare a key "wrd" with value "" in the hashmap "nodeLst"
    }

    // Next line;       minecraft:obtainedBy minecraft:Stick_bamboo_recipe_Instance
    if strings.HasPrefix(line, "minecraft:EdgeName minecraft:TargetName"){
        var tempEdge DataEdge
        temp := strings.TrimPrefix(line, "minecraft:") //obtainedBy minecraft:Stick_bamboo_recipe_Instance
        firstName := getFirstWord(temp) // Get the first word in trimmed line; "obtainedBy"
        temp := strings.TrimPrefix(temp, "EdgeName minecraft:")
        secondName := getFirstWord(temp) // Get the second word in trimmed line; "Stick_bamboo_recipe_Instance"

        tempEdge.EdgeName = firstName
        tempEdge.TargetName = secondName

        nodeLst[wrd] = append(nodeLst[wrd], tempEdge)  // Append tempEdge to array of edges in node "wrd"
    }
```

## Architecture

```mermaid
sequenceDiagram
    participant API
    participant cent as Central server

    note over API: where request start


    note over cent: A server which listens to http request<br/>and answer/forward queries.<br/>All servers are of this type, but their<br/>databases very<br/>and the central database is quite<br/>different with only one node*.

    API->>+cent: Pick/crafted_by*
    note left of cent: sent over http api

    loop until not stored
        create participant DB_A as DB
        cent->>DB_A: 
        note over DB_A: This database is just a node<br/>which edges point to items in other servers
        destroy DB_A
        DB_A->>cent: 
    end
    note left of DB_A: won't loop since only one node



    participant Tool_company as The tool company
    cent->>+Tool_company: Pick/crafted_by*


    loop until not stored
        create participant DB_B as DB
        Tool_company->> DB_B: 
        note over DB_B: A graph data base<br/>keeps track of nodes<br/>and relations between them.<br/>Also have node to denote<br/>when to pass the query to<br/>different server.
        destroy DB_B
        DB_B->>Tool_company: 
    end

    note left of DB_B: Pickaxe don't have<br/>nodes stone or stick<br/>returns server contact<br/>information.

    
    note right of Tool_company: outgoing queries can be sent in parallel
    participant Mason_LTD as Masons LTD
    Tool_company->>+Mason_LTD: Stone/crafted_by*


    loop until not stored
        create participant DB_C as DB
        Mason_LTD->>DB_C: 
        destroy DB_C 
        DB_C->>Mason_LTD: 
    end
    note left of DB_C: Stone don't have crafted by<br/>and is att end of query, return stone


    Mason_LTD->>-Tool_company: Data of Stone

    participant Wood_INC as Wood INC

    Tool_company->>+Wood_INC: Stick/crafted_by*
    loop until not stored
        create participant DB_D as DB
        Wood_INC->>DB_D: 
        destroy DB_D
        DB_D->>Wood_INC: 
    end
    note left of DB_D: This query can be resolved<br/>without traversing to other<br/>servers<br/>Stick -> Plank -> Log

    Wood_INC->>-Tool_company: Data of Log
    Tool_company->>-cent: [Data of Stone, Data of Log]
    cent->>-API: [Data of Stone, Data of Log]

```
<!-- TODO Explain in words what happens in the sequence diagram -->

## Query structure

The query structure was designed for simplicity and not fines, the goal was an easy way to write path expressions with loops.

```mermaid
graph LR;
    S -->|Pickaxe| Pickaxe;
    S -->|Stick| Stick;

    Pickaxe -->|foundAt| Mineshaft
    Pickaxe -->|obtainedBy| Pickaxe_From_Stick_And_Stone_Recipe;
    Pickaxe_From_Stick_And_Stone_Recipe -->|hasInput| Stick;
    Pickaxe_From_Stick_And_Stone_Recipe -->|hasInput| Cobblestone;

    Mineshaft -->|rarity| Rare
    Pickaxe_From_Stick_And_Stone_Recipe -->|rarity| Common;

    Stick -->|obtainedBy| Stick_From_Planks_Recipe;
    Stick_From_Planks_Recipe -->|hasInput| Plank;

    Plank -->|obtainedBy| Plank_From_Logs_Recipe;
    Plank_From_Logs_Recipe -->|hasInput| Log
```

### Example 1, Simple traversal

To follow a simple path, first have the starting node (s in this case a we have not implemented a dht to resolve node location) followed by the edges name separated by `/`

`S/Pickaxe/obtainedBy/crafting_recipe/hasInput`

The example will start att pickaxe and follow edge `obtainedBy` to `Pickaxe_From_Stick_And_Stone_Recipe`
where the query will split and go to both `Cobblestone` and `stick`.
Since this is the end of the query they are returned.

### Example 2, groups {}

In the above example the query each operation was evaluated left to right, in some cases this might not be desired when using more complex operators such a loop (aka match zero or more)

``S/Pick/made_of/Crafting_recipie*``
``S/Pick/{made_of/Crafting_recipie}*``

In other cases we might want to do more complex
operations, for example an AND or an XOR operation between edges, those are explained in further examples.

<!-- ### example arguments (), TO BE DECIDED

arguments could be added to loop operator? -->

### Example 3, Loop

Looping expressions, matching more than once, allowing for following a path of unknown length. The syntax is the to add a star around a group ``{...}*`` or if only a single edge requires looping the group can be omitted

``S/Pickaxe/{obtainedBy/hasInput}*``

will loop down by the edges obtainedBy/hasInput until it reaches the end

```text
Pickaxe --> Pickaxe_From_Stick_And_Stone_Recipe --> Stick --> Stick_From_Planks_Recipe --> Plank --> Plank_From_Logs_Recipe --> Log
Pickaxe --> Pickaxe_From_Stick_And_Stone_Recipe --> Cobblestone
```

Where both Cobblestone and Log would be returned.

``S/Pickaxe/{obtainedBy/hasInput}*/rarity``

after taking hasInput it will loop if available and go check the rarity if available if neither edge exist it returns as normal.

### Example 4, Or

Allows a path traversal to follow either edge

``S/Pickaxe/{obtainedBy/rarity|foundAt}/rarity``

```text
Pickaxe --> Pickaxe_From_Stick_And_Stone_Recipe --> Common
Pickaxe --> Mineshaft --> Rare
```

### Example 5, AND

Only allows the query to continue if both edges exist on the node, both are traversed

``S/Pickaxe/{obtainedBy & foundAt}/rarity`` would return

```text
Pickaxe --> Pickaxe_From_Stick_And_Stone_Recipe --> Common
Pickaxe --> Mineshaft --> Rare
```

``S/Stick/{obtainedBy & foundAt}/rarity`` would return nothing as stick don't have the edge foundAt.

### Example 6, XOR

Allows the query to continue, only if one of the edges exist

``S/Pickaxe/{obtainedBy ^ foundAt}/rarity`` would return

```text
Pickaxe --> Pickaxe_From_Stick_And_Stone_Recipe --> Common
Pickaxe --> Mineshaft --> Rare
```

``S/Stick/{obtainedBy ^ foundAt}/rarity`` would go down obtainedBy as it does not have the edge foundAt

## Parsing and constructing the evaluationTree

When constructing the evaluationTree the code calls the function func grow_tree(str string, parent Node, id *int) (Node, error)
providing the query string, and the parent node, it returns a node and an error, the error is nil if it did not encounter an error in the function.
It then passes the string for some formatting, removing whitespaces, newlines and so on.

It then calls the function slit_q, where it separates operators from non-operators(edges, bracket-sections or remainder), and returns them.
for example it starts matching traverseNodes, "/", and then finds another operator
everything to the right of the last traverseNode will be treated as a remainder and added to the non-operators.

grow-tree then matches the returned operator and creates a matching node.
The node then takes all the parts and recursively calls growTree on them, then assigns them as a child and appends it to its slice of children.

<https://github.com/user-attachments/assets/4149fed5-c0e7-4b85-a21d-64de8eccd0b8>

## Traversing the tree

When traversing the tree it will call  NextNode(caller Node, availablePaths []string) []*LeafNode,
that takes in the caller node, all the available paths on the server, and returns a slice of leafNode pointers.

Different nodes behave differently when NextNode is called on them but the general behavior is that it calls the next child it has,
and if it was the last child that called it calls its parents nextNode.

The recursive nextNode goes down to the firs leaf it finds (some nodes look for more then one leafNode) and returns a pointer to it.
When NextNode is called from a leaf it will find the next leaf in the evaluation order.

When the query is passed to a new server it has to recreate the tree and must then get a new pointer to the last visited leaf in the newly constructed tree, it then calls GetLeaf(id int) *LeafNode,
that takes in the id of a leaf and returns a pointer to it.

<https://github.com/user-attachments/assets/dee7cb16-68cf-401e-9b38-bfee2e433ea0>

## go style pseudo code

Note this is an example of part of the tree structure.
This example will use the traverse- and leaf- node as an example.

```go
// interface type that all nodes need to implement
// it add the possibility to query for the next node, returns the leaf nodes that are next in the query
// And find the next leaf, returning a pointer to it
type Node interface {
 NextNode(Node, []string) []*LeafNode
 GetLeaf(int) *LeafNode
}

// A Traverse Node represent a traversal from right to left
type TraverseNode struct {
 Parent Node
 Children []Node 
}

// will simple pass by and return a matching leafNode pointer or null (see leafNode)
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

// node that implements the traverse function.
// If the parent calls it checks the fist node
// if a child calls it checks the next node
// if it is the last child it calls parents nextNode nad gets it's next child.
// returns a slice of leafNode pointers, empty if no matches or logic stops it.
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
 // until we reach the last child where we call the parent
 } else if caller == t.Children[len(t.Children)-1] {
  leafs = append(leafs, t.Parent.NextNode(t, availablePaths)...)
 } else {
  panic("Should not happen!")
 }
 return leafs
}

// A leaf node represent en edge in the query, these are also the leafs in the evaluation tree
type LeafNode struct {
 Parent Node
 Value  string
 ID     int
}

// will check if the id matches and return a pointer to itself if it does, nil if it doesn't
func (l *LeafNode) GetLeaf(id int) *LeafNode {
 if l.ID == id {
  return l
 }
 return nil
}

// Returns pointer to self in slice if parent called it
// if call came from nil take parents nextNode instead
func (l *LeafNode) NextNode(caller Node, availablePaths []string) []*LeafNode {
 if caller == l.Parent {
  return []*LeafNode{l}
 } else if caller == nil {
  return l.Parent.NextNode(l, availablePaths)
 } else {
  panic("leafNode nextNode panic")
 }
}

```

## Example of internal structure of a query

<!-- Note to readers, this look incredibly like the state machines that regex compiles to -->
```mermaid
graph TD;
    r([root]);
    r -->|left| s
    r -->|right| 2

    2([/]);
    2 -->|left| Pickaxe
    2 -->|right| 3

    3([\*]);
    3 -->|left| 4
    3 -->|right| NULL

    4([/])
    4 -->|left| obtainedBy
    4 -->|right| hasInput
```

### An example of evaluation

Lets take an example query of show its internal evaluation

``S/Pickaxe/{obtainedBy/hasInput}*``

This is then converted to a tree structure of operations, where the leafs are edges and nodes .

lets say that we are on edge ``obtainedBy``, and we want to know whats next.
By looking at the parent we know that we are on the left side of an *traverse*
and the next edge is the one on the right of the traverse, ``hasInput``

if whe should get the next node from ``hasInput`` we can again look att the parentO
and se that we are on the right side of the *traverse*,
to find the next node we need to look higher, the *traverse*'s parent.
This gives us the knowledge that we are on the left side of *loop* operator (aka *zero or more*)
We then have two possible options continue right or redo the left side.
by evaluating the left side we get ``obtainedBy`` again, showing us that the *loop* works.
the right sides gives us NULL, the end of the query an valid position to return.

## Current limitations and future development of the query structure

<!-- TODO -->

### Syntax Validation

As mentioned previously, checking the syntax is very straightforward. We simply only need to check for invalid operator combinations with a few edge cases such as what starting/ending operators are allowed.

## The query wrapper

The evaluation tree, on its own, does not provide complete utility.
Therefore, a wrapper is constructed to include additional information necessary for graph traversal.
This wrapper specifies the next node to traverse and the edge along which to traverse.

With the aid of the evaluation tree, the next edges and subsequent nodes can be determined, and the traversal path can be updated,
provided we have access to the graph.
This process may necessitate splitting the query into multiple sub-queries in certain cases.
By applying this approach recursively, we can effectively navigate the graph.

Currently, the wrapper, which is responsible for initializing values and invoking the grow tree function,
lacks error handling and query syntax validation. For a system like this to be useful,
stability is crucial, which is currently not the case.

Implementing error handling for simple syntax errors is relatively straightforward.
However, more advanced error handling may require additional time due to potential edge cases,
particularly those involving Unicode.

### Passing the query to dirent servers

To pass data and queries to different servers an common protocol is needed, this is described bellow

#### The Common Header

The common header for all queries Include general information that is needed for the system.
These include an magic to determine if its acutely an PETS query, Query type, an identifying UUID and TTL.
The body to this header is determined by the type.

There might be further changes to this standard to include multiple variable with authorizations tokens in further development. This would allow advanced functions such as to hide internal paths to non authorized users, but still propagate the query thru the network.

```mermaid
packet-beta
title PETS common header
0-31: "Magic 'PETS' "
32-47: "PETS Type"
48-63: "TTL"
64-191: "Query identifier (UUID)"
```

#### Payload for recursive mermaid query (type 0x1)

Since the data might not be stored on the same server,
we need the ability to send the query to the next server and return the results.
Each node that is not stored on the server has a ``false node`` with an edge labeled ``pointsToServer``.
This edge allows us to obtain the contact information needed to forward the query.

The query is then converted from its internal representation to a format suitable for transmission in the payload, in this case a simple string:

``QueryString;NextNode;AlongEdge``

- QueryString: The first part before the ; separator is the query itself, such as ``S/Pickaxe/{obtainedBy/hasInput}*``.
- NextNode: The second part indicates the intended destination server, for example, ``Cobblestone``.
- AlongEdge: The final part is an index into the query, indicating which edge within the query is used to reach the NextNode. For instance, an index of 3 would denote the edge ``hasInput``.

This information is sufficient to reconstruct the query and its state. In the current implementation,
if state values are missing, the query is assumed to be new. The starting node is the first part of the query, and the edge is the same as the starting node.

The return of *recursive mermaid query* is always a valid mermaid string, when error are encountered during the traversal its converted to an valid mermaid syntax,
shush as node with custom styling containing the error message, if possible an arrow to the node that the error occurred in is also drawn.

## Webserver

A webserver is beneficial for our system because it acts as a bridge between users and the linked data processing.
It allows us to interact with our system from anywhere using a simple HTTP request and it provides a unified interface for querying and retrieving linked data.

<!-- TODO write about rendering mermaid -->
