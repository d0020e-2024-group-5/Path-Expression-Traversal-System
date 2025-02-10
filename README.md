# Path-Expression-Traversal-System

![The logo, depicting a cat slowly becoming a mesh network](./readme_images/PETS.png)

PETS, a system to store linked distributed data with traversal functions

## Ontologies

An Ontology is a way to describe a reletionship with a stucture of subject, predicate and object. And our data is therefore a list of these structures which can be describe as following:

```mermaid
    graph LR;
        Subject -->|Predicate| Object;
```

We call all subjects and objects nodes and predicates edges.

```mermaid
    graph LR;
        Node1 -->|Edge| Node2;
```

What we want to do is to search such an ontology strucure using a queary where this structure is spread over several servers.

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

    
    note right of Tool_company: outgoing querys can be sent in parallel
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
<!-- Explain in words what happens in the sequence diagram -->

## Node ontologies

```mermaid
graph TD;

        Stick_Plank_made_Instance -->|obtainedBy| Stick_Planks_recipe_Instance

        Stick_Bamboo_made_Instance -->|obtainedBy| Stick_bamboo_recipe_Instance



        Pickaxe_Instance_Henry -->|obtainedBy| PickaxeRecipe_Instance


        Plank_Instance -->|obtainedBy| Plannks_From_Logs_Recipe_Instance


        PickaxeRecipe_Instance -->|hasInput| Stick_Plank_made_Instance
        PickaxeRecipe_Instance -->|hasInput| Stick_Bamboo_made_Instance
        PickaxeRecipe_Instance -->|hasInput| Cobblestone_Bob
        PickaxeRecipe_Instance -->|hasOutput| Pickaxe_Instance_Henry
        PickaxeRecipe_Instance -->|usedInStation| CraftingTable_Instance

        Stick_bamboo_recipe_Instance -->|hasInput| Bamboo_Instance
        Stick_bamboo_recipe_Instance -->|hasOutput| Stick_Bamboo_made_Instance
        Stick_bamboo_recipe_Instance -->|usedInStation| CraftingTable_Instance

        Stick_Planks_recipe_Instance -->|hasInput| Plank_Instance
        Stick_Planks_recipe_Instance -->|hasOutput| Stick_Plank_made_Instance
        Stick_Planks_recipe_Instance -->|usedInStation| CraftingTable_Instance

        Plannks_From_Logs_Recipe_Instance -->|hasInput| Log_Instance
        Plannks_From_Logs_Recipe_Instance -->|hasOutput| Plank_Instance
        Plannks_From_Logs_Recipe_Instance -->|usedInStation| CraftingTable_Instance
```

## Node ontologies distributed

```mermaid
graph TD;

        root

        subgraph Server A
            Stick_Bamboo_made_Instance
            Stick_bamboo_recipe_Instance
            Bamboo_Instance
            CraftingTable_Instance
        end

        subgraph Server B
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

        root -->|Pickaxe_Instance_Henry| Pickaxe_Instance_Henry
        root -->|Cobblestone_Bob| Cobblestone_Bob

        Stick_Plank_made_Instance -->|obtainedBy| Stick_Planks_recipe_Instance

        Stick_Bamboo_made_Instance -->|obtainedBy| Stick_bamboo_recipe_Instance

        Pickaxe_Instance_Henry -->|obtainedBy| PickaxeRecipe_Instance


        Plank_Instance -->|obtainedBy| Plannks_From_Logs_Recipe_Instance


        PickaxeRecipe_Instance -->|hasInput| Stick_Plank_made_Instance
        PickaxeRecipe_Instance -->|hasInput| Stick_Bamboo_made_Instance
        PickaxeRecipe_Instance -->|hasInput| Cobblestone_Bob
        PickaxeRecipe_Instance -->|hasOutput| Pickaxe_Instance_Henry
        PickaxeRecipe_Instance -->|usedInStation| CraftingTable_Instance

        Stick_bamboo_recipe_Instance -->|hasInput| Bamboo_Instance
        Stick_bamboo_recipe_Instance -->|hasOutput| Stick_Bamboo_made_Instance
        Stick_bamboo_recipe_Instance -->|usedInStation| CraftingTable_Instance

        Stick_Planks_recipe_Instance -->|hasInput| Plank_Instance
        Stick_Planks_recipe_Instance -->|hasOutput| Stick_Plank_made_Instance
        Stick_Planks_recipe_Instance -->|usedInStation| CraftingTable_Instance

        Plannks_From_Logs_Recipe_Instance -->|hasInput| Log_Instance
        Plannks_From_Logs_Recipe_Instance -->|hasOutput| Plank_Instance
        Plannks_From_Logs_Recipe_Instance -->|usedInStation| CraftingTable_Instance
```

