#!/bin/bash

set -e
set -x

# Get latest tag
latest_tag=$(git tag --sort=-v:refname | head -n 1)

if [[ -z "$latest_tag" ]]; then
  # If no tags, start from v1.0.0
  next_tag="v1.0.0"
else
  # Split into components
  IFS='.' read -r major minor patch <<< "${latest_tag#v}"
  patch=$((patch + 1))
  next_tag="v${major}.${minor}.${patch}"
fi

# Tag and push
git tag "$next_tag"
git push origin "$next_tag"

echo "Released $next_tag"