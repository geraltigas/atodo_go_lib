package table

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"
)

type TaskStatus int

const (
	Todo TaskStatus = iota
	Suspended
	Done
)

func (status *TaskStatus) String() (string, error) {
	names := [...]string{
		"Todo",
		"Suspended",
		"Done",
	}
	if *status < Todo || *status > Done {
		return "Unknown", fmt.Errorf("unknown TaskStatus")
	}
	return names[*status], nil
}

func (status *TaskStatus) FromString(status2 string) {
	switch status2 {
	case "Todo":
		*status = Todo
	case "Suspended":
		*status = Suspended
	case "Done":
		*status = Done
	default:
		*status = Todo
	}
}

type Task struct {
	ID                   int       `gorm:"primaryKey;autoIncrement"`
	RootTask             int       `gorm:"column:root_task"`
	Name                 string    `gorm:"type:text"`
	Goal                 string    `gorm:"type:text"`
	Deadline             time.Time `gorm:"type:timestamp"`
	InWorkTime           bool      `gorm:"column:in_work_time"`
	Status               TaskStatus
	ParentTask           int    `gorm:"column:parent_task"`
	PositionX            int    `gorm:"column:position_x"`
	PositionY            int    `gorm:"column:position_y"`
	DependencyConstraint string `gorm:"column:dependency_constraint"`
	SubtaskConstraint    string `gorm:"column:subtask_constraint"`
}

func (Task) TableName() string {
	return "task"
}

func (task Task) Equal(other Task) bool {
	return task.ID == other.ID &&
		task.RootTask == other.RootTask &&
		task.Name == other.Name &&
		task.Goal == other.Goal &&
		task.Deadline == other.Deadline &&
		task.InWorkTime == other.InWorkTime &&
		task.Status == other.Status &&
		task.ParentTask == other.ParentTask
}

func InitTaskTable() error {
	err := DB.AutoMigrate(&Task{})
	if err != nil {
		log.Fatal("Failed to migrate Task table: ", err)
		return err
	}
	log.Println("Task table migrated")
	return nil
}

func AddTask(task Task) int {
	err := DB.Create(&task).Error
	if err != nil {
		log.Fatal("Failed to add task: ", err)
		return -1
	}
	log.Println("Task added: ", task.ID)
	return task.ID
}

func DeleteTask(id int) error {
	err := DB.Delete(&Task{}, id).Error
	if err != nil {
		log.Fatal("Failed to delete task: ", err)
		return err
	}
	log.Println("Task deleted: ", id)
	return nil
}

func ClearAllTasks() error {
	// get table name from model
	err := DB.Exec("DELETE FROM " + (&Task{}).TableName()).Error
	if err != nil {
		log.Fatal("Failed to delete all tasks: ", err)
		return err
	}
	log.Println("All tasks deleted")
	return nil
}

func UpdateTaskName(id int, name string) error {
	err := DB.Model(&Task{}).Where("id = ?", id).Update("name", name).Error
	if err != nil {
		log.Fatal("Failed to update task name: ", err)
		return err
	}
	log.Println("Task name updated: ", id, name)
	return nil
}

func UpdateTaskGoal(id int, goal string) error {
	err := DB.Model(&Task{}).Where("id = ?", id).Update("goal", goal).Error
	if err != nil {
		log.Fatal("Failed to update task goal: ", err)
		return err
	}
	log.Println("Task goal updated: ", id, goal)
	return nil
}

func UpdateTaskDeadline(id int, deadline int64) error {
	err := DB.Model(&Task{}).Where("id = ?", id).Update("deadline", deadline).Error
	if err != nil {
		log.Fatal("Failed to update task deadline: ", err)
		return err
	}
	log.Println("Task deadline updated: ", id, deadline)
	return nil
}

func UpdateTaskInWorkTime(id int, inWorkTime bool) error {
	err := DB.Model(&Task{}).Where("id = ?", id).Update("in_work_time", inWorkTime).Error
	if err != nil {
		log.Fatal("Failed to update task in work time: ", err)
		return err
	}
	log.Println("Task in work time updated: ", id, inWorkTime)
	return nil
}

func UpdateTaskStatus(id int, status TaskStatus) error {
	err := DB.Model(&Task{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		log.Fatal("Failed to update task status: ", err)
		return err
	}
	log.Println("Task status updated: ", id, status)
	return nil
}

func UpdateTaskParentTask(id int, parentTask int) error {
	err := DB.Model(&Task{}).Where("id = ?", id).Update("parent_task", parentTask).Error
	if err != nil {
		log.Fatal("Failed to update task parent task: ", err)
		return err
	}
	log.Println("Task parent task updated: ", id, parentTask)
	return nil
}

func GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := DB.Find(&tasks).Error
	if err != nil {
		log.Fatal("Failed to get all tasks: ", err)
		return nil, err
	}
	log.Println("All tasks: ", tasks)
	return tasks, nil
}

func GetTaskByID(id int) (Task, error) {
	var task Task
	err := DB.First(&task, id).Error
	if err != nil {
		log.Fatal("Failed to get task by ID: ", err)
		return task, err
	}
	return task, nil
}

func GetTasksByRootTask(rootTask int) ([]Task, error) {
	var tasks []Task
	err := DB.Where("root_task = ?", rootTask).Find(&tasks).Error
	if err != nil {
		log.Fatal("Failed to get tasks by root task: ", err)
		return nil, err
	}
	log.Println("Tasks by root task: ", tasks)
	return tasks, nil
}

func GetTasksByParentTask(parentTask int) ([]Task, error) {
	var tasks []Task
	err := DB.Where("parent_task = ?", parentTask).Find(&tasks).Error
	if err != nil {
		log.Fatal("Failed to get tasks by parent task: ", err)
		return nil, err
	}
	return tasks, nil
}

func CreateTask(name string, goal string, deadline int64, inWorkTime bool) (int, error) {
	task := Task{
		Name:       name,
		Goal:       goal,
		RootTask:   0,
		Deadline:   time.UnixMilli(deadline),
		InWorkTime: inWorkTime,
		Status:     Todo,
	}
	fmt.Println("Task created: ", task)
	task.ID = AddTask(task)
	err := UpdatePosition(task.ID, 0, 0)
	if err != nil {
		return -1, err
	}
	err = UpdateConstraints(task.ID, "", "")
	if err != nil {
		return -1, err
	}
	return task.ID, nil
}

func EliminateTask(id int) error {
	task, err := GetTaskByID(id)
	if err != nil || task.ID == -1 {
		return err
	}
	tasks, err := GetTasksByParentTask(id)
	err = DeleteAllRelatedTaskRelations(id)
	if err != nil {
		return err
	}
	err = DeleteTaskTriggersByID(id)
	if err != nil {
		return err
	}
	err = DeleteTaskAfterEffectByID(id)
	if err != nil {
		return err
	}
	err = DeleteSuspendedTasks(id)
	if err != nil {
		return err
	}
	for _, task := range tasks {
		err := EliminateTask(task.ID)
		if err != nil {
			return err
		}
	}
	return DeleteTask(id)
}

type TaskDetail struct {
	Task struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Goal       string `json:"goal"`
		Deadline   int64  `json:"deadline"`
		InWorkTime bool   `json:"in_work_time"`
		Status     string `json:"status"`
	} `json:"task"`
	TriggerTypes       []string `json:"trigger_type"`
	AfterEffectTypes   []string `json:"after_effect_type"`
	SuspendedTaskTypes []string `json:"suspended_task_type"`
	SuspendedTask      struct {
		ResumeTime string   `json:"resume_time"`
		Email      string   `json:"email"`
		Keywords   []string `json:"keywords"`
	} `json:"suspended_task"`
	Trigger struct {
		EventName        string `json:"event_name"`
		EventDescription string `json:"event_description"`
	} `json:"trigger"`
	AfterEffect struct {
		NowAt     int   `json:"now_at"`
		Period    int   `json:"period"`
		Intervals []int `json:"intervals"`
	} `json:"after_effect"`
	TaskConstraint struct {
		DependencyConstraint string `json:"dependency_constraint"`
		SubtaskConstraint    string `json:"subtask_constraint"`
	} `json:"task_constraint"`
}

func GetDetailedTask(id int) (TaskDetail, error) {
	task, err := GetTaskByID(id)
	if err != nil {
		return TaskDetail{}, err
	}
	taskDetail := TaskDetail{}
	taskDetail.Task.ID = task.ID
	taskDetail.Task.Name = task.Name
	taskDetail.Task.Goal = task.Goal
	taskDetail.Task.Deadline = task.Deadline.UnixMilli()
	taskDetail.Task.InWorkTime = task.InWorkTime
	statusString, err := task.Status.String()
	if err != nil {
		return TaskDetail{}, err
	}
	taskDetail.Task.Status = statusString
	triggers, err := GetTaskTriggersByID(id)
	if err != nil {
		return TaskDetail{}, err
	}
	for _, trigger := range triggers {
		triggerTypeString, err := trigger.Type.String()
		if err != nil {
			return TaskDetail{}, err
		}
		taskDetail.TriggerTypes = append(taskDetail.TriggerTypes, triggerTypeString)
		eventInfo, err := trigger.GetEventInfo()
		if err != nil {
			return TaskDetail{}, err
		}
		taskDetail.Trigger.EventName = eventInfo.EventName
		taskDetail.Trigger.EventDescription = eventInfo.EventDescription
	}

	afterEffects, err := GetTaskAfterEffectsByID(id)
	if err != nil {
		return TaskDetail{}, err
	}
	for _, afterEffect := range afterEffects {
		afterEffectTypeString, err := afterEffect.Type.String()
		if err != nil {
			return TaskDetail{}, err
		}
		taskDetail.AfterEffectTypes = append(taskDetail.AfterEffectTypes, afterEffectTypeString)
		periodicInfo, err := afterEffect.GetPeriodicInfo()
		if err != nil {
			return TaskDetail{}, err
		}
		taskDetail.AfterEffect.NowAt = periodicInfo.NowAt
		taskDetail.AfterEffect.Period = periodicInfo.Period
		taskDetail.AfterEffect.Intervals = periodicInfo.Intervals
	}

	if task.Status == Suspended {
		suspendedTask, err := GetSuspendedTask(id)
		if err != nil {
			return TaskDetail{}, err
		}
		if suspendedTask.ID != -1 {
			if suspendedTask.Type == Time {
				suspendedTimeInfo, err := suspendedTask.GetTimeInfo()
				if err != nil {
					return TaskDetail{}, err
				}
				taskDetail.SuspendedTask.ResumeTime = strconv.FormatInt(suspendedTimeInfo.Timestamp, 10)
			} else if suspendedTask.Type == Email {
				suspendedEmailInfo, err := suspendedTask.GetEmailInfo()
				if err != nil {
					return TaskDetail{}, err
				}
				taskDetail.SuspendedTask.Email = suspendedEmailInfo.Email
				taskDetail.SuspendedTask.Keywords = suspendedEmailInfo.Keywords
			} else {
				return TaskDetail{}, fmt.Errorf("unknown suspended task type")
			}
		}
	}

	taskDetail.TaskConstraint.DependencyConstraint = task.DependencyConstraint
	taskDetail.TaskConstraint.SubtaskConstraint = task.SubtaskConstraint

	return taskDetail, nil
}

func updateTaskTriggers(taskDetail TaskDetail) error {
	err := DeleteTaskTriggersByID(taskDetail.Task.ID)
	if err != nil {
		return err
	}
	for _, triggerType := range taskDetail.TriggerTypes {
		if triggerType == "Event" {
			taskTrigger := TaskTrigger{
				ID:   taskDetail.Task.ID,
				Type: Event,
			}
			err := taskTrigger.SetEventInfo(taskDetail.Trigger.EventName, taskDetail.Trigger.EventDescription)
			if err != nil {
				return err
			}
			err = AddOrUpdateTaskTrigger(taskTrigger)
			if err != nil {
				return err
			}
		} else if triggerType == "Dependency" {
			log.Fatal("Dependency trigger not implemented")
		}
	}
	return nil
}

func updateTaskAfterEffects(taskDetail TaskDetail) error {
	err := DeleteTaskAfterEffectByID(taskDetail.Task.ID)
	if err != nil {
		return err
	}
	for _, afterEffectType := range taskDetail.AfterEffectTypes {
		if afterEffectType == "Periodic" {
			taskAfterEffect := TaskAfterEffect{
				ID:   taskDetail.Task.ID,
				Type: Periodic,
			}
			err := taskAfterEffect.SetPeriodicInfo(PeriodicT{
				NowAt:     taskDetail.AfterEffect.NowAt,
				Period:    taskDetail.AfterEffect.Period,
				Intervals: taskDetail.AfterEffect.Intervals,
			})
			if err != nil {
				return err
			}
			err = AddOrUpdateTaskAfterEffect(taskAfterEffect)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateSuspendedTask(taskDetail TaskDetail) error {
	err := DeleteSuspendedTasks(taskDetail.Task.ID)
	if err != nil {
		return err
	}
	if taskDetail.Task.Status == "Suspended" {
		if len(taskDetail.SuspendedTaskTypes) != 0 {
			if taskDetail.SuspendedTaskTypes[0] == "Time" {
				parse, err := time.Parse(time.RFC3339, taskDetail.SuspendedTask.ResumeTime)
				if err != nil {
					return err
				}
				suspendedTimeInfo := SuspendedTimeInfo{
					Timestamp: parse.UnixMilli(),
				}
				jsonData, err := json.Marshal(suspendedTimeInfo)
				err = AddOrUpdateSuspendedTask(SuspendedTask{
					ID:   taskDetail.Task.ID,
					Type: Time,
					Info: jsonData,
				})
				if err != nil {
					return err
				}
			} else if taskDetail.SuspendedTaskTypes[0] == "Email" {
				suspendedEmailInfo := SuspendedEmailInfo{
					Email:    taskDetail.SuspendedTask.Email,
					Keywords: taskDetail.SuspendedTask.Keywords,
				}
				jsonData, err := json.Marshal(suspendedEmailInfo)
				err = AddOrUpdateSuspendedTask(SuspendedTask{
					ID:   taskDetail.Task.ID,
					Type: Email,
					Info: jsonData,
				})
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("unknown suspended task type")
			}
		}
	}
	return nil
}

func SetDetailedTask(taskDetail TaskDetail) error {
	task, err := GetTaskByID(taskDetail.Task.ID)
	if err != nil {
		return err
	}
	if task.ID == -1 {
		return fmt.Errorf("task not found")
	}
	task.Name = taskDetail.Task.Name
	task.Goal = taskDetail.Task.Goal
	task.Deadline = time.UnixMilli(taskDetail.Task.Deadline)
	task.InWorkTime = taskDetail.Task.InWorkTime
	task.Status.FromString(taskDetail.Task.Status)
	task.ParentTask, err = GetNowViewingTask()
	if err != nil {
		return err
	}
	err = updateTaskTriggers(taskDetail)
	if err != nil {
		return err
	}
	err = updateTaskAfterEffects(taskDetail)
	if err != nil {
		return err
	}
	err = updateSuspendedTask(taskDetail)
	if err != nil {
		return err
	}

	err = DB.Save(&task).Error
	if err != nil {
		return err
	}
	return nil
}

func HaveSubTasks(id int) bool {
	var count int64
	err := DB.Model(&Task{}).Where("parent_task = ?", id).Count(&count).Error
	if err != nil {
		log.Fatal("Failed to get subtasks count: ", err)
		return false
	}
	return count > 0
}

func GetSubTasks(id int) ([]Task, error) {
	var tasks []Task
	err := DB.Where("parent_task = ?", id).Find(&tasks).Error
	if err != nil {
		log.Fatal("Failed to get subtasks: ", err)
		return nil, err
	}
	return tasks, nil
}

func CheckParentStatus(id int) bool {
	task, err := GetTaskByID(id)
	if err != nil {
		return false
	}

	parentTaskID := task.ParentTask
	var count int64
	err = DB.Model(&Task{}).Where("parent_task = ? AND status != ?", parentTaskID, Done).Count(&count).Error
	if err != nil {
		log.Fatal("Failed to get subtasks count: ", err)
		return false
	}

	if count == 0 {
		err := UpdateTaskStatus(parentTaskID, Done)
		if err != nil {
			return false
		}
		CheckParentStatus(parentTaskID)
	}
	return true
}

type UpdateTaskUIs struct {
	TaskUIs []struct {
		ID       int `json:"id"`
		Position struct {
			X int `json:"x"`
			Y int `json:"y"`
		} `json:"position"`
	} `json:"task_uis"`
}

func UpdatePositions(updateTaskUIs UpdateTaskUIs) error {
	nowViewingTask, err := GetNowViewingTask()
	if err != nil {
		return err
	}
	if nowViewingTask == -1 {
		return nil
	}

	for _, taskUI := range updateTaskUIs.TaskUIs {
		err := DB.Model(&Task{}).Where("id = ?", taskUI.ID).Updates(Task{PositionX: taskUI.Position.X, PositionY: taskUI.Position.Y}).Error
		if err != nil {
			log.Fatal("Failed to update positions")
			return err
		}
	}

	return nil
}

func UpdatePosition(id, positionX, positionY int) error {
	err := DB.Model(&Task{}).Where("id = ?", id).Updates(Task{PositionX: positionX, PositionY: positionY}).Error
	if err != nil {
		log.Fatal("Failed to update position")
		return err
	}
	return nil
}

func UpdateConstraints(id int, dependencyConstraint, subtaskConstraint string) error {
	err := DB.Model(&Task{}).Where("id = ?", id).Updates(Task{DependencyConstraint: dependencyConstraint, SubtaskConstraint: subtaskConstraint}).Error
	if err != nil {
		log.Fatal("Failed to update constraints")
		return err
	}
	return nil
}

func checkParentStatus(id int) {
	task, err := GetTaskByID(id)
	if err != nil {
		return
	}
	parentTaskID := task.ParentTask
	var count int64
	err = DB.Model(&Task{}).Where("parent_task = ? AND status != ?", parentTaskID, Done).Count(&count).Error
	if err != nil {
		return
	}
	if count == 0 {
		err := UpdateTaskStatus(parentTaskID, Done)
		if err != nil {
			return
		}
		checkParentStatus(parentTaskID)
	}
}

const deltaTime int64 = 60 * 60 * 24 * 1000

func CompleteTask(id int) error {
	task, err := GetTaskByID(id)
	if err != nil {
		return err
	}
	if task.ID == -1 {
		return fmt.Errorf("task not found")
	}
	err = UpdateTaskStatus(id, Done)
	if err != nil {
		return err
	}
	checkParentStatus(id)
	affect, err := GetTaskAfterEffectsByID(id)
	if len(affect) == 0 {
		return nil
	}
	afterEffect := affect[0]
	if afterEffect.Type == Periodic {
		periodicInfo, err := afterEffect.GetPeriodicInfo()
		if err != nil {
			return err
		}
		if len(periodicInfo.Intervals) == 0 {
			err = DeleteTaskAfterEffectByID(id)
			if err != nil {
				return err
			}
			return nil
		}
		periodicT := PeriodicT{
			NowAt:     periodicInfo.NowAt,
			Period:    periodicInfo.Period,
			Intervals: periodicInfo.Intervals,
		}
		if periodicT.NowAt == len(periodicT.Intervals)-1 {
			periodicInfo.NowAt = 0
			periodicInfo.Period++
			err := UpdateTaskStatus(id, Todo)
			if err != nil {
				return err
			}
			err = UpdateTaskDeadline(id, task.Deadline.UnixMilli()+deltaTime)
			if err != nil {
				return err
			}
		} else {
			periodicInfo.NowAt++
			err = UpdateTaskDeadline(id, task.Deadline.UnixMilli()+int64(periodicT.Intervals[periodicT.NowAt]))
			if err != nil {
				return err
			}
		}
		err = afterEffect.SetPeriodicInfo(periodicT)
		if err != nil {
			return err
		}
		err = AddOrUpdateTaskAfterEffect(afterEffect)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetSubTasksConnectedToEnd(id int) ([]int, error) {
	var subTasksConnectedToEnd []int
	tasks, err := GetTasksByParentTask(id)
	if err != nil {
		return nil, err
	}

	relations, err := GetRelationByParentTask(id)
	if err != nil {
		return nil, err
	}

	sourceSet := make(map[int]bool)
	targetSet := make(map[int]bool)
	connectedMap := make(map[int]bool)
	for _, relation := range relations {
		sourceSet[relation.Source] = true
		targetSet[relation.Target] = true
		connectedMap[relation.Source] = true
		connectedMap[relation.Target] = true
	}

	intersection := make(map[int]bool)
	for k := range sourceSet {
		if targetSet[k] {
			intersection[k] = true
		}
	}

	sourceTarget := make(map[int]bool)
	for k := range sourceSet {
		if !intersection[k] {
			sourceTarget[k] = true
		}
	}

	targetSource := make(map[int]bool)
	for k := range targetSet {
		if !intersection[k] {
			targetSource[k] = true
		}
	}

	for k := range targetSource {
		subTasksConnectedToEnd = append(subTasksConnectedToEnd, k)
	}

	for _, task := range tasks {
		if !connectedMap[task.ID] {
			subTasksConnectedToEnd = append(subTasksConnectedToEnd, task.ID)
		}
	}

	sort.Slice(subTasksConnectedToEnd, func(i, j int) bool { return subTasksConnectedToEnd[i] < subTasksConnectedToEnd[j] })

	return subTasksConnectedToEnd, nil
}
