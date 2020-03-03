package status

import (
	"sort"
	"strings"

	"github.com/grasparv/codeindex/stats"
)

type statList []stats.FileStat

func (n statList) Len() int {
	return len(n)
}

func (n statList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n statList) Less(i, j int) bool {
	return n[i].GetScore() > n[j].GetScore()
}

func Status(fs *stats.FileStats) (string, error) {
	b := strings.Builder{}
	list := make(statList, 0, len(fs.Entries))
	for _, v := range fs.Entries {
		list = append(list, *v)
	}
	sort.Sort(list)
	for _, f := range list {
		b.WriteString(f.Description())
	}
	return b.String(), nil
}
