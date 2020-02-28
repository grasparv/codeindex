package main

import (
	"flag"
	"fmt"

	"github.com/grasparv/codeindex/index"
	"github.com/grasparv/codeindex/stats"
)

func main() {
	var ending string
	flag.StringVar(&ending, "s", ".go", "extension to index")
	flag.Parse()
	dirs := flag.Args()
	if len(dirs) != 2 {
		fmt.Printf("invalid usage\n")
		return
	}

	switch dirs[0] {
	case "index":
		indexer := &codeindex.Indexer{
			Ending: ending,
		}
		stats, err := stats.Read()
		if err == nil {
			err = indexer.Run(stats, dirs[1])
		}
		if err != nil {
			fmt.Printf("error %s\n", err.Error())
			return
		}
	case "use":
		stats, err := stats.Read()
		if err == nil {
			err = stats.Update(dirs[1])
		}
		if err == nil {
			err = stats.Write()
		}
		if err != nil {
			fmt.Printf("error %s\n", err.Error())
			return
		}
	default:
		fmt.Printf("unknown command\n")
	}
}
