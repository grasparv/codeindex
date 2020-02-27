package main

import (
	"flag"
	"fmt"

	"github.com/grasparv/codeindex/index"
	"github.com/grasparv/codeindex/stats"
)

func main() {
	var ending string
	flag.BoolVar
	flag.StringVar(&ending, "s", ".go", "extension to index")
	flag.Parse()
	dirs := flag.Args()
	if len(dirs) != 2 {
		fmt.Printf("invalid usage\n")
		return
	}

	if dirs[0] == "scan" {
		indexer := &codeindex.Indexer{
			Ending: ending,
		}
		err := indexer.Run(dirs[1], "")
		if err != nil {
			fmt.Printf("error %s\n", err.Error())
			return
		}
	} else if dirs[0] == "stats" {
		err := stats.Update(dirs[1])
		if err != nil {
			fmt.Printf("error %s\n", err.Error())
			return
		}
	}
}
