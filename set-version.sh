#!/bin/sh
echo "Updating version file with new tag: $GORELEASER_CURRENT_TAG"
echo "package sato" > src/version.go
echo "" >> src/version.go
echo "const Version = \"$GORELEASER_CURRENT_TAG\"" >> src/version.go
