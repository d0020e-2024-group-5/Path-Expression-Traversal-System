package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("Hello world, from docker container running go")
		time.Sleep(2 * time.Second)
	}
}
