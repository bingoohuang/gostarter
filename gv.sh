#!/usr/bin/env bash

# https://stackoverflow.com/questions/10376206/what-is-the-preferred-bash-shebang
# from https://github.com/CodisLabs/codis/blob/release3.2/version

if ! version=$(git log --date=iso --pretty=format:"%cd @%H" -1); then
  version="unknown version"
fi

if ! compile=$(date +"%F %T %z")" by "$(go version); then
  compile="unknown datetime"
fi

if ! describe=$(git describe --tags 2>/dev/null); then
  version="${version} @${describe}"
fi

mkdir -p util
cat <<EOF | gofmt >util/v.go
package util

const (
    Version = "$version"
    Compile = "$compile"
)
EOF
