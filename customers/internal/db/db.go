package db

import (
	"customers/internal/model"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(os.Getenv("MYSQL_DSN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	_ = db.AutoMigrate(&model.Customer{})
	_ = db.AutoMigrate(&model.ProcessedTransaction{})

	return db
}
