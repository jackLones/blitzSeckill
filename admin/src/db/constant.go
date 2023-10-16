package db

type TaskEnum int

const (
	TASK_STATUS_PENDING    TaskEnum = 1
	TASK_STATUS_PROCESSING TaskEnum = 2
	TASK_STATUS_SUCCESS    TaskEnum = 3
	TASK_STATUS_FAILED     TaskEnum = 4
)

const (
	// 优先1年
	MAX_PRIORITY = 3600 * 24 * 30 * 12
)

// IsValidStatus 任务状态合法性检验
func IsValidStatus(status TaskEnum) bool {
	if status == TASK_STATUS_PENDING {
		return true
	}
	if status == TASK_STATUS_PROCESSING {
		return true
	}
	if status == TASK_STATUS_SUCCESS {
		return true
	}
	if status == TASK_STATUS_FAILED {
		return true
	}
	return false
}
