# go-starter
build a starter project for golang, including log, config, ctl, recover, statiq and etc.


## build

for release:

1. `./build.sh` for local
1. `./build.sh -t linux` for linux version
 
```bash
$ ./build.sh -h
Usage: ./build.sh [OPTION]...

  -t target   linux/local, default local
  -u yes/no   enable upx compression if upx is available or not
  -b          binary name, default notify4g
  -h          display help
```

for dev:

1. `go get github.com/bingoohuang/statiq`
1. `./buildres.sh`
1. `statiq -src=res`
1. `go fmt ./...;sh build.sh`