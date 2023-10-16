package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

var TaskTypeCfgNsp TaskScheduleCfg

/*
 * @desc: task schedule config
 * @TaskType：任务类型
 * @ScheduleLimit：每次取多少个任务来执行
 * @MaxProcessingTime：单个任务一次最大执行时间，单位秒
 * @MaxRetryNum:  短任务特有属性，最大重试次数，表示最大重试多少次
 * @MaxRetryInterval: 短任务特有属性，表示短任务重试间隔，本期暂时保留
 * @
 */

// TaskScheduleCfg cfg
type TaskScheduleCfg struct {
	TaskType          string
	ScheduleLimit     int
	ScheduleInterval  int
	MaxProcessingTime int64
	MaxRetryNum       int
	MaxRetryInterval  int
	CreateTime        *time.Time
	ModifyTime        *time.Time
}

// TableName 表名
func (p *TaskScheduleCfg) TableName() string {
	return "t_schedule_cfg"
}

// Create 创建记录
func (p *TaskScheduleCfg) Create(db *gorm.DB, task *TaskScheduleCfg) error {
	err := db.Table(p.TableName()).Create(task).Error
	return err
}

// Save 保存记录
func (p *TaskScheduleCfg) Save(db *gorm.DB, task *TaskScheduleCfg) error {
	err := db.Table(p.TableName()).Save(task).Error
	return err
}

// GetTaskTypeCfg 获取记录
func (p *TaskScheduleCfg) GetTaskTypeCfg(db *gorm.DB, taskType string) (*TaskScheduleCfg, error) {
	var cfg = new(TaskScheduleCfg)
	err := db.Table(p.TableName()).Where("task_type = ?", taskType).First(&cfg).Error
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// GetTaskTypeCfgList 获取记录列表
func (p *TaskScheduleCfg) GetTaskTypeCfgList(db *gorm.DB) ([]*TaskScheduleCfg, error) {
	var taskTypeCfgList = make([]*TaskScheduleCfg, 0)
	db = db.Table(p.TableName())
	err := db.Find(&taskTypeCfgList).Error
	if err != nil {
		return nil, err
	}
	return taskTypeCfgList, nil
}
