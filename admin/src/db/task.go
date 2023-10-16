package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/niuniumart/gosdk/martlog"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var TaskNsp Task

const (
	BEGIN_TABLE_POS = 1
)

// Task 任务表
type Task struct {
	Id               int       `gorm:"primary_key;AUTO_INCREMENT:15"`
	UserId           string    `gorm:"type:varchar(256);column:user_id;not null;default:'';index:idx_user_id"`
	TaskId           string    `gorm:"type:varchar(256);column:task_id;not null;default:'';unique_index:idx_task_id"`
	TaskType         string    `gorm:"type:varchar(128);column:task_type;not null"`
	TaskStage        string    `gorm:"type:varchar(128);column:task_stage;not null"`
	Status           int       `gorm:"type:tinyint(3);column:status;not null;default:1;index:idx_tasktype_status_modify_time"`
	Priority         int       `gorm:"type:int(11);column:priority;not null;default:0;comment:'优先级'"`
	CrtRetryNum      int       `gorm:"type:int(11);column:crt_retry_num;not null;default:0;comment:'已经重试几次了'"`
	MaxRetryNum      int       `gorm:"type:int(11);column:max_retry_num;not null;default:0;comment:'最大能重试几次'"`
	MaxRetryInterval int       `gorm:"type:int(11);column:max_retry_interval;not null;default:0;comment:'最大重试间隔'"`
	ScheduleLog      string    `gorm:"type:varchar(4096);column:schedule_log;not null;default:'';comment:'调度信息记录'"`
	TaskContext      string    `gorm:"type:varchar(8192);column:task_context;not null;default:'';comment:'任务上下文，用户自定义'"`
	OrderTime        int64     `gorm:"type:int(20);column:order_time;not null;default:0;comment:'调度时间，越小调度越优先';index:idx_tasktype_status_modify_time"`
	CreateTime       time.Time `gorm:"type:datetime;column:create_time;not null;autoCreateTime;"`
	ModifyTime       time.Time `gorm:"type:datetime;column:modify_time;not null;autoCreateTime;"`
}

func (p *Task) getTableName(taskType, pos string) string {
	return fmt.Sprintf("t_%s_%s_%s", taskType, p.TableName(), pos)
}

// TableName 表名
func (p *Task) TableName() string {
	return "task"
}

func (p *Task) GetTaskTypeFromPosTaskType(posTaskType string) string {
	// t_lark_task
	s := strings.Split(posTaskType, "_")
	s[1] = strings.Replace(s[1], "-", "_", -1)
	return s[1]
}

// GenTaskId 生成对应taskId
func (p *Task) GenTaskId(taskType, pos string) string {
	taskType = strings.Replace(taskType, "_", "-", -1)
	return fmt.Sprintf("%+v_%s_%s", uuid.New(), taskType, pos)
}
func (p *Task) getTablePosFromTaskId(taskId string) (string, string) {
	s := strings.Split(taskId, "_")
	switch len(s) {
	case 3:
		s[1] = strings.Replace(s[1], "-", "_", -1)
		return s[1], s[2]
	default:
		martlog.Errorf("big error taskId %s have not match _", taskId)
		return "", ""
	}
}

// BatchSetStatus batch set
func (p *Task) BatchSetStatus(db *gorm.DB,
	taskIdList []string, status TaskEnum) error {
	var dic = map[string]interface{}{
		"status":      status,
		"modify_time": time.Now().Format("2006-01-02 15:04:05"),
	}
	tmpTaskId := taskIdList[0]
	taskType, pos := p.getTablePosFromTaskId(tmpTaskId)
	db = db.Table(p.getTableName(taskType, pos)).Where("task_id in (?)", taskIdList).
		UpdateColumns(dic)
	err := db.Error
	if err != nil {
		return err
	}
	return nil
}

// CreateTable 创建任务信息表（需要任务类型，已经表的position）
func (p *Task) CreateTable(db *gorm.DB, taskType, pos string) error {
	// 构造表名
	newTableName := p.getTableName(taskType, pos)
	return db.Table(newTableName).AutoMigrate(&Task{}).Error
}

// Find 查找记录
func (p *Task) Find(db *gorm.DB, taskId string) (*Task, error) {
	var data = &Task{}
	taskType, pos := p.getTablePosFromTaskId(taskId)
	err := db.Table(p.getTableName(taskType, pos)).Where("task_id = ?", taskId).First(data).Error
	return data, err
}

// Create 创建记录
func (p *Task) Create(db *gorm.DB, taskType, pos string, task *Task) error {
	err := db.Table(p.getTableName(taskType, pos)).Create(task).Error
	return err
}

// Save 保存记录
func (p *Task) Save(db *gorm.DB, task *Task) error {
	taskType, pos := p.getTablePosFromTaskId(task.TaskId)
	err := db.Table(p.getTableName(taskType, pos)).Save(task).Error
	return err
}

// GetTaskList 获取记录列表
func (p *Task) GetTaskList(db *gorm.DB,
	pos string, taskType string, status TaskEnum, limit int) ([]*Task, error) {
	var taskList = make([]*Task, 0)
	err := db.
		Table(p.getTableName(taskType, pos)).
		Where("status = ?", status).
		Order("order_time").
		Limit(limit).
		Find(&taskList).Error
	if err != nil {
		return nil, err
	}
	return taskList, nil
}

// GetAliveTaskList 获取处于激活状态的任务列表
func (p *Task) GetAliveTaskList(db *gorm.DB, taskType, pos string, limit int) ([]*Task, error) {
	var taskList = make([]*Task, 0)
	var statusSet = []TaskEnum{TASK_STATUS_PENDING, TASK_STATUS_PROCESSING}
	err := db.
		Table(p.getTableName(taskType, pos)).
		Order("modify_time").
		Limit(limit).
		Where("status in (?)", statusSet).
		Find(&taskList).Error
	if err != nil {
		return nil, err
	}
	return taskList, nil
}

// GetAliveTaskCount 获取处于激活状态的任务数
func (p *Task) GetAliveTaskCount(db *gorm.DB, taskType, pos string) (int, error) {
	return p.getTaskCount(db, taskType, pos,
		[]TaskEnum{TASK_STATUS_PENDING, TASK_STATUS_PROCESSING})
}

// GetAllTaskCount 获取所有任务数
func (p *Task) GetAllTaskCount(db *gorm.DB, taskType, pos string) (int, error) {
	return p.GetTableCount(db, taskType, pos)
}

// GetTaskCountByStatus 根据状态获取任务数
func (p *Task) GetTaskCountByStatus(db *gorm.DB, taskType, pos string, status int) (int, error) {
	var count int
	err := db.Table(p.getTableName(taskType, pos)).Where("status = ?", status).Count(&count).Error
	if err != nil {
		return count, err
	}
	return count, nil
}

// GetFinishTaskCount 获取处于完成状态的任务数
func (p *Task) GetFinishTaskCount(db *gorm.DB, taskType, pos string) (int, error) {
	// 任务失败和成功都算完成
	return p.getTaskCount(db, taskType, pos,
		[]TaskEnum{TASK_STATUS_FAILED, TASK_STATUS_SUCCESS})

}

func (p *Task) getTaskCount(db *gorm.DB, taskType, pos string, statusSet []TaskEnum) (int, error) {
	var count int
	err := db.Table(p.getTableName(taskType, pos)).Where("status in (?)", statusSet).Count(&count).Error
	if err != nil {
		return count, err
	}
	return count, nil
}

// 获取表记录总数
func (p *Task) GetTableCount(db *gorm.DB, taskType, pos string) (int, error) {
	var count int
	err := db.Table(p.getTableName(taskType, pos)).Count(&count).Error
	if err != nil {
		return count, err
	}
	return count, nil
}

// SetStatusPending 设置任务为等待状态
func (p *Task) SetStatusPending(db *gorm.DB, taskId string) error {
	return p.SetStatus(db, taskId, TASK_STATUS_PENDING)
}

// SetStatusSucc 设置任务为成功状态
func (p *Task) SetStatusSucc(db *gorm.DB, taskId string) error {
	return p.SetStatus(db, taskId, TASK_STATUS_SUCCESS)
}

// SetStatusFailed 设置任务为失败状态
func (p *Task) SetStatusFailed(db *gorm.DB, taskId string) error {
	return p.SetStatus(db, taskId, TASK_STATUS_FAILED)
}

// SetStatus 设置任务对应状态
func (p *Task) SetStatus(db *gorm.DB, taskId string, status TaskEnum) error {
	var dic = map[string]interface{}{
		"status": status,
	}
	taskType, pos := p.getTablePosFromTaskId(taskId)
	err := db.Table(p.getTableName(taskType, pos)).Where("task_id = ?", taskId).Updates(dic).Error
	if err != nil {
		return err
	}
	return nil
}

// SetStatusAndRetryNumIncrement 设置任务对应状态，并且重试次数加1
func (p *Task) SetStatusAndRetryNumIncrement(db *gorm.DB, taskId string, status TaskEnum) error {
	var dic = map[string]interface{}{
		"status":        status,
		"crt_retry_num": p.CrtRetryNum + 1,
	}
	taskType, pos := p.getTablePosFromTaskId(taskId)
	err := db.Table(p.getTableName(taskType, pos)).Where("task_id = ?", taskId).Updates(dic).Error
	if err != nil {
		return err
	}
	return nil
}

// SetStatusWithOutModifyTime 设置对应任务状态，不更新modify_time
func (p *Task) SetStatusWithOutModifyTime(db *gorm.DB, taskId string, status TaskEnum) error {
	taskType, pos := p.getTablePosFromTaskId(taskId)
	err := db.Table(p.getTableName(taskType, pos)).Where("task_id = ?", taskId).UpdateColumn("status", status).Error
	if err != nil {
		return err
	}
	return nil
}

// SetContext 设置任务的上下文
func (p *Task) SetContext(db *gorm.DB, taskId, context string) error {
	var dic = map[string]interface{}{
		"task_context": context,
	}
	taskType, pos := p.getTablePosFromTaskId(taskId)
	err := db.Table(p.getTableName(taskType, pos)).Where("task_id = ?", taskId).Updates(dic).Error
	if err != nil {
		return err
	}
	return nil
}

// GetLongTimeProcessing 获取超过最大执行时间的任务列表
func (p *Task) GetLongTimeProcessing(db *gorm.DB,
	taskType, pos string, maxProcessTime int64, limit int) ([]*Task, error) {
	var Tasks = make([]*Task, 0)
	err := db.Table(p.getTableName(taskType, pos)).
		Where("status = ?", TASK_STATUS_PROCESSING).
		Where("unix_timestamp(modify_time) + ? < ?", maxProcessTime, time.Now().Unix()).
		Limit(limit).
		Find(&Tasks).
		Error
	if err != nil {
		return nil, err
	}
	return Tasks, nil
}

// ModifyTimeoutPending 更新超过最大执行时间的任务为等待状态
func (p *Task) ModifyTimeoutPending(db *gorm.DB, taskType, pos string, maxProcessTime int64) error {
	var dic = map[string]interface{}{
		"status": TASK_STATUS_PENDING,
	}
	err := db.Table(p.getTableName(taskType, pos)).
		Where("status = ?", TASK_STATUS_PROCESSING).
		Where("unix_timestamp(modify_time) + ? < ?", maxProcessTime, time.Now().Unix()).
		Updates(dic).
		Error
	if err != nil {
		return err
	}
	return nil
}

// IncreaseCrtRetryNum 增加一次对应任务的重试次数
func (p *Task) IncreaseCrtRetryNum(db *gorm.DB, taskId string) error {
	taskType, pos := p.getTablePosFromTaskId(taskId)
	return db.Table(p.getTableName(taskType, pos)).
		Where("task_id = ?", taskId).
		Update("crt_retry_num", gorm.Expr("crt_retry_num + ?", 1)).Error
}

// BeforeCreate 创建之前的回调函数
func (p *Task) BeforeCreate(scope *gorm.Scope) error {
	now := time.Now()
	scope.SetColumn("create_time", now)
	scope.SetColumn("modify_time", now)
	return nil
}

// UpdateTask 更新任务
func (p *Task) UpdateTask(db *gorm.DB) error {
	taskType, pos := p.getTablePosFromTaskId(p.TaskId)
	tableName := p.getTableName(taskType, pos)
	p.ModifyTime = time.Now()
	err := db.Table(tableName).Where("task_id = ?", p.TaskId).
		Where("status <> ? and status <> ?", TASK_STATUS_SUCCESS, TASK_STATUS_FAILED).Updates(p).Error
	if err != nil {
		return err
	}
	return nil
}

// SetScheduleLog 设置任务的调度信息
func (p *Task) SetScheduleLog(db *gorm.DB, ScheduleLog string) error {
	p.ScheduleLog = ScheduleLog
	return p.UpdateTask(db)
}

func (p *Task) BatchSetOwnerStatusWithPendingOutModify(db *gorm.DB,
	taskIdList []string, owner string, status TaskEnum) (int64, error) {
	var dic = map[string]interface{}{
		"status": status,
	}
	if owner != "" {
		dic["owner"] = owner
	}
	tmpTaskId := taskIdList[0]
	taskType, pos := p.getTablePosFromTaskId(tmpTaskId)
	db = db.Table(p.getTableName(taskType, pos)).Where("task_id in (?)", taskIdList).
		Where("status = ?", TASK_STATUS_PENDING).UpdateColumns(dic)
	err := db.Error
	if err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

// GetAssignTasksByOwnerStatus 在指定任务列表中获取对应归宿和状态的任务列表
func (p *Task) GetAssignTasksByOwnerStatus(db *gorm.DB,
	taskIdList []string, owner string, status TaskEnum, limit int64) ([]*Task, error) {
	if len(taskIdList) == 0 {
		martlog.Infof("taskId list is empty")
		return nil, nil
	}
	var Tasks = make([]*Task, 0)
	tmpTaskId := taskIdList[0]
	taskType, pos := p.getTablePosFromTaskId(tmpTaskId)
	err := db.Table(p.getTableName(taskType, pos)).
		Where("task_id in (?)", taskIdList).
		Where("owner = ? and status = ?", owner, status).
		Limit(limit).
		Find(&Tasks).
		Error
	if err != nil {
		return nil, err
	}
	return Tasks, nil
}

// ConventTaskIdList 任务信息列表转换成对应任务ID列表
func ConventTaskIdList(tasks []*Task) []string {
	taskIds := make([]string, 0, len(tasks))
	for _, task := range tasks {
		if task != nil {
			taskIds = append(taskIds, task.TaskId)
		}
	}
	return taskIds
}
