#!/bin/bash

set -x #echo on

target=local
upx=yes
bin=`basename "$PWD"`

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

echo bin:${bin}
echo target:${target}
echo upx:${upx}

# notice how we avoid spaces in $now to avoid quotation hell in go build
now=$(date +'%Y-%m-%d_%T')

if [[ ${target} = "linux" ]]; then
    export GOOS=linux
    export GOARCH=amd64
    bin=${bin}_linux_amd64
fi

go fmt ./...
go build -ldflags "-w -s -X main.sha1ver=`git rev-parse HEAD` -X main.builtTime=$now" -o "${bin}"
if [[ ${upx} = "yes" ]] && type upx > /dev/null 2>&1; then
    upx ${bin}
fi
