package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

func run(args []string) error {

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		dryRun, overwrite bool
		target, dest      string
		tags              multiFlag
	)
	flags.StringVar(&target, "t", "", "target directory")
	flags.StringVar(&dest, "d", "", "destination directory")
	flags.BoolVar(&dryRun, "q", false, "dry run doesn't actually move files")
	flags.BoolVar(&overwrite, "w", false, "overwrite destination file, if it exists")
	flags.Var(&tags, "tag", "tags to add to images (can be specified multiple times)")

	err := flags.Parse(args[1:])
	check(err)
	fmt.Println("target =", target)
	fmt.Println("dest =", dest)
	fmt.Println("Dry Run =", dryRun)
	if len(tags) > 0 {
		fmt.Println("Tags =", tags)
	}
	fmt.Println()
	if target == "" || dest == "" {
		flags.Usage()
		return errors.New("")
	}

	files, err := os.ReadDir(target)
	check(err)
	for _, file := range files {
		fmt.Println()
		if !file.IsDir() &&
			(strings.HasSuffix(strings.ToLower(file.Name()), ".jpg") ||
			strings.HasSuffix(strings.ToLower(file.Name()), ".jpeg")) {

			fmt.Printf("processing %v/%v \n", target, file.Name())
			filePath := path.Join(target, file.Name())
			f, err := os.Open(filePath)
			check(err)

			//get the create date
			x, err := exif.Decode(f)
			if err != nil {
				fmt.Println("skipping, error reading EXIF data:", err)
				f.Close()
				continue
			}
			tm, err := x.DateTime()
			if err != nil {
				fmt.Println("skipping, error reading DateTime:", err)
				f.Close()
				continue
			}
			f.Close()
			fmt.Println("Taken:", tm)

			//format {year}/{month}/{filename}
			parts := strings.Split(tm.String(), "-")
			year := parts[0]
			month := parts[1]
			destPath := path.Join(dest, year, month, file.Name())
			
			// Apply tags if any were specified
			if len(tags) > 0 {
				if dryRun {
					// In dry run mode, show the exiftool command that would be executed
					exifArgs := []string{"exiftool"}
					for _, tag := range tags {
						exifArgs = append(exifArgs, fmt.Sprintf("-XMP-dc:Subject+=%s", tag))
					}
					exifArgs = append(exifArgs, filePath)
					fmt.Printf("Would run: %s (dry run)\n", strings.Join(exifArgs, " "))
				} else {
					if err := applyTags(filePath, tags); err != nil {
						fmt.Printf("Error applying tags to %s: %v\n", file.Name(), err)
					} else {
						fmt.Printf("Applied XMP tags to %s: %v\n", file.Name(), tags)
					}
				}
			}
			fmt.Printf("moving to %v", destPath)
			if dryRun {
				fmt.Println(" (dry run)")
			}
			fmt.Println()
			
			if fileExists(destPath) {
				fmt.Println("destination file already exists")
				if !overwrite {
					continue
				}
			}
			
			if !dryRun {
				containingDir := path.Dir(destPath)
				err = os.MkdirAll(containingDir, os.ModePerm)
				check(err)
				
				err = os.Rename(filePath, destPath)
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

// multiFlag is a custom flag type that can be specified multiple times
type multiFlag []string

func (f *multiFlag) String() string {
	return strings.Join(*f, ", ")
}

func (f *multiFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

// applyTags uses exiftool to write tags to the image file using XMP format
func applyTags(imagePath string, tags []string) error {
	// Create separate -XMP-dc:Subject+= arguments for each tag
	args := []string{}
	for _, tag := range tags {
		args = append(args, fmt.Sprintf("-XMP-dc:Subject+=%s", tag))
	}
	
	// Add the image path as the last argument
	args = append(args, imagePath)
	
	cmd := exec.Command("exiftool", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exiftool error: %v, output: %s", err, string(output))
	}
	
	return nil
}
