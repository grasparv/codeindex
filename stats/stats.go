package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"time"
)

const oneweek = 120.0
const pruneratio = 5

type FileStat struct {
	Count int       `json:"count"`
	Date  time.Time `json:"timestamp"`
}

type FileStats struct {
	Entries map[string]FileStat `json:"entries"`
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
			Entries: make(map[string]FileStat),
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
		st.Entries = make(map[string]FileStat)
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

	if v, ok := st.Entries[absname]; ok {
		st.Entries[absname] = FileStat{
			Count: v.Count + 1,
			Date:  time.Now(),
		}
	} else {
		st.Entries[absname] = FileStat{
			Count: 1,
			Date:  time.Now(),
		}
	}

	// auto-prune old or rarely-used files
	for k, v := range st.Entries {
		age := time.Since(v.Date)
		if age.Hours() > oneweek && // last use older than one week
			v.Count/int(math.Round(age.Hours())) < pruneratio { // low count/age ratio
			delete(st.Entries, k)
		}
	}

	return nil
}
