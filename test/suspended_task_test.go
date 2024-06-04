package test

import (
	"atodo_go/table"
	"testing"
)

func TestAddSuspendedTask(t *testing.T) {
	err := table.InitDB()
	if err != nil {
		return
	}
	var task table.SuspendedTask
	task.ID = 1
	task.Type = table.Time
	err = task.SetTimeInfo(table.SuspendedTimeInfo{Timestamp: 1234567890})
	if err != nil {
		return
	}

	if table.AddSuspendedTask(task) != 1 {
		t.Error("Failed to add suspended task")
	}

	suspendedTask, err := table.GetSuspendedTask(1)
	if err != nil {
		return
	}

	if !task.Equal(suspendedTask) {
		t.Error("Failed to get suspended task")
	}

	task.ID = 2
	task.Type = table.Email
	err = task.SetEmailInfo(table.SuspendedEmailInfo{Email: "test", Keywords: []string{"test"}})
	if err != nil {
		return
	}

	if table.AddSuspendedTask(task) != 2 {
		t.Error("Failed to add suspended task")
	}

	suspendedTask, err = table.GetSuspendedTask(2)
	if err != nil {
		return
	}

	if !task.Equal(suspendedTask) {
		t.Error("Failed to get suspended task")
	}

	err = table.DeleteSuspendedTasks(1)
	if err != nil {
		return
	}

	err = table.DeleteSuspendedTasks(2)
	if err != nil {
		return
	}
}
