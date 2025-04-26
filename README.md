# photo-metadata-organizer

Organizes photos based on their EXIF metadata.

- currently supports jpeg format
- creates directories based on year/month (e.g. `2022/01`) in photos metadata
- can add tags to images using exiftool

## Usage

```
-d string
    destination directory
-q	dry run doesn't actually move files
-s string
    source directory
-t string
    shortcut for -tag
-tag string
    tags to add to images (can be specified multiple times)
-w	overwrite destination file, if it exists
```

### Examples

```
# Organize photos with tags
./app -s /path/to/photos -d /destination/path -tag Wedding -tag Family

# Using the shortcut for tags
./app -s /path/to/photos -d /destination/path -t Wedding -t Family

# Dry run (doesn't move files or add tags)
./app -s /path/to/photos -d /destination/path -t Holiday -q
```

## Development

```
 Choose a make command to run

  vet           vet code
  test          run unit tests
  build         build a binary
  autobuild     auto build when source files change
  dockerbuild   build project into a docker container image
  start         build and run local project
  deploy        build code into a container and deploy it to the cloud dev environment
```
