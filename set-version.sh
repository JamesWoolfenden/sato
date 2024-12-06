#!/bin/sh
# set-version.sh: Updates version.go file with latest git tag
# Usage: ./set-version.sh
# Requires: git
set -e  # Exit on error

if ! latesttag=$(git describe --tags); then
    echo "Error: Failed to get git tag" >&2
    exit 1
fi

if [ -z "$latesttag" ]; then
    echo "Error: No git tags found" >&2
    exit 1
fi

if ! echo "$latesttag" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+'; then
    echo "Error: Invalid version tag format" >&2
    exit 1
fi

echo "Updating version file with new tag: $latesttag"
echo "package version" > src/version/version.go
echo "" >> src/version/version.go
echo "const Version = \"$latesttag\"" >> src/version/version.go
