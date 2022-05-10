# gostarter

a starter project for golang

## features

- [x] yaml 配置文件使用示例
- [ ] db 使用示例
- [x] golog 使用示例
- [x] [gin](https://github.com/gin-gonic/gin) 框架使用
- [x] pprof cpu 支持

    * `echo 5m > jj.cpu; kill -USR1 {pid}` 5分钟后，获取 cpu.profile 文件
    * `go tool pprof -http :9402 cpu.profile` 本机打开 web 界面查看

- [x] pprof mem 支持

    * `touch jj.mem; kill -USR1 ; kill -USR1 {pid}` 立即获取 mem.profile 文件
    * `go tool pprof -http :9402 mem.profile` 本机打开 web 界面查看

## build

1. `make install` for local
1. `make linux` for linux version
