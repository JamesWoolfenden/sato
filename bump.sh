#!/bin/bash

# Parameters
message="${1:-new release}"

# Version pattern
versionPattern='^[0-9]+\.[0-9]+\.[0-9]+$'
version=''

# Get the current version
version=$(git describe --tags --abbrev=0 2>/dev/null)
version=${version//v}
if [[ ! $version =~ $versionPattern ]]; then
    echo "Invalid version format. Expected: x.y.z"
    exit 1
fi

# Split the version and increment the build number
IFS='.' read -r major minor build <<< "$version"
newBuild=$((build + 1))
newVersion="$major.$minor.$newBuild"

if [[ ! "$newVersion" > "$version" ]]; then
    echo "New version must be greater than current version"
    exit 1
fi

# Output the current and new version
echo "Current version: $version"
echo "New version: $newVersion"
echo "Creating new tag..."

# Create a new tag and push it
git tag -a "v$newVersion" -m "$message"
git push origin "v$newVersion"
