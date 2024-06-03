package table

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() error {
	// check if db file exists
	// if not, create it
	// if exists, open it
	var err error
	if DB == nil {
		// try to open db
		// if failed, create it
		DB, err = gorm.Open(sqlite.Open("./data.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to open db: ", err)
			return err
		}
		log.Println("DB opened")
	}

	{
		err := InitAppStateTable()
		if err != nil {
			return err
		}
		err = InitTaskTable()
		if err != nil {
			return err
		}
		err = InitSuspendedTaskTable()
		if err != nil {
			return err
		}
		err = InitTaskAfterEffectTable()
		if err != nil {
			return err
		}
		err = InitTaskRelationTable()
		if err != nil {
			return err
		}
		err = InitTaskTriggerTable()
		if err != nil {
			return err
		}
	}

	return nil
}
