package main

import (
	"flag"
	"fmt"

	"github.com/grasparv/codeindex/index"
	"github.com/grasparv/codeindex/stats"
	"github.com/grasparv/codeindex/status"
)

func main() {
	var ending string
	flag.StringVar(&ending, "s", ".go", "extension to index")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 || len(args) > 2 {
		fmt.Printf("invalid usage\n")
		return
	}

	switch args[0] {
	case "index":
		if len(args) != 2 {
			fmt.Printf("invalid usage\n")
			return
		}
		indexer := &codeindex.Indexer{
			Ending: ending,
		}
		stats, err := stats.Read()
		if err == nil {
			err = indexer.Run(stats, args[1])
		}
		if err != nil {
			fmt.Printf("error %s\n", err.Error())
			return
		}
	case "use":
		if len(args) != 2 {
			fmt.Printf("invalid usage\n")
			return
		}
		stats, err := stats.Read()
		if err == nil {
			err = stats.Update(args[1])
		}
		if err == nil {
			err = stats.Write()
		}
		if err != nil {
			return
		}
	case "status":
		var out string
		stats, err := stats.Read()
		if err == nil {
			out, err = status.Status(stats)
		}
		if err != nil {
			fmt.Printf("error %s\n", err.Error())
			return
		}
		fmt.Print(out)
	default:
		fmt.Printf("unknown command\n")
	}
}
