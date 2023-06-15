#!/bin/sh
latesttag=$(git describe --tags)
echo "Updating version file with new tag: $latesttag"
echo "package version" > src/version/version.go
echo "" >> src/version/version.go
echo "const Version = \"$latesttag\"" >> src/version/version.go
