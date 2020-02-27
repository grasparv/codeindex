package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileStats struct {
	Entries map[string]int `json:"entries"`
}

func getStatsFilename() string {
	return fmt.Sprintf("%s/.go.stats", os.Getenv("HOME"))
}

func Read() (*FileStats, error) {
	statsfile := getStatsFilename()

	_, err := os.Stat(statsfile)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		st := &FileStats{
			Entries: make(map[string]int),
		}
		return st, nil
	}

	rh, err := os.Open(statsfile)
	if err != nil {
		return nil, err
	}
	defer rh.Close()

	data, err := ioutil.ReadAll(rh)
	if err != nil {
		return nil, err
	}

	var st FileStats
	err = json.Unmarshal(data, &st)
	if err != nil {
		return nil, err
	}

	if st.Entries == nil {
		st.Entries = make(map[string]int)
	}

	return &st, nil
}

func (st *FileStats) Write() error {
	data, err := json.Marshal(&st)
	if err != nil {
		return err
	}

	statsfile := getStatsFilename()
	wh, err := os.Create(statsfile)
	if err != nil {
		return err
	}
	defer wh.Close()

	n, err := wh.Write(data)
	if err != nil {
		return err
	}

	if n != len(data) {
		return errors.New("did not write all data")
	}

	return nil
}

func (st *FileStats) Update(filename string) error {
	absname, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	finfo, err := os.Stat(absname)
	if err != nil {
		return err
	}

	if finfo.IsDir() {
		return errors.New("will not count directories")
	}

	if _, ok := st.Entries[absname]; ok {
		st.Entries[absname]++
	} else {
		st.Entries[absname] = 1
	}

	return nil
}
