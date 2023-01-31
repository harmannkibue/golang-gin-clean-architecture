#!/usr/bin/env bash

# shellcheck disable=SC2006
BRANCH=`git rev-parse --abbrev-ref HEAD`

echo "Pushing to branch ${BRANCH}"

git add .

echo -n "Enter commit message: "
# shellcheck disable=SC2162
read COMMIT

git commit -m "$COMMIT"
git push --verbose origin "${BRANCH}"
