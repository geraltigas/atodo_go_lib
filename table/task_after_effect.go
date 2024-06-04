package table

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type TaskAfterEffect struct {
	ID   int             `gorm:"primaryKey"`
	Type AfterEffectType `gorm:"primaryKey"`
	Info datatypes.JSON  `gorm:"column:info"`
}

type AfterEffectType int

const (
	Periodic AfterEffectType = iota
)

func (t AfterEffectType) String() (string, error) {
	names := [...]string{
		"Periodic",
	}
	if t < Periodic || t > Periodic {
		return "Unknown", nil
	}
	return names[t], nil
}

func (*TaskAfterEffect) TableName() string {
	return "task_after_effect"
}

func InitTaskAfterEffectTable() error {
	err := DB.AutoMigrate(&TaskAfterEffect{})
	if err != nil {
		return err
	}
	return nil
}

func AddOrUpdateTaskAfterEffect(tae TaskAfterEffect) error {
	err := DB.Create(&tae).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteTaskAfterEffect(id int, t AfterEffectType) error {
	err := DB.Delete(&TaskAfterEffect{}, id, t).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteTaskAfterEffectByID(id int) error {
	err := DB.Delete(&TaskAfterEffect{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func GetTaskAfterEffect(id int, t AfterEffectType) (*TaskAfterEffect, error) {
	tae := TaskAfterEffect{}
	err := DB.First(&tae, id, t).Error
	if err != nil {
		return nil, err
	}
	return &tae, nil
}

func GetTaskAfterEffectsByID(id int) ([]TaskAfterEffect, error) {
	var taes []TaskAfterEffect
	err := DB.Find(&taes, id).Error
	if err != nil {
		return nil, err
	}
	return taes, nil
}

type PeriodicT struct {
	NowAt     int
	Period    int
	Intervals []int
}

func (p *PeriodicT) Equal(other PeriodicT) bool {
	if p.NowAt != other.NowAt {
		return false
	}
	if p.Period != other.Period {
		return false
	}
	if len(p.Intervals) != len(other.Intervals) {
		return false
	}
	for i := range p.Intervals {
		if p.Intervals[i] != other.Intervals[i] {
			return false
		}
	}
	return true
}

func (tae *TaskAfterEffect) SetPeriodicInfo(info PeriodicT) error {
	infoBytes, err := json.Marshal(info)
	if err != nil {
		return err
	}
	tae.Info = infoBytes
	return nil
}

func (tae *TaskAfterEffect) GetPeriodicInfo() (*PeriodicT, error) {
	info := PeriodicT{}
	err := json.Unmarshal(tae.Info, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (tae *TaskAfterEffect) Equal(effect TaskAfterEffect) bool {
	if tae.ID != effect.ID {
		return false
	}
	if tae.Type != effect.Type {
		return false
	}
	info, err := tae.GetPeriodicInfo()
	if err != nil {
		return false
	}
	effectInfo, err := effect.GetPeriodicInfo()
	if err != nil {
		return false
	}

	if !info.Equal(*effectInfo) {
		return false
	}

	DeleteTaskAfterEffect(tae.ID, tae.Type)

	return true
}