## Truly distributed data

```mermaid
graph TD;
subgraph Server_a
Server_a_minecraft:Stick_Bamboo_made_Instance([Stick_Bamboo_made_Instance])
end
Server_a_minecraft:Stick_Bamboo_made_Instance-->|obtainedBy|Server_a_minecraft:Stick_bamboo_recipe_Instance

subgraph Server_a
Server_a_minecraft:Server_a[(Server_a)]
end

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

subgraph Server_b
Server_b_minecraft:Stick_Plank_made_Instance([Stick_Plank_made_Instance])
end
Server_b_minecraft:Stick_Plank_made_Instance-->|obtainedBy|Server_b_minecraft:Stick_Planks_recipe_Instance

subgraph Server_b
Server_b_minecraft:Server_a[(Server_a)]
end

subgraph Server_b
Server_b_minecraft:Server_b[(Server_b)]
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
Server_b_minecraft:Stick_Planks_recipe_Instance-->|hasOutput|Server_b_minecraft:Stick_Plank_made_Instance
Server_b_minecraft:Stick_Planks_recipe_Instance-->|usedInStation|Server_b_minecraft:CraftingTable_Instance

subgraph Server_b
Server_b_minecraft:Plannks_From_Logs_Recipe_Instance[Plannks_From_Logs_Recipe_Instance]
end
Server_b_minecraft:Plannks_From_Logs_Recipe_Instance-->|inter_server|Server_c_minecraft:Plannks_From_Logs_Recipe_Instance
Server_b_minecraft:Plannks_From_Logs_Recipe_Instance-->|pointsToServer|Server_b_minecraft:Server_c

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
Server_c_minecraft:Server_c[(Server_c)]
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
Server_c_minecraft:Stick_bamboo_recipe_Instance[Stick_bamboo_recipe_Instance]
end
Server_c_minecraft:Stick_bamboo_recipe_Instance-->|inter_server|Server_a_minecraft:Stick_bamboo_recipe_Instance
Server_c_minecraft:Stick_bamboo_recipe_Instance-->|pointsToServer|Server_c_minecraft:Server_a

subgraph Server_c
Server_c_minecraft:Plannks_From_Logs_Recipe_Instance([Plannks_From_Logs_Recipe_Instance])
end
Server_c_minecraft:Plannks_From_Logs_Recipe_Instance-->|hasInput|Server_c_minecraft:Log_Instance
Server_c_minecraft:Plannks_From_Logs_Recipe_Instance-->|hasOutput|Server_c_minecraft:Plank_Instance
Server_c_minecraft:Plannks_From_Logs_Recipe_Instance-->|usedInStation|Server_c_minecraft:CraftingTable_Instance
```

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

### Example 2, Loop

Looping expressions, matching more than once, allowing for following a path of unknown length. The syntax is the to add a star around a group ``{...}*``

``S/Pickaxe/{obtainedBy/hasInput}*``

Will see what pick is made of recursively down to its minimal component, the paths are

```text
Pickaxe --> Pickaxe_From_Stick_And_Stone_Recipe --> Stick --> Stick_From_Planks_Recipe --> Plank --> Plank_From_Logs_Recipe --> Log
Pickaxe --> Pickaxe_From_Stick_And_Stone_Recipe --> Cobblestone
```

Where both Cobblestone and Log would be returned.

### Example 3, Or

Allows a path traversal to follow either edge

``S/Pickaxe/{obtainedBy/rarity|foundAt}/rarity``

```text
Pickaxe --> Pickaxe_From_Stick_And_Stone_Recipe --> Common
Pickaxe --> Mineshaft --> Rare
```

### Example 4, AND

Only allows the query to continue if both edges exist on the node, both are traversed

``S/Pickaxe/{obtainedBy & foundAt}/rarity`` would return

```text
Pickaxe --> Pickaxe_From_Stick_And_Stone_Recipe --> Common
Pickaxe --> Mineshaft --> Rare
```

