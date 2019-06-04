# go-starter

<del>build a starter project for golang, including log, config, ctl, recover, statiq and etc.</del>

## functions

1. viper 聚合配置使用
1. logrus 日志使用
1. res 资源内嵌
1. gin 框架使用
1. pprof 支持

    * `./bin --pprof-addr localhost:6060`
    * open `http://localhost:6060/debug/pprof` in explorer
    * or 可视化数据（火焰图），见如下：
    * `curl http://localhost:6060/debug/pprof/heap > heap.prof`
    * `go get -u github.com/google/pprof`
    * `pprof -http=:8080heap.prof`

## build

for release:

1. `./build.sh` for local
1. `./build.sh -t linux` for linux version
 
```bash
$ ./build.sh -h
Usage: ./build.sh [OPTION]...

  -t target   linux/local, default local
  -u yes/no   enable upx compression if upx is available or not
  -b          binary name, default go-starter
  -h          display help
```

for dev:

1. `go get github.com/bingoohuang/statiq`
1. `./buildres.sh`
1. `statiq -src=res`
1. `go fmt ./...;sh build.sh`

## start

run `./go-starter -o=false -u`.

