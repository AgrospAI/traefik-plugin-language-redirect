#!/bin/bash

VERSION_FILE=".version"
README_FILE="README.md"

if [ ! -f "$VERSION_FILE" ]; then
  echo "0.0.0" > "$VERSION_FILE"
fi

CURRENT_VERSION=$(cat "$VERSION_FILE")
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"
PATCH=$((PATCH + 1))
NEW_VERSION="$MAJOR.$MINOR.$PATCH"
echo "$NEW_VERSION" > "$VERSION_FILE"
sed -i '' "s/version: $CURRENT_VERSION/version: $NEW_VERSION/" "$README_FILE"
echo "Bumped version from $CURRENT_VERSION to $NEW_VERSION"
