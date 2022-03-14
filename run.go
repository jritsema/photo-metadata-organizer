package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

func run(args []string, stdout io.Writer) error {

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		dryRun, overwrite bool
		target, dest      string
	)
	flags.StringVar(&target, "t", "", "target directory")
	flags.StringVar(&dest, "d", "", "destination directory")
	flags.BoolVar(&dryRun, "q", false, "dry run doesn't actually move files")
	flags.BoolVar(&overwrite, "w", false, "overwrite destination file, if it exists")
	err := flags.Parse(args[1:])
	check(err)
	fmt.Println("target =", target)
	fmt.Println("dest =", dest)
	fmt.Println("Dry Run =", dryRun)
	fmt.Println()
	if target == "" || dest == "" {
		flags.Usage()
		return errors.New("")
	}

	files, err := ioutil.ReadDir(target)
	check(err)
	for _, file := range files {
		fmt.Println()
		if !file.IsDir() &&
			strings.HasSuffix(file.Name(), ".jpg") ||
			strings.HasSuffix(file.Name(), ".jpeg") ||
			strings.HasSuffix(file.Name(), ".JPG") {

			fmt.Printf("processing %v/%v \n", target, file.Name())
			f, err := os.Open(path.Join(target, file.Name()))
			check(err)

			//get the create date
			x, err := exif.Decode(f)
			check(err)
			tm, err := x.DateTime()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Taken:", tm)

			//format {year}/{month}/{filename}
			parts := strings.Split(tm.String(), "-")
			year := parts[0]
			month := parts[1]
			dest := path.Join(dest, year, month, file.Name())
			t := path.Join(target, file.Name())
			fmt.Printf("moving to %v \n", dest)
			if fileExists(dest) {
				fmt.Println("already exists")
				if !overwrite {
					continue
				}
			}
			if !dryRun {
				containingDir := path.Dir(dest)
				err = os.MkdirAll(containingDir, os.ModePerm)
				check(err)
				err = os.Rename(t, dest)
				check(err)
			}
		}
	}

	return nil
}

func fileExists(f string) bool {
	if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
