package main

import (
	"fmt"

	"github.com/sliceking/galleria/rand"
)

func main() {
	fmt.Println(rand.String(10))
	fmt.Println(rand.RememberToken())
}