``S/Stick/{obtainedBy & foundAt}/rarity`` would return nothing as stick dont have the edge foundAt.

### Example 5, groups {}

TODO EXPLAIN MORE

S/Pick/{(made_of & Crafting_recipie)/made_of}

### example arguments (), TO BE DECIDED

arguments could be added to loop operator?

## Example of internal structure of a query

Lets take an example query af show its internal evaluation

``S/Pickaxe/{obtainedBy/hasInput}*``

This is then converted to a tree structure of operations, where the leafs are edges and.

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

lets say that we are on edge ``obtainedBy``, and we want to know whats next.
By looking at the parent we know that we are on the left side of an *traverse*
and the next edge is the one on the right of the traverse, ``hasInput``

if whe should get the next node from ``hasInput`` we can again look att the parent
and se that we are on the right side of the *traverse*,
to find the next node we need to look higher, the *traverse*'s parent.
This gives us the knowledge that we are on the left side of *loop* operator (aka *zero or more*)
We then have two possible options continue right or redo the left side.
by evaluating the left side we get ``obtainedBy`` again, showing us that the *loop* works.
the right sides gives us NULL, the end of the query an valid position to return.

## go style pseudo code

Note that this pseudo code

```go
type TraverseNode Struct{
    Parent  *Node
    Left    *Node
    Right   *Node
}

// when calling this function we need to know where this was called from, was it our parent, left or right, there for passing a pointer to caller is necessary
// the function return array of pointers to the leafs/query edges, which can be used to determine the next node(s)
func (self TraverseNode) nextEdge(caller *Node) []*LeafNode {
    // if the caller is parent, we should deced into the left branch,
    if caller == self.parent {
        return self.left.nextEdge(&self)
    }
    // when the left branch has evaluated it will call us again
    // an we than have to evaluate the right branch
    else if caller == self.left {
        self.right.nextEdge(&self)
    }
    // when the right brach has evaluated it will call us again
    // we then know we have been fully evaluated and can call our parent saying we are done
    else if caller == self.right {
        self.parent.nextEdege(&self)
    } else {
        log.fatal("i dont know what should happen here?")
    }
}

type LeafNode Struct {
    Parent      *Node
    edgeName    string
}

// we are asked what the next edge is, this LeafNode represent that edge
func (self LeafNode) nextEdge(caller *Node) []*LeafNode {
    if caller == self.parent {
        return [&self]
    } else {
        log.fatal("i dont know what should happen here?")
    }
}

type LoopNode Struct{
    Parent  *Node
    Left    *Node
    Right   *Node
}

func (self LoopNode) nextEdge(caller *Node) []*LeafNode {
    // if the caller is parent, the possible outcomes are that we match zero of the edges and move on with the right branch
    // or that we match whatever in the left brach, therefore we return the next edges
    // an therefore return 
    // if this was match one and more instead of zero or more, calling right would not be the right option as it then would progress forward without having matched anything on the left
    // maybe add an + operator which is match one or more?
    if caller == self.parent {
        return [self.left.nextEdge(&self),self.right.nextEdge(&self)]
    }
    // when the left branch has evaluated it will call us again
    // we can continue the loop so left is an valid option, but we could also exit
    // this leads to the same output as the caller was the parent
    else if caller == self.left {
        return [self.left.nextEdge(&self),self.right.nextEdge(&self)]
    }
    // when the right brach has evaluated it will call us again
    // we then know we have been fully evaluated and can call our parent saying we are done
    else if caller == self.right {
        self.parent.nextEdge(&self)
    } else {
        log.fatal("i dont know what should happen here?")
    }
}

type OrNode Struct{
    Parent  *Node
    Left    *Node
    Right   *Node
}

func (self OrNode) nextEdge(caller *Node) []*LeafNode {
    // if the parent calls us we could either match each side, so both sides are an alternative
    if caller == self.parent {
        return [self.left.nextEdge(&self),self.right.nextEdge(&self)]
    }
    // if the left side calls us we have then completed one of the options and are fully evaluated, an the call our parent saying we are done, and let them get the next edge
    else if caller == self.left {
        self.parent.nextEdge(&self)
    }
    // same ass above
    else if caller == self.right {
        self.parent.nextEdge(&self)
    } else {
        log.fatal("i dont know what should happen here?")
    }
}

```
