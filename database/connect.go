package database

import (
	"github.com/bbaktaeho/block-catcher/repository"
	"gorm.io/gorm"
)

func ConnectGORM(dialector gorm.Dialector, opt *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, opt)
	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&repository.Block{}); err != nil {
		return nil, err
	}

	return db, nil
}

func CloseGORM(db *gorm.DB) error {
	d, err := db.DB()
	if err != nil {
		return err
	}

	return d.Close()
}
