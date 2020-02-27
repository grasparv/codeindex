package main

import (
	"flag"
	"fmt"

	"github.com/grasparv/codeindex/index"
)

func main() {
	indexer := &codeindex.Indexer{}
	flag.StringVar(&indexer.Ending, "s", ".go", "extension to index")
	flag.Parse()
	dirs := flag.Args()
	if len(dirs) != 1 {
		fmt.Printf("no directory given\n")
		return
	}

	err := indexer.Run(dirs[0], "")
	if err != nil {
		fmt.Printf("error %s\n", err.Error())
		return
	}
}
