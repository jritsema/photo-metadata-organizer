package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
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
		if !file.IsDir() {

			fmt.Printf("processing %v/%v \n", target, file.Name())
			filePath := path.Join(target, file.Name())

			//get the create date
			stdout, stderr, err := execute("exiftool", []string{
				"-CreateDate",
				"-DateTimeOriginal",
				filePath,
			})
			check(err)
			if stderr != "" {
				check(errors.New(stderr))
			}
			fmt.Println(stdout)
			if stdout == "" {
				fmt.Println("create date not found!")
				continue
			}

			//could return one or both "Create Date", "Date/Time Original", or nothing
			lines := strings.Split(stdout, "\n")

			//just take 1st line
			parts := strings.Split(lines[0], " : ")
			createDate := strings.TrimSpace(parts[1])
			tm, err := time.Parse("2006:01:02 15:04:05", createDate)
			check(err)
			fmt.Println("Taken:", tm)

			//format {year}/{month}/{filename}
			parts = strings.Split(tm.String(), "-")
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

//executes a command and returns its stdout and stderr
func execute(name string, args []string) (string, string, error) {

	// Command to execute
	cmd := exec.Command(name, args...)

	// Run the command and capture its stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		return "", "", err
	}

	// Print the stdout and stderr
	return stdout.String(), stderr.String(), nil
}
