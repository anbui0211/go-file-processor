package initialize

import (
	"fmt"
	"gofile/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql() {
	dsn := "user:pass@tcp(127.0.0.1:3306)/accounting_system?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("MySQL connected")
	global.Mdb = db
}
