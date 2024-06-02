package table

import (
	"testing"
	"time"
)

func TestAddTask(t *testing.T) {
	err := InitDB()
	err = InitTaskTable()
	if err != nil {
		return
	}
	task := Task{
		Name:       "Test Task",
		Goal:       "Test Goal",
		RootTask:   0,
		Deadline:   time.Now().Unix(),
		InWorkTime: true,
		Status:     Todo,
		ParentTask: 0,
	}

	err = AddTask(task)
	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}

}
