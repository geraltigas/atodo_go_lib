package table

import "testing"

func TestRootTask(t *testing.T) {
	err := InitDB()
	if err != nil {
		t.Fatal(err)
	}

	rootTask, err := GetRootTask()
	if err != nil {
		t.Fatal(err)
	}

	setRootTask := 100
	err = SetRootTask(setRootTask)
	if err != nil {
		t.Fatal(err)
	}

	rootTaskN, err := GetRootTask()
	if err != nil {
		t.Fatal(err)
	}

	if rootTaskN != setRootTask {
		t.Fatal("Failed to set root task")
	}

	err = SetRootTask(rootTask)
	if err != nil {
		t.Fatal(err)
	}
}
