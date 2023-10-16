package cache

import (
	"blitzSeckill/admin/src/config"
	"errors"
	"fmt"
	"github.com/niuniumart/gosdk/goredis"
	"github.com/niuniumart/gosdk/martlog"
	"golang.org/x/net/context"
	"time"
)

var (
	rdb        *goredis.RedisCli
	prefix     = "blitz_seckill_admin"
	expireTime = time.Hour * 24 // 默认为24小时

	storeDeductionScriptLua = "local c_s = redis.call('get', KEYS[1])\n" +
		"if not c_s or tonumber(c_s) < tonumber(KEYS[2]) then\n" +
		"return 0\n" +
		"end\n" +
		"redis.call('decrby',KEYS[1], KEYS[2])\n" +
		"return 1"

	storeDeductionScriptSha1 string
)

func InitCache() error {
	goredis.Factory.MaxIdleConn = config.Conf.Redis.MaxIdle
	goredis.Factory.IdleTimeout = time.Second * time.Duration(config.Conf.Redis.IdleTimeout)
	goredis.Factory.MaxConn = config.Conf.Redis.MaxActive

	redisCli, err := goredis.Factory.CreateRedisCli(config.Conf.Redis.Auth, config.Conf.Redis.Url)
	if err != nil {
		martlog.Errorf("Redis connection error: %v", err)
		return err
	}
	rdb = redisCli

	if config.Conf.Redis.CacheTimeoutDay != 0 {
		// 设置为config.Conf.Redis.CacheTimeoutDay 天的过期时间
		expireTime = expireTime * time.Duration(config.Conf.Redis.CacheTimeoutDay)
	}

	//在系统启动时，将脚本预加载到Redis中，并返回一个加密的字符串，下次只要传该加密窜，即可执行对应脚本，减少了Redis的预编译
	// 执行 SCRIPT LOAD 命令获取脚本的 SHA1 值
	sha1, err := rdb.RedisPool.ScriptLoad(context.Background(), storeDeductionScriptLua).Result()
	if err != nil {
		martlog.Errorf("Redis ScriptLoad Get SHA1 error: %v", err)
		return err
	}

	storeDeductionScriptSha1 = sha1
	return nil
}

// Evalsha
// 调用Lua脚本,不需要每次都传入Lua脚本，只需要传入预编译返回的sha1即可
func Evalsha(key string, buyNum int) int64 {
	// 使用 EvalSha 方法执行脚本
	result, err := rdb.RedisPool.EvalSha(context.Background(), storeDeductionScriptSha1, []string{key}, buyNum).Result()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	value, ok := result.(int64)
	if !ok {
		panic(errors.New("返回值类型错误"))
	}
	return value
}
