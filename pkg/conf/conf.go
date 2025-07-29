package conf

import (
	"github.com/bingoohuang/gg/pkg/fla9"
	"time"
)

type DbConf struct {
	Driver             string
	DataSource         string
	MaxOpenConnections int
	MaxIdleConnections int
	ConnMaxLifetime    int
	ConnMaxIdleTime    int
}

// conf Command line options
type conf struct {
	// 下面 3 项 为建议固定设置，不需要改动
	Conf    string `flag:"conf" usage:"yaml 配置文件路径" yaml:"-"`
	Init    bool   `usage:"创建样例配置文件和 ctl 文件，然后退出"  yaml:"-"`
	Version bool   `flag:",V" usage:"打印版本号"  yaml:"-"`

	// 以下配置为示例配置（exportable)，根据实际业务需要进行调整
	// 其中默认值只有在 conf.yml 文件中响应参数没必要定义时才生效
	Level  string   `val:"all" usage:"Output level"` // 也可以通过环境变量进行设置（例如：export GG_LEVEL=INFO)
	Output []string // 演示切片行参数
	Log    struct {
		Spec   string
		Layout string
	}
	Db       DbConf
	Duration time.Duration `flag:"d"`                                             // 命令行参数 -d，演示时长参数类型
	MyFlag   MyFlag        `flag:"my"`                                            // 命令行参数 -my，演示自定义解析类型参数
	Port     int           `flag:"p" val:"1234"`                                  // 命令行参数 -p 默认值 1234
	Other    string        `flag:"-"`                                             // 忽略
	V        int           `count:"true"`                                         // 命令行上，支持 -v 或者 -vv 等计数形式的参数
	Size     uint64        `flag:",s" size:"true" val:"10MiB" yaml:",label=size"` // 字节大小，短名 s
	Pmem     float32

	IgnorePaths []string
}

// MyFlag 演示自定义 flag 的使用
type MyFlag struct {
	Value string
}

// 确保 MyFlag 实现了 flag.Value 接口
var _ fla9.Value = (*MyFlag)(nil)

func (i *MyFlag) String() string { return i.Value }

func (i *MyFlag) Set(value string) error {
	i.Value = value
	return nil
}

var Conf = &conf{}
