package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"douyin/service/user/internal/config"
)


func InitGorm(c config.DbConfig) (*gorm.DB, error) {
	m := config.Mysql{DbConfig: c}
	mysqlConfig := mysql.Config{
		DSN: m.Dsn(),
	}

	db, err := gorm.Open(mysql.New(mysqlConfig))
	if err != nil {
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db, nil
	}
}
