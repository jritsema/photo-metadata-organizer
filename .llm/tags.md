add a flag to this go project to allow multiple tags to be passed in.  for example:

./app -t AtlantaShower -t Wedding -t WeddingShower

then add code to read this list of tags and use exiftool to write them to the actual .jpg or image file. honor the `-q` flag and only do this if `-q` is not specified. but if `-q` is specified print out the actual exiftool command that you'd run

make sure that code compiles using the `make build` command

then run a test using the following command:

```sh
./app -d /Volumes/photos \
  -t '/Volumes/photos/Wedding/Atlanta Shower' \
  -tag Wedding -tag WeddingShower -tag AtlantaShower -tag JohnAndLindseysWedding \
  -q
```
