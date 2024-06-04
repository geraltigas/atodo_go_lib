package test

import (
	"atodo_go/task_show"
	"testing"
)

func TestGetShowData(t *testing.T) {
	// test
	data, err := task_show.GetShowData()
	if err != nil {
		t.Fatal(err)
	}
	//if len(data.Nodes) != 0 {
	//	t.Fatal("len(data.Nodes) != 0")
	//}
	if len(data.Edges) != 0 {
		t.Fatal("len(data.Edges) != 0")
	}
}
