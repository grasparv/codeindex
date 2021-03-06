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

const minexist = 1.0
const pruneratio = 0.1
const recentuse = 48

type FileStat struct {
	Shortname string    `json:"name"`
	Count     int       `json:"count"`
	Date      time.Time `json:"timestamp"`
}

type FileStats struct {
	Entries map[string]*FileStat `json:"entries"`
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
			Entries: make(map[string]*FileStat),
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
		st.Entries = make(map[string]*FileStat)
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
		st.Entries[absname] = &FileStat{
			Shortname: absname,
			Count:     v.Count + 1,
			Date:      time.Now(),
		}
	} else {
		st.Entries[absname] = &FileStat{
			Shortname: absname,
			Count:     1,
			Date:      time.Now(),
		}
	}

	// auto-prune old or rarely-used files
	for k, v := range st.Entries {
		if v.IsTooOld() {
			delete(st.Entries, k)
		}
	}

	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (e FileStat) ratio() float64 {
	age := time.Since(e.Date)
	hours := math.Round(age.Hours())
	if hours == 0 {
		return 4.0
	}
	return float64(e.Count) / hours
}

func (e FileStat) IsTooOld() bool {
	age := time.Since(e.Date)
	hours := math.Round(age.Hours())
	if hours > minexist && e.ratio() < pruneratio {
		return true
	}
	return false
}

func (e FileStat) GetScore() int {
	//frequency := e.Count
	//dur := time.Since(e.Date).Hours()
	//if dur < recentuse {
	//	factor := (float64(recentuse) - dur) // e.g. 7.2 for most-recent hour
	//	addition := e.Count * int(factor)    // e.g. + 7.2 * 145 = 1044
	//	frequency = frequency + addition
	//	fmt.Printf("for %s, factor=%f, count=%d, addition=%d, frequency=%d\n", f.Name(), factor, v.Count, addition, frequency)
	//}
	//return frequency
	if e.Count > 5 && time.Since(e.Date).Hours() < recentuse {
		return int(time.Since(e.Date) / time.Minute)
	}
	return -1
}

func (e FileStat) Description() string {
	dur := time.Since(e.Date).Hours()
	return fmt.Sprintf("%5dx %4.1fh %5.1fx/h %5dpts %s\n", e.Count, dur, e.ratio(), e.GetScore(), e.Shortname)
}
