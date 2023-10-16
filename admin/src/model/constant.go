package model

type ActivityEnum int

const (
	TASK_STATUS_PENDING    ActivityEnum = 1
	TASK_STATUS_PROCESSING ActivityEnum = 2
	TASK_STATUS_SUCCESS    ActivityEnum = 3
	TASK_STATUS_FAILED     ActivityEnum = 4
)

const (
	// 优先1年
	MAX_PRIORITY = 3600 * 24 * 30 * 12
)

// IsValidStatus 任务状态合法性检验
func IsValidStatus(status ActivityEnum) bool {
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
