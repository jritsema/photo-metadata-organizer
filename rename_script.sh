#!/bin/bash

# Script to rename files by removing "_original" suffix
# Directory path
DIR="/Volumes/photos/Wedding/Atlanta Shower"

# Loop through all files with "_original" suffix
for file in "$DIR"/*_original; do
    # Check if file exists
    if [ -f "$file" ]; then
        # Get the new name by removing "_original" suffix
        newname="${file%_original}"
        
        # Rename the file
        echo "Renaming: $file -> $newname"
        mv "$file" "$newname"
    fi
done

echo "Renaming complete!"
