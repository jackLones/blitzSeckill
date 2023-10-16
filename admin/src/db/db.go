package db

import (
	"blitzSeckill/admin/src/config"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/niuniumart/gosdk/gormcli"
)

var DB *gorm.DB

// InitDB 初始化DB
func InitDB() error {
	var err error
	gormcli.Factory = gormcli.GormFactory{
		MaxIdleConn: config.Conf.MySQL.MaxIdleConn,
		MaxConn:     config.Conf.MySQL.MaxConn,
		IdleTimeout: config.Conf.MySQL.IdleTimeout,
	}

	DB, err = gormcli.Factory.CreateGorm(config.Conf.MySQL.User,
		config.Conf.MySQL.Pwd, config.Conf.MySQL.Url, config.Conf.MySQL.Dbname)
	if err != nil {
		return err
	}

	// 尝试发送 ping 请求
	err = DB.DB().Ping()
	if err != nil {
		fmt.Println("Database connection is not available:", err)
		return err
	}

	fmt.Println("Database connection is available")

	return nil
}

const (
	GORM_DUPLICATE_ERR_KEY = "Duplicate entry"
)

// IsDupErr 重复记录错误判定
func IsDupErr(err error) bool {
	return strings.Contains(err.Error(), GORM_DUPLICATE_ERR_KEY)
}

// GetTaskTableName 获取short task 表名
func GetTaskTableName(taskType string) string {
	taskTableName := fmt.Sprintf("t_%s_task", taskType)
	return taskTableName
}
