#!/bin/sh
latesttag=$(git describe --tags)
echo "Updating version file with new tag: $latesttag"
echo "package sato" > src/version.go
echo "" >> src/version.go
echo "const Version = \"$latesttag\"" >> src/version.go
