#!/usr/bin/env bash

set -x #echo on

target=local
upx=yes
bin=$(basename "$PWD")

function usage {
	cat <<EOM
Usage: $0 [OPTION]...
  -t target   linux/local, default local
  -u yes/no   enable upx compression if upx is available or not
  -b          binary name, default ${bin}
  -h          display help
EOM
}

while getopts "t:b:u:h-:" optKey; do
  case ${optKey} in
    t) target=$OPTARG ;;
    u) upx=$OPTARG ;;
    b) bin=$OPTARG ;;
    h|*) usage; exit 0;;
    esac
done

echo bin:"${bin}"
echo target:"${target}"
echo upx:"${upx}"

if [[ ${target} = "linux" ]]; then
    export GOOS=linux
    export GOARCH=amd64
fi

./gv.sh

go build -o "${bin}"
if [[ ${upx} = "yes" ]] && type upx > /dev/null 2>&1; then
    upx "${bin}"
fi
