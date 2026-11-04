package main

import (
	"log"
	"mini-bilibili/internal/auth"
	"mini-bilibili/internal/config"
	"mini-bilibili/internal/database"
	"mini-bilibili/internal/server"
	"mini-bilibili/internal/user"
)

func main() {

	// TODO 1：加载配置
	cfg := config.Load()

	// TODO 2：连接 MySQL
	// 发生错误时使用 log.Fatal 结束程序
	db, err := database.OpenMySQL(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	// TODO 3：
	// 使用 GORM AutoMigrate 创建 users 表。
	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatal("自动迁移失败", err)
	}
	// 传入的是 User 模型的指针。

	// TODO 4：创建 Router
	// TODO：创建 Repository
	userRepository := user.NewRepository(db)

	// TODO：创建 Service
	userService := user.NewService(userRepository)
	tokenManger, err := auth.NewTokenManager(cfg.JWT.Secret, cfg.JWT.TTL)
	if err != nil {
		// TODO：使用 log.Fatal 结束启动
		log.Fatal("初始化token失败", err)
		// 不要输出密钥
	}
	// TODO：创建 Handler
	userHandler := user.NewHandler(userService, tokenManger)

	// TODO：把 Handler 传给 NewRouter
	r := server.NewRouter(userHandler, tokenManger)
	// TODO 5：监听 8000 端口
	if err := r.Run(":8000"); err != nil {
		log.Fatal("服务启动失败: ", err)
	}

}
