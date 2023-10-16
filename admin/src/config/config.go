package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var Conf *TomlConfig

// TomlConfig 配置
type TomlConfig struct {
	Common commonConfig
	MySQL  mysqlConfig
	Redis  redisConfig
	Task   TaskConfig
}

type commonConfig struct {
	Port    int  `toml:"port"`
	OpenTLS bool `toml:"open_tls"`
}

type mysqlConfig struct {
	Url         string `toml:"url"`
	User        string `toml:"user"`
	Pwd         string `toml:"pwd"`
	Dbname      string `toml:"db_name"`
	MaxIdleConn int    `toml:"max_idle"`
	MaxConn     int    `toml:"max_active"`
	IdleTimeout int    `toml:"idle_timeout"`
}

type redisConfig struct {
	Url         string `toml:"url"`
	Auth        string `toml:"auth"`
	MaxIdle     int    `toml:"max_idle"`
	MaxActive   int    `toml:"max_active"`
	IdleTimeout int    `toml:"idle_timeout"`
	// CacheTimeout           int    `toml:"cache_timeout"`
	// CacheTimeoutVerifyCode int    `toml:"cache_timeout_verify_code"`
	CacheTimeoutDay int `toml:"cache_timeout_day"`
}

type TaskConfig struct {
	TableMaxRows        int   `toml:"table_max_rows"`        // 表最大行数
	AliveThreshold      int   `toml:"alive_threshold"`       // 任务存活阈值
	SplitInterval       int   `toml:"split_interval"`        // 分表间隔
	LongProcessInterval int   `toml:"long_process_interval"` // 长时间处理间隔
	MoveInterval        int   `toml:"move_interval"`         // 更新begin下标的时间间隔
	MaxProcessTime      int64 `toml:"max_process_time"`      // 最大处理时间
}

// LoadConfig 导入配置
func (c *TomlConfig) LoadConfig() {

	if _, err := os.Stat(GetConfigPath()); err != nil {
		panic(err)
	}

	if _, err := toml.DecodeFile(GetConfigPath(), &c); err != nil {
		panic(err)
	}
}

func Init() {
	// 初始化配置
	initConf()
}

// InitConf 初始化配置
func initConf() {
	Conf = new(TomlConfig)
	Conf.LoadConfig()
}

// 项目主目录
var rootDir string

func GetConfigPath() string {
	return rootDir + "/config/config.toml"
}

func init() {
	inferRootDir()
	// 初始化配置
}

// 推断 Root目录（copy就行）
func inferRootDir() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var infer func(string) string
	infer = func(dir string) string {
		if exists(dir + "/main.go") {
			return dir
		}

		// 查看dir的父目录
		parent := filepath.Dir(dir)
		return infer(parent)
	}

	rootDir = infer(pwd)
}

func exists(dir string) bool {
	// 查找主机是不是存在 dir
	_, err := os.Stat(dir)
	return err == nil || os.IsExist(err)
}
