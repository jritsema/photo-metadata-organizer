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
		dryRun       bool
		target, dest string
	)
	flags.StringVar(&target, "t", "", "target directory")
	flags.StringVar(&dest, "d", "", "destination directory")
	flags.BoolVar(&dryRun, "q", false, "dry run doesn't actually move files")
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
			check(err)
			fmt.Println("Taken:", tm)

			//format {year}/{month}/{filename}
			parts := strings.Split(tm.String(), "-")
			year := parts[0]
			month := parts[1]
			dest := path.Join(dest, year, month, file.Name())
			t := path.Join(target, file.Name())
			fmt.Printf("moving to %v \n", dest)

			if !dryRun {
				containingDir := path.Dir(dest)
				err = os.MkdirAll(containingDir, os.ModePerm)
				check(err)
				err = os.Rename(t, dest)
				check(err)
			}

			fmt.Println()
		}
	}

	return nil
}
