package main

import (
	"log"

	"github.com/matthewaveryusa/funfile"
)

func main() {
	var f funfile.FunFile
	if err := f.Connect("file:memory:?mode=memory", 1); err != nil {
		log.Fatal(err)
	}
	f.AddFile("a/b/c", []byte{1, 2, 3})
	f.GetFile("a/b/c")
	f.GetFile("a")
	defer f.Disconnect()
}
