package stats

import (
	"os"
)

type fileStat struct {
	count    int
	filename string
}

type fileStats struct {
	entries []fileStat `json:"entries"`
}

func Update(filename string) {
	dir, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	statsfile := fmt.Sprintf("%s/.go.stats", os.Getenv("HOME"))
}
