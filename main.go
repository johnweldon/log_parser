package main

import (
	"os"
)

func main() {
	p := NewParser(os.Stdin)
	p.Parse()
}
