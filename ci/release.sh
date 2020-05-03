#!/usr/bin/env bash

versionFile=$1
packageVersion=$2

# workaround while we try to discover how to the get the
# git tags within git action environment
function getLatestTag() {
  echo "cloning into the temp folder"
  mkdir temp && cd temp || exit 1
  git clone https://github.com/azure-open-tools/event-hubs-cli.git
  cd event-hubs-cli || exit 1
  git fetch --all
}

function deleteTempFolder() {
  echo "deleting temp folder"
  cd ../../
  rm -rf temp/
}

if [[ $versionFile == *"package"* ]];
then
  version="$packageVersion"
else
  version=$(go run "$versionFile")
fi

getLatestTag
checkTag=$(git --no-pager tag -l | grep "$version" | xargs)
if [[ $checkTag != "" ]];
then
  echo "$checkTag already exist, skipping release."
  deleteTempFolder
  exit 0
fi

latestTag="$(git --no-pager tag -l | tail -1)"
changeLog="$(git --no-pager log --oneline "$latestTag"...HEAD)"

echo "New Version: $version"
echo "Latest Tag: $latestTag"
echo -e "Change Log Since Latest Tag: \n$changeLog"

hub release create -m "Azure Event Hubs Lib $version" -m "$changeLog" "$version"

deleteTempFolder