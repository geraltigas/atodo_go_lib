package schedule

import (
	"atodo_go/table"
)

type TaskShow struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Goal       string `json:"goal"`
	Deadline   int64  `json:"deadline"`
	InWorkTime bool   `json:"in_work_time"`
}

type SuspendedInfo interface {
	GetType() string
}

type SuspendedTimeInfo struct {
	Time int64 `json:"time"`
}

type SuspendedEmailInfo struct {
	Email    string   `json:"email"`
	Keywords []string `json:"keywords"`
}

func (SuspendedTimeInfo) GetType() string {
	return "time"
}

func (SuspendedEmailInfo) GetType() string {
	return "email"
}

type SuspendedTaskShow struct {
	Id         int           `json:"id"`
	Name       string        `json:"name"`
	Goal       string        `json:"goal"`
	Deadline   int64         `json:"deadline"`
	InWorkTime bool          `json:"in_work_time"`
	Type       string        `json:"type"`
	Info       SuspendedInfo `json:"info"`
}

type EventTriggerTaskShow struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Goal             string `json:"goal"`
	Deadline         int64  `json:"deadline"`
	InWorkTime       bool   `json:"in_work_time"`
	EventName        string `json:"event_name"`
	EventDescription string `json:"event_description"`
}

type TSchedule struct {
	Tasks            []TaskShow             `json:"tasks"`
	SuspendedTasks   []SuspendedTaskShow    `json:"suspended_tasks"`
	EventTriggerTask []EventTriggerTaskShow `json:"event_trigger_tasks"`
}

func suspendedTaskPreprocess(task table.Task) {

}

// use generic type

func GetFirstElementFromSet[T comparable](set map[T]bool) *T {
	for key := range set {
		return &key
	}
	return nil
}

func Schedule() (*TSchedule, error) {
	tasks := make([]TaskShow, 0)
	suspendedTasks := make([]SuspendedTaskShow, 0)
	eventTriggerTasks := make([]EventTriggerTaskShow, 0)
	nowViewingTask, err := table.GetRootTask()
	if err != nil {
		return nil, err
	}
	waitForViewing := make(map[int]bool)
	waitForViewing[nowViewingTask] = true
	sourceTasks := make([]int, 0)
	subTasks := make([]int, 0)
	for len(waitForViewing) > 0 {
		taskId := *GetFirstElementFromSet(waitForViewing)
		task, err := table.GetTaskByID(taskId)
		if err != nil {
			return nil, err
		}
		if task.Status == table.Suspended {
			suspendedTaskPreprocess(task)
		}
		switch task.Status {
		case table.Suspended:
			suspendedTaskShow := SuspendedTaskShow{
				Id:         task.ID,
				Name:       task.Name,
				Goal:       task.Goal,
				Deadline:   task.Deadline.UnixMilli(),
				InWorkTime: task.InWorkTime,
			}
			suspendedTaskInfo, err := table.GetSuspendedTask(task.ID)
			if err != nil {
				return nil, err
			}
			switch suspendedTaskInfo.Type {
			case table.Time:
				timeInfo, err := suspendedTaskInfo.GetTimeInfo()
				if err != nil {
					return nil, err
				}
				suspendedTaskShow.Type = SuspendedTimeInfo{}.GetType()
				suspendedTaskShow.Info = SuspendedTimeInfo{
					Time: timeInfo.Timestamp,
				}
			case table.Email:
				emailInfo, err := suspendedTaskInfo.GetEmailInfo()
				if err != nil {
					return nil, err
				}
				suspendedTaskShow.Type = SuspendedEmailInfo{}.GetType()
				suspendedTaskShow.Info = SuspendedEmailInfo{
					Email:    emailInfo.Email,
					Keywords: emailInfo.Keywords,
				}
			}
			suspendedTasks = append(suspendedTasks, suspendedTaskShow)
			delete(waitForViewing, taskId)
			continue
		case table.Todo:
			sourceTasks = sourceTasks[:0]
			subTasks = subTasks[:0]
			sourceTasks, err = table.GetSourceTasks(task.ID)
			if err != nil {
				return nil, err
			}
			newTasks := make([]int, 0)
			for _, taskId := range sourceTasks {
				task, err := table.GetTaskByID(taskId)
				if err != nil {
					return nil, err
				}
				if task.Status != table.Done {
					newTasks = append(newTasks, taskId)
				}
			}
			sourceTasks = newTasks
			if len(sourceTasks) != 0 {
				for _, taskId := range sourceTasks {
					waitForViewing[taskId] = true
				}
				delete(waitForViewing, taskId)
				continue
			}

			if table.HaveSubTasks(taskId) {
				subTasks, err = table.GetSubTasksConnectedToEnd(taskId)
				for _, taskId := range subTasks {
					waitForViewing[taskId] = true
				}
				delete(waitForViewing, taskId)
				continue
			}

			eventTriggerTasks = eventTriggerTasks[:0]
			taskTriggers, err := table.GetTaskTriggersByID(taskId)
			if err != nil {
				return nil, err
			}
			if len(taskTriggers) == 1 {
				triggerInfo, err := taskTriggers[0].GetEventInfo()
				if err != nil {
					return nil, err
				}
				eventTriggerTasks = append(eventTriggerTasks, EventTriggerTaskShow{
					Id:               task.ID,
					Name:             task.Name,
					Goal:             task.Goal,
					Deadline:         task.Deadline.UnixMilli(),
					InWorkTime:       task.InWorkTime,
					EventName:        triggerInfo.EventName,
					EventDescription: triggerInfo.EventDescription,
				})
			} else {
				tasks = append(tasks, TaskShow{
					Id:         task.ID,
					Name:       task.Name,
					Goal:       task.Goal,
					Deadline:   task.Deadline.UnixMilli(),
					InWorkTime: task.InWorkTime,
				})
			}
			delete(waitForViewing, taskId)
			continue
		case table.Done:
			delete(waitForViewing, taskId)
			continue
		}
	}
	return &TSchedule{
		Tasks:            tasks,
		SuspendedTasks:   suspendedTasks,
		EventTriggerTask: eventTriggerTasks,
	}, nil
}
