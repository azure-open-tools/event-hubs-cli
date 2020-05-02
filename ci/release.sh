#!/usr/bin/env bash

versionFile=$1
packageVersion=$2

if [[ $versionFile == *"package"* ]];
then
  version="$packageVersion"
else
  version=$(go run "$versionFile")
fi

checkTag=$(git --no-pager tag -l | grep "$version" | xargs)

if [[ $checkTag != "" ]];
then
  echo "$checkTag already exist, skipping release."
  exit 0
fi

latestTag="$(git --no-pager tag -l | tail -1)"
changeLog="$(git --no-pager log --oneline "$latestTag"...HEAD)"

echo "New Version: $version"
echo "Latest Tag: $latestTag"
echo -e "Change Log Since Latest Tag: \n$changeLog"

hub release create -m "Azure Event Hubs Lib $version" -m "$changeLog" "$version"