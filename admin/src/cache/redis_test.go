package cache

import (
	"blitzSeckill/admin/src/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestInitDB(t *testing.T) {
	config.Init()
	err := InitCache()
	if err != nil {

	}
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456",
		DB:       0,
	})

	// 执行 SCRIPT LOAD 命令获取脚本的 SHA1 值
	sha1, err := client.ScriptLoad(context.Background(), "return 'Hello World!'").Result()
	if err != nil {
		panic(err)
	}

	// 使用 EvalSha 方法执行脚本
	result, err := client.EvalSha(context.Background(), sha1, nil).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(result) // 输出 Hello World!

}
