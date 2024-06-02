package table

import (
	"fmt"
	"log"
	"time"
)

const defaultAppStateID = 0

type AppState struct {
	ID              int       `gorm:"primaryKey;check:id=0"`
	RootTask        int       `gorm:"column:root_task"`
	NowViewingTask  int       `gorm:"column:now_viewing_task"`
	NowSelectedTask int       `gorm:"column:now_selected_task"`
	WorkTime        time.Time `gorm:"column:work_time"`
	NowIsWorkTime   bool      `gorm:"column:now_is_work_time"`
	NowDoingTask    int       `gorm:"column:now_doing_task"`
}

func (AppState) TableName() string {
	return "app_state"
}

func InitAppStateTable() error {
	err := DB.AutoMigrate(&AppState{})
	if err != nil {
		log.Fatal("Failed to migrate AppState table: ", err)
		return err
	}
	log.Println("AppState table migrated")
	return nil
}

func getAppState() (AppState, error) {
	var appState AppState
	// find by id
	err := DB.First(&appState, defaultAppStateID).Error
	if err != nil {
		return appState, err
	}
	return appState, nil
}

func SetRootTask(rootTask int) error {
	// where id = 1
	DB.Model(&AppState{}).Where(fmt.Sprintf("id = %d", defaultAppStateID)).Update("root_task", rootTask)
	return nil
}

func GetRootTask() (int, error) {
	appState, err := getAppState()
	if err != nil {
		return -1, err
	}
	return appState.RootTask, nil
}

func SetNowViewingTask(nowViewingTask int) error {
	DB.Model(&AppState{}).Where(fmt.Sprintf("id = %d", defaultAppStateID)).Update("now_viewing_task", nowViewingTask)
	return nil
}

func GetNowViewingTask() (int, error) {
	appState, err := getAppState()
	if err != nil {
		return -1, err
	}
	return appState.NowViewingTask, nil
}

func SetNowSelectedTask(nowSelectedTask int) error {
	DB.Model(&AppState{}).Where(fmt.Sprintf("id = %d", defaultAppStateID)).Update("now_selected_task", nowSelectedTask)
	return nil
}

func GetNowSelectedTask() (int, error) {
	appState, err := getAppState()
	if err != nil {
		return -1, err
	}
	return appState.NowSelectedTask, nil
}

func BackToParentTask() error {
	nowViewingTask, err := GetNowViewingTask()
	if nowViewingTask == -1 {
		return err
	}
	task, err := GetTaskByID(nowViewingTask)
	if err != nil || task.ParentTask == -1 {
		return err
	}
	err = SetNowViewingTask(task.ParentTask)
	if err != nil {
		return err
	}
	return nil
}

func SetWorkTime(workTime int64) error {
	DB.Model(&AppState{}).Where(fmt.Sprintf("id = %d", defaultAppStateID)).Update("work_time", time.Unix(workTime, 0))
	return nil
}

func GetWorkTime() (int64, error) {
	appState, err := getAppState()
	if err != nil {
		return -1, err
	}
	return appState.WorkTime.Unix(), nil
}

func SetNowDoingTask(nowDoingTask int) error {
	DB.Model(&AppState{}).Where(fmt.Sprintf("id = %d", defaultAppStateID)).Update("now_doing_task", nowDoingTask)
	return nil
}

func GetNowDoingTask() (int, error) {
	appState, err := getAppState()
	if err != nil {
		return -1, err
	}
	return appState.NowDoingTask, nil
}

func SetNowIsWorkTime(nowIsWorkTime bool) error {
	DB.Model(&AppState{}).Where(fmt.Sprintf("id = %d", defaultAppStateID)).Update("now_is_work_time", nowIsWorkTime)
	return nil
}

func GetNowIsWorkTime() (bool, error) {
	appState, err := getAppState()
	if err != nil {
		return false, err
	}
	return appState.NowIsWorkTime, nil
}
