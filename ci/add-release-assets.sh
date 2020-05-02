#!/usr/bin/env bash

git fetch --all
git branch
latestTag="$(git --no-pager tag -l | tail -1)"

# git for-each-ref --sort=creatordate --format '%(creatordate)'
echo "$PWD"
echo "Latest Tag: $latestTag"
#find . -type f -iname "ehs-*" -exec ls -lah {} \;

# sender
find . -type f -iname "ehs-*" -exec hub release edit -m "" -a {} "$latestTag" \;

# receiver
find . -type f -iname "ehr-*" -exec hub release edit -m "" -a {} "$latestTag" \;