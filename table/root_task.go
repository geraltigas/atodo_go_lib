package table

type RootTask struct {
	ID int `gorm:"primaryKey"`
}

func (RootTask) TableName() string {
	return "root_task"
}

func InitRootTaskTable() error {
	err := DB.AutoMigrate(&RootTask{})
	if err != nil {
		return err
	}
	return nil
}
