package table

import (
	"testing"
	"time"
)

func TestAddTask(t *testing.T) {
	err := InitDB()
	if err != nil {
		t.Fatal(err)
	}
	err = InitTaskTable()
	task := Task{
		RootTask:   1,
		Name:       "Test Task",
		Goal:       "Test Goal",
		Deadline:   time.Now().UTC(),
		InWorkTime: false,
		Status:     Todo,
		ParentTask: 0,
	}
	id := AddTask(task)
	if id == -1 {
		t.Fatal("Failed to add task")
	}

	task.ID = id
	task2, err := GetTaskByID(id)
	if err != nil {
		t.Fatal(err)
	}

	if !task.Equal(task2) {
		t.Fatal("Failed to add task")
	}

	err = DeleteTask(id)
	if err != nil {
		t.Fatal(err)
	}
}
