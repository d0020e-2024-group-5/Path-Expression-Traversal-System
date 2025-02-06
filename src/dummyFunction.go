package main

import "math/rand"

var dummy_data = [3]string{
	"No Result",
	"Pickaxe",
	"Stick",
}

func sendQuery(queryString string) string {
	return dummy_data[rand.Intn(3)]
}
