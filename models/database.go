package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
	var err error
	dsn := "root:selopia123@tcp(127.0.0.1:3306)/lmsdb?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	fmt.Println("Database connection established successfully")
	return nil
}
