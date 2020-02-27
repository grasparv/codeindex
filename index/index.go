package codeindex

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/grasparv/codeindex/stats"
)

const pad = 256
const columns = 233
const spacing = 3

type Indexer struct {
	Ending string

	nodes nodelist
}

type node struct {
	name     string
	sort     string
	relative string
	new      string
}

type nodelist []node

func (n nodelist) Len() int {
	return len(n)
}

func (n nodelist) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n nodelist) Less(i, j int) bool {
	if n[i].sort == n[j].sort {
		return n[i].name < n[j].name
	}
	return n[i].sort < n[j].sort
}

func (p *Indexer) Run(st *stats.FileStats, dir string) error {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	linksfile := fmt.Sprintf("%s/.go.links", os.Getenv("HOME"))

	p.nodes = make([]node, 0, 4096)
	err = p.run(st, dir, "")
	if err != nil {
		return err
	}
	sort.Sort(p.nodes)

	longest := 0
	for _, f := range p.nodes {
		if len(f.name) > longest {
			longest = len(f.name)
		}
	}

	bld := strings.Builder{}

	for _, f := range p.nodes {
		target := filepath.Join(f.relative, f.name)
		bld.WriteString(target)
		bld.WriteString("\n")
	}

	var fh *os.File

	_, err = os.Stat(linksfile)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		err = os.Remove(linksfile)
		if err != nil {
			return err
		}
	}

	fh, err = os.Create(linksfile)
	if err != nil {
		return err
	}

	n, err := fh.WriteString(bld.String())
	fh.Close()
	if err != nil {
		return err
	}

	if n != len(bld.String()) {
		return errors.New("not all contents written")
	}

	err = os.Chmod(linksfile, 0400)
	if err != nil {
		return err
	}

	return nil
}

func (p *Indexer) run(st *stats.FileStats, dir string, relative string) error {
	finfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range finfo {
		if f.IsDir() {
			err := p.run(st, filepath.Join(dir, f.Name()), filepath.Join(relative, f.Name()))
			if err != nil {
				return err
			}
		}
	}
	for _, f := range finfo {
		if !f.IsDir() {
			if strings.HasSuffix(f.Name(), p.Ending) {
				fullname := filepath.Join(dir, f.Name())
				frequency := 0
				if v, ok := st.Entries[fullname]; ok {
					frequency = v
				} else {
					frequency = 999999999
				}
				fs := fmt.Sprintf("%010d", frequency)
				p.nodes = append(p.nodes, node{
					name:     f.Name(),
					sort:     fmt.Sprintf("%s%s%s", fs, relative, strings.Repeat("z", pad-len(relative))),
					relative: relative,
				})
			}
		}
	}

	return nil
}

func stringReverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}
