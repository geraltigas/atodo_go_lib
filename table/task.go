package table

import (
	"fmt"
	"log"
	"time"
)

type TaskStatus int

const (
	Todo      TaskStatus = 0
	Suspended TaskStatus = 3
	Done      TaskStatus = 4
)

type Task struct {
	ID         int       `gorm:"primaryKey;autoIncrement"`
	RootTask   int       `gorm:"column:root_task"`
	Name       string    `gorm:"type:text"`
	Goal       string    `gorm:"type:text"`
	Deadline   time.Time `gorm:"type:timestamp"`
	InWorkTime bool      `gorm:"column:in_work_time"`
	Status     TaskStatus
	ParentTask int `gorm:"column:parent_task"`
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
	log.Println("Task by ID: ", task)
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
	log.Println("Tasks by parent task: ", tasks)
	return tasks, nil
}

func CreateTask(name string, goal string, deadline int64, inWorkTime bool) error {
	task := Task{
		Name:       name,
		Goal:       goal,
		RootTask:   0,
		Deadline:   time.UnixMilli(deadline),
		InWorkTime: inWorkTime,
		Status:     Todo,
	}
	fmt.Println("Task created: ", task)
	// TODO
	//	id := AddTask(task)
	//	int64_t task_id = add_task(task);
	//	task_ui::add_or_update_task_ui(task_id, task.parent_task, 0, 0);
	//	task_constraint::set_task_constraint({task_id, "", ""});
	return nil
}

func EliminateTask(id int) error {
	_, err := GetTaskByID(id)
	if err != nil {
		return err
	}
	tasks, err := GetTasksByParentTask(id)
	// TODO
	//task_ui::delete_task_ui(task_id);
	//task_relation::remove_all_related_relations(task_id);
	//task_trigger::delete_task_triggers_by_id(task_id);
	//task_after_effect::delete_after_effect_by_id(task_id);
	//suspended_task::delete_suspended_task(task_id);
	//task_constraint::delete_task_constraint(task_id);
	for _, task := range tasks {
		err := EliminateTask(task.ID)
		if err != nil {
			return err
		}
	}
	return DeleteTask(id)
}

// TODO: task::task_detail_t task::get_detailed_task(int64_t task_id)
// TODO: bool task::set_detailed_task(const task_detail_t &task)

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

// TODO: bool task::complete_task(int64_t task_id)
// TODO: std::vector<int64_t> task::get_sub_tasks_connected_to_end(int64_t task_id)
