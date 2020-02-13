package codeindex

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const pad = 256

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

func (p *Indexer) Run(dir string, relative string) error {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	linksdir := fmt.Sprintf("%s/src/links", os.Getenv("HOME"))

	p.nodes = make([]node, 0, 4096)
	err = p.run(dir, relative)
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

	err = os.RemoveAll(linksdir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(linksdir, 0755)
	if err != nil {
		return err
	}

	format := fmt.Sprintf("%%04d   %%-%ds %%s", longest+3)
	for i, f := range p.nodes {
		noslashes := strings.ReplaceAll(f.relative, "/", "／")
		tmp := fmt.Sprintf(format, i+1, f.name, noslashes)
		other := strings.Repeat(" ", 233-len(tmp))
		list := []string{tmp, other}
		f.new = strings.Join(list, "")
		target := filepath.Join(dir, f.relative, f.name)
		source := filepath.Join(linksdir, f.new)

		err := os.Symlink(target, source)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Indexer) run(dir string, relative string) error {
	finfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range finfo {
		if f.IsDir() {
			err := p.run(filepath.Join(dir, f.Name()), filepath.Join(relative, f.Name()))
			if err != nil {
				return err
			}
		}
	}
	for _, f := range finfo {
		if !f.IsDir() {
			if strings.HasSuffix(f.Name(), p.Ending) {
				p.nodes = append(p.nodes, node{
					name:     f.Name(),
					sort:     fmt.Sprintf("%s%s", relative, strings.Repeat("z", pad-len(relative))),
					relative: relative,
				})
			}
		}
	}

	return nil
}
