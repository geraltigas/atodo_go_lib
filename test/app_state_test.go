package test

import "testing"

import "atodo_go/table"

func TestRootTask(t *testing.T) {
	err := table.InitDB()
	if err != nil {
		t.Fatal(err)
	}

	rootTask, err := table.GetRootTask()
	if err != nil {
		t.Fatal(err)
	}

	setRootTask := 100
	err = table.SetRootTask(setRootTask)
	if err != nil {
		t.Fatal(err)
	}

	rootTaskN, err := table.GetRootTask()
	if err != nil {
		t.Fatal(err)
	}

	if rootTaskN != setRootTask {
		t.Fatal("Failed to set root task")
	}

	err = table.SetRootTask(rootTask)
	if err != nil {
		t.Fatal(err)
	}
}
