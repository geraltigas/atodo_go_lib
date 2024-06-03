package test

import (
	"atodo_go/table"
	"testing"
	"time"
)

func TestAddTask(t *testing.T) {
	err := table.InitDB()
	if err != nil {
		t.Fatal(err)
	}
	err = table.InitTaskTable()
	task := table.Task{
		RootTask:   1,
		Name:       "Test Task",
		Goal:       "Test Goal",
		Deadline:   time.Now().UTC(),
		InWorkTime: false,
		Status:     table.Todo,
		ParentTask: 0,
	}
	id := table.AddTask(task)
	if id == -1 {
		t.Fatal("Failed to add task")
	}

	task.ID = id
	task2, err := table.GetTaskByID(id)
	if err != nil {
		t.Fatal(err)
	}

	if !task.Equal(task2) {
		t.Fatal("Failed to add task")
	}

	err = table.DeleteTask(id)
	if err != nil {
		t.Fatal(err)
	}
}
