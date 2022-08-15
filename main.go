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
		os.Exit(2)
	}

	if _, err := os.Stat(cfg.Source); err != nil {
		fmt.Printf("File %v does not exist\n", cfg.Source)
	}

	// fsysSource := os.DirFS(cfg.Source)
	// err := fs.WalkDir(fsysSource, ".",
	// 	func(p string, d fs.DirEntry, err error) error {
	// 		fmt.Print(d.Name)
	// 		d.Info()
	// 		tm, e := time.Parse(cfg.Layout, d.Name())
	// 		if err != nil {
	// 			return err
	// 		}
	// 		year := string(tm.Year())
	// 		month := string(int(tm.Month()))
	// 		e = os.Rename(path.Join(cfg.Source, d.Name()), path.Join(cfg.Destination, year, month, d.Name()))
	// 		if e != nil {
	// 			return err
	// 		}
	// 		fmt.Printf("File %v mode to %v \r", path.Join(cfg.Source, d.Name()), path.Join(cfg.Destination, year, month, d.Name()))
	// 		return nil
	// 	})

	files, err := os.ReadDir(cfg.Source)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	num := strings.Count(cfg.Layout, "*")
	layout := cfg.Layout[num:]
	count := 0
	for _, file := range files {
		fn := file.Name()
		if len(fn) < num+len(layout) {
			continue
		}
		tm, e := time.Parse(layout, fn[num:num+len(layout)])
		if e != nil {
			fmt.Println(e)
			os.Exit(1)
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
				fmt.Println(err)
				os.Exit(2)
			}
		}

		e = os.Rename(path.Join(cfg.Source, fn), path.Join(cfg.Destination, year, month, fn))
		if e != nil {
			fmt.Println(e)
			os.Exit(1)
		}
		count++
		fmt.Printf("File %v mode to %v. All moved files %v \r", path.Join(cfg.Source, fn), path.Join(cfg.Destination, year, month, fn), count)
	}

	if err != nil {
		fmt.Printf("Error: %v", err)
	}

}
