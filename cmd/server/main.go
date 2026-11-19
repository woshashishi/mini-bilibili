package main

import (
	"log"
	"mini-bilibili/internal/auth"
	"mini-bilibili/internal/config"
	"mini-bilibili/internal/database"
	"mini-bilibili/internal/like"
	"mini-bilibili/internal/server"
	"mini-bilibili/internal/storage"
	"mini-bilibili/internal/upload"
	"mini-bilibili/internal/user"
	"mini-bilibili/internal/video"
	"os"
	"path/filepath"
	"runtime"
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
	if err := db.AutoMigrate(&user.User{}, &video.Video{}, &like.Like{}); err != nil {
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
	videoRepository := video.NewRepository(
		db,
	)

	videoService := video.NewService(
		videoRepository,
	)

	videoHandler := video.NewHandler(
		videoService,
	)
	likeRepository := like.NewRepository(
		db,
	)

	likeService := like.NewService(
		likeRepository,
		videoRepository,
	)

	likeHandler := like.NewHandler(
		likeService,
	)
	localStorage, err := storage.NewLocalStorage(
		/* TODO：存储目录，填写 "storage/videos" */
		"storage/videos",
	)
	if err != nil {
		// TODO：log.Fatal，结束启动
		log.Fatal()
	}

	ffmpegPath := os.Getenv("FFMPEG_PATH")
	if ffmpegPath == "" {
		if runtime.GOOS == "windows" {
			ffmpegPath = filepath.Join("tools", "ffmpeg.exe")
		} else {
			ffmpegPath = "ffmpeg"
		}
	}
	transcoder, err := upload.NewVideoTranscoder(ffmpegPath)
	if err != nil {
		log.Fatal("初始化视频转码失败: ", err)
	}
	uploadService := upload.NewService(
		/* TODO：传入 localStorage */
		localStorage,
		transcoder,
	)

	uploadHandler := upload.NewHandler(
		/* TODO：传入 uploadService */
		uploadService,
	)
	// TODO：把 Handler 传给 NewRouter
	r := server.NewRouter(userHandler, tokenManger, videoHandler, likeHandler, uploadHandler)
	// TODO 5：监听 8000 端口
	if err := r.Run(":8000"); err != nil {
		log.Fatal("服务启动失败: ", err)
	}

}
