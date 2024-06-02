package table

import (
	"encoding/json"
	"testing"
)

func TestAddSuspendedTask(t *testing.T) {
	err := InitDB()
	if err != nil {
		return
	}
	var task SuspendedTask
	task.ID = 1
	task.Type = Time

	timeInfo := SuspendedTimeInfo{Timestamp: 1234567890}
	jsonData, err := json.Marshal(timeInfo)
	if err != nil {
		t.Error("Failed to marshal timeInfo")
	}

	task.Info = jsonData

	if AddSuspendedTask(task) != 1 {
		t.Error("Failed to add suspended task")
	}

	var task2 SuspendedTask
	task2.ID = 2
	task2.Type = Email

	emailInfo := SuspendedEmailInfo{
		Email:    "test",
		Keywords: []string{"test1", "test2"},
	}

	jsonData, err = json.Marshal(emailInfo)
	if err != nil {
		t.Error("Failed to marshal emailInfo")
	}

	task2.Info = jsonData

	if AddSuspendedTask(task2) != 2 {
		t.Error("Failed to add suspended task")
	}
}
