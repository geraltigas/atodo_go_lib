package test

import (
	"atodo_go/table"
	"testing"
)

func TestAddOrUpdateTaskAfterEffect(t *testing.T) {
	taskAfterEffect := table.TaskAfterEffect{
		ID:   1,
		Type: table.Periodic,
	}

	err := taskAfterEffect.SetPeriodicInfo(table.PeriodicT{
		Period:    1,
		NowAt:     2,
		Intervals: []int{3, 4},
	})
	if err != nil {
		return
	}

	err = table.AddOrUpdateTaskAfterEffect(taskAfterEffect)
	if err != nil {
		return
	}

	tae, err := table.GetTaskAfterEffect(taskAfterEffect.ID, taskAfterEffect.Type)
	if err != nil {
		t.Error(err)
	}

	if !tae.Equal(taskAfterEffect) {
		t.Error("TaskAfterEffect not equal")
	}
}
