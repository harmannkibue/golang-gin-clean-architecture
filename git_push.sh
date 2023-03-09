#!/usr/bin/env bash

# shellcheck disable=SC2006
BRANCH=`git rev-parse --abbrev-ref HEAD`

echo "PUSHING TO BRANCH ${BRANCH}"

git add .

git commit

git push --verbose origin "${BRANCH}"
