package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mini-bilibili/internal/config"
	"time"
)

func OpenMySQL(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// TODO 1：
	// 使用：
	// gorm.Open(mysql.Open(...), &gorm.Config{})
	// 建立 GORM 连接。
	//
	dsn := cfg.DSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open mysql error:%w", err)
	}
	// 连接失败时，使用 fmt.Errorf 包装错误并返回。
	// 包装错误时保留原始错误：%w。

	// TODO 2：
	// 通过 db.DB() 获得底层 *sql.DB。
	//
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("open mysql error:%w", err)
	}
	// 为什么需要底层 sql.DB？
	// 因为连接池由 database/sql 管理，不是由 GORM 模型管理。

	// TODO 3：
	// 设置连接池：
	// MaxOpenConns    = 20
	// MaxIdleConns    = 10
	// ConnMaxLifetime = 30 分钟
	//
	// 查询下面三个方法：
	// SetMaxOpenConns
	sqlDB.SetMaxOpenConns(20)
	// SetMaxIdleConns
	sqlDB.SetMaxIdleConns(10)
	// SetConnMaxLifetime
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// TODO 4：
	// 调用 Ping()，确认数据库真的可以访问。
	//
	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("数据库 Ping 失败:%w", err)
	}
	// gorm.Open 成功只表示对象初始化成功，
	// 不一定表示 MySQL 当前真的可访问。

	return db, nil
}
