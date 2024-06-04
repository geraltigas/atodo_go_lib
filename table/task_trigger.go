package table

import (
	"encoding/json"
	"errors"
	"gorm.io/datatypes"
)

type TaskTrigger struct {
	ID   int             `gorm:"primaryKey"`
	Type TaskTriggerType `gorm:"primaryKey"`
	Info datatypes.JSON  `gorm:"column:info"`
}

type TaskTriggerType int

const (
	Dependency TaskTriggerType = iota
	Event
)

func (t TaskTriggerType) String() (string, error) {
	names := [...]string{
		"Dependency",
		"Event",
	}
	if t < Dependency || t > Event {
		return "Unknown", errors.New("unknown TaskTriggerType")
	}
	return names[t], nil
}

func (TaskTrigger) TableName() string {
	return "task_trigger"
}

type DependencyInfo struct {
	Source int `json:"source"`
}

type EventInfo struct {
	EventName        string `json:"event_name"`
	EventDescription string `json:"event_description"`
}

func InitTaskTriggerTable() error {
	err := DB.AutoMigrate(&TaskTrigger{})
	if err != nil {
		return err
	}
	return nil
}

func AddOrUpdateTaskTrigger(taskTrigger TaskTrigger) error {
	err := DB.Create(&taskTrigger).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteTaskTriggersByID(id int) error {
	err := DB.Delete(&TaskTrigger{}, "id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}

func GetTaskTriggersByID(id int) ([]TaskTrigger, error) {
	var taskTriggers []TaskTrigger
	err := DB.Find(&taskTriggers, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return taskTriggers, nil
}

func (t TaskTrigger) SetDependencyInfo(source int) error {
	dependencyInfo := DependencyInfo{Source: source}
	marshal, err := json.Marshal(dependencyInfo)
	if err != nil {
		return err
	}
	t.Info = marshal
	return nil
}

func (t TaskTrigger) GetDependencyInfo() (*DependencyInfo, error) {
	var dependencyInfo DependencyInfo
	err := json.Unmarshal(t.Info, &dependencyInfo)
	if err != nil {
		return nil, err
	}
	return &dependencyInfo, nil
}

func (t TaskTrigger) SetEventInfo(eventName string, eventDescription string) error {
	eventInfo := EventInfo{EventName: eventName, EventDescription: eventDescription}
	marshal, err := json.Marshal(eventInfo)
	if err != nil {
		return err
	}
	t.Info = marshal
	return nil
}

func (t TaskTrigger) GetEventInfo() (*EventInfo, error) {
	var eventInfo EventInfo
	err := json.Unmarshal(t.Info, &eventInfo)
	if err != nil {
		return nil, err
	}
	return &eventInfo, nil
}
