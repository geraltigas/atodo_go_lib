package table

import "errors"

type TaskRelation struct {
	ParentTask int `gorm:"column:parent_task"`
	Source     int `gorm:"primaryKey"`
	Target     int `gorm:"primaryKey"`
}

func (TaskRelation) TableName() string {
	return "task_relation"
}

func InitTaskRelationTable() error {
	err := DB.AutoMigrate(&TaskRelation{})
	if err != nil {
		return err
	}
	return nil
}

func AddRelation(parentTask, source, target int) error {
	err := DB.Create(&TaskRelation{ParentTask: parentTask, Source: source, Target: target}).Error
	if err != nil {
		return err
	}
	return nil
}

func AddRelationDefault(source, target int) error {
	nowViewingTask, err2 := GetNowViewingTask()
	if err2 != nil {
		return err2
	}
	if nowViewingTask == -1 {
		return errors.New("no task is being viewed, add relation failed")
	}
	err := DB.Create(&TaskRelation{ParentTask: nowViewingTask, Source: source, Target: target}).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteRelation(source, target int) error {
	err := DB.Delete(&TaskRelation{}, "source = ? AND target = ?", source, target).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllRelatedTaskRelations(task int) error {
	err := DB.Delete(&TaskRelation{}, "source = ? OR target = ?", task, task).Error
	if err != nil {
		return err
	}
	return nil
}

func GetTargetTasks(source int) ([]int, error) {
	var targets []int
	err := DB.Model(&TaskRelation{}).Where("source = ?", source).Pluck("target", &targets).Error
	if err != nil {
		return nil, err
	}
	return targets, nil
}

func GetSourceTasks(target int) ([]int, error) {
	var sources []int
	err := DB.Model(&TaskRelation{}).Where("target = ?", target).Pluck("source", &sources).Error
	if err != nil {
		return nil, err
	}
	return sources, nil
}

func GetRelationByParentTask(parentTask int) ([]TaskRelation, error) {
	var relations []TaskRelation
	err := DB.Find(&relations, "parent_task = ?", parentTask).Error
	if err != nil {
		return nil, err
	}
	return relations, nil
}
