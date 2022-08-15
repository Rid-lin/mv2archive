// Binary bot-auth-manual implements example of custom session storage and
// manually setting up client options without environment variables.

package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rid-lin/mv2archive/internal/config"
)

func main() {
	cfg := config.New()

	if cfg.S != "" {
		cfg.Source = cfg.S
	} else if cfg.Source == "" {
		fmt.Println("Source folder do not specifity")
		os.Exit(1)
	}
	if cfg.D != "" {
		cfg.Destination = cfg.S
	} else if cfg.Destination == "" {
		fmt.Println("Destination folder do not specifity")
		os.Exit(1)
	}

	if _, err := os.Stat(cfg.Source); err != nil {
		fmt.Printf("File %v does not exist\n", cfg.Source)
		os.Exit(2)
	}

	files, err := os.ReadDir(cfg.Source)
	if err != nil {
		fmt.Printf("Cannot read folder %v: %v", cfg.Source, err)
		os.Exit(2)
	}

	num := strings.Count(cfg.Layout, "*")
	layout := cfg.Layout[num:]
	count := 0

	for _, file := range files {
		fn := file.Name()
		if len(fn) < num+len(layout) {
			fmt.Printf("\n File not contains date in right format")
			continue
		}
		tm, e := time.Parse(layout, fn[num:num+len(layout)])
		if e != nil {
			fmt.Printf("\n Cannot parse time from file %v: %v", fn, e)
			continue
		}
		year := fmt.Sprint(tm.Year())
		month := fmt.Sprint(int(tm.Month()))
		if len(month) == 1 {
			month = "0" + month
		}
		destinationPath := path.Join(cfg.Destination, year, month)
		if _, err := os.Stat(destinationPath); err != nil {
			err := os.MkdirAll(destinationPath, 0755)
			if err != nil {
				fmt.Printf("\n Cannot create folder %v: %v", destinationPath, err)
				continue
			}
		}

		e = os.Rename(path.Join(cfg.Source, fn), path.Join(cfg.Destination, year, month, fn))
		if e != nil {
			fmt.Printf("\n Cannot move file %v to %v: %v", path.Join(cfg.Source, fn), path.Join(cfg.Destination, year, month, fn), e)
			continue
		}
		count++
		fmt.Printf("File %v mode to %v. All moved files %v \r", path.Join(cfg.Source, fn), path.Join(cfg.Destination, year, month, fn), count)
	}

	if err != nil {
		fmt.Printf("Error: %v", err)
	}

}
