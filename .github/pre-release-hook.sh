#!/bin/sh
set -e

# Update the version in the README from GitHub Tags
sed -i "s/v0.0.0/${GITHUB_REF_NAME}/" internal/metadata/metadata.go
