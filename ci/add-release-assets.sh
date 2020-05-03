#!/usr/bin/env bash

latestTag=""

function deleteTempFolder() {
  echo "deleting temp folder"
  cd ../../
  rm -rf temp/
}

# workaround while we try to discover how to the get the
# git tags within git action environment
function getLatestTag() {
  echo "cloning into the temp folder"
  mkdir temp && cd temp || exit 1

  git clone https://github.com/azure-open-tools/event-hubs-cli.git
  cd event-hubs-cli || exit 1

  git fetch --all
  latestTag="$(git --no-pager tag -l | tail -1)"
  deleteTempFolder
}

getLatestTag
# git for-each-ref --sort=creatordate --format '%(creatordate)'
echo "$PWD"
echo "Latest Tag: $latestTag"
#find . -type f -iname "ehs-*" -exec ls -lah {} \;

# sender
find . -type f -iname "ehs-*" -exec hub release edit -m "" -a {} "$latestTag" \;

# receiver
find . -type f -iname "ehr-*" -exec hub release edit -m "" -a {} "$latestTag" \;