package table

import (
	"gorm.io/datatypes"
)

type SuspendedTask struct {
	ID   int `gorm:"primaryKey"`
	Type SuspendedTaskType
	Info datatypes.JSON `gorm:"column:info;type:json"`
}

func (SuspendedTask) TableName() string {
	return "suspended_task"
}

type SuspendedTaskType int

const (
	Time SuspendedTaskType = iota
	Email
)

type SuspendedTimeInfo struct {
	Timestamp int64
}

type SuspendedEmailInfo struct {
	Email    string
	Keywords []string
}

func InitSuspendedTaskTable() error {
	err := DB.AutoMigrate(&SuspendedTask{})
	if err != nil {
		return err
	}
	return nil
}

func AddSuspendedTask(task SuspendedTask) int {
	err := DB.Create(&task).Error
	if err != nil {
		return -1
	}
	return task.ID
}
