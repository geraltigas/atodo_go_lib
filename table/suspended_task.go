package table

import (
	"encoding/json"
	"errors"
	"gorm.io/datatypes"
)

type SuspendedTask struct {
	ID   int `gorm:"primaryKey"`
	Type SuspendedTaskType
	Info datatypes.JSON `gorm:"column:info"`
}

func (st *SuspendedTask) SetTimeInfo(info SuspendedTimeInfo) error {
	infoBytes, err := json.Marshal(info)
	if err != nil {
		return err
	}
	st.Info = infoBytes
	return nil
}

func (st *SuspendedTask) SetEmailInfo(info SuspendedEmailInfo) error {
	infoBytes, err := json.Marshal(info)
	if err != nil {
		return err
	}
	st.Info = infoBytes
	return nil
}

func (st *SuspendedTask) GetTimeInfo() (*SuspendedTimeInfo, error) {
	info := SuspendedTimeInfo{}
	err := json.Unmarshal(st.Info, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (st *SuspendedTask) GetEmailInfo() (*SuspendedEmailInfo, error) {
	info := SuspendedEmailInfo{}
	err := json.Unmarshal(st.Info, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

type SuspendedTaskInfo interface {
	GetTimeInfo() (*SuspendedTimeInfo, error)
	GetEmailInfo() (*SuspendedEmailInfo, error)
}

func (*SuspendedTask) TableName() string {
	return "suspended_task"
}

func (st *SuspendedTask) Equal(other SuspendedTask) bool {
	if st.ID != other.ID {
		return false
	}
	if st.Type != other.Type {
		return false
	}
	if st.Type == Time {
		info1, _ := st.GetTimeInfo()
		info2, _ := other.GetTimeInfo()
		return info1.Equal(*info2)
	} else {
		info1, _ := st.GetEmailInfo()
		info2, _ := other.GetEmailInfo()
		return info1.Equal(*info2)
	}
}

type SuspendedTaskType int

const (
	Time SuspendedTaskType = iota
	Email
)

func (t SuspendedTaskType) String() (string, error) {
	names := [...]string{
		"Time",
		"Email",
	}
	if t < Time || t > Email {
		return "Unknown", errors.New("unknown SuspendedTaskType")
	}
	return names[t], nil
}

type SuspendedTimeInfo struct {
	Timestamp int64
}

func (info *SuspendedTimeInfo) GetTimeInfo() (*SuspendedTimeInfo, error) {
	return info, nil
}

func (info *SuspendedTimeInfo) GetEmailInfo() (*SuspendedEmailInfo, error) {
	return nil, errors.New("not an email info")
}

func (info *SuspendedTimeInfo) Equal(other SuspendedTimeInfo) bool {
	return info.Timestamp == other.Timestamp
}

type SuspendedEmailInfo struct {
	Email    string
	Keywords []string
}

func (info *SuspendedEmailInfo) GetTimeInfo() (*SuspendedTimeInfo, error) {
	return nil, errors.New("not a time info")
}

func (info *SuspendedEmailInfo) GetEmailInfo() (*SuspendedEmailInfo, error) {
	return info, nil
}

func (info *SuspendedEmailInfo) Equal(other SuspendedEmailInfo) bool {
	if info.Email != other.Email {
		return false
	}
	if len(info.Keywords) != len(other.Keywords) {
		return false
	}
	for i, keyword := range info.Keywords {
		if keyword != other.Keywords[i] {
			return false
		}
	}
	return true
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

func DeleteSuspendedTasks(id int) error {
	err := DB.Delete(&SuspendedTask{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func GetSuspendedTask(id int) (SuspendedTask, error) {
	var task SuspendedTask
	err := DB.First(&task, id).Error
	return task, err
}

func AddOrUpdateSuspendedTask(task SuspendedTask) error {
	err := DB.Save(&task).Error
	if err != nil {
		return err
	}
	return nil
}
