package config

import (
	"fmt"
	"os"
	"time"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}
type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
}
type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

func Load() Config {
	return Config{
		Database: DatabaseConfig{
			Host:     getEnv("MYSQL_HOST", "127.0.0.1"),
			Port:     getEnv("MYSQL_PORT", "3309"),
			User:     getEnv("MYSQL_USER", "root"),
			Password: getEnv("MYSQL_PASSWORD", "123456"), // 和 docker-compose 中一致
			Name:     getEnv("MYSQL_DATABASE", "mini_bilibili"),
		},
		JWT: JWTConfig{
			// TODO：通过 os.Getenv 读取 JWT_SECRET
			// 不提供默认密钥；缺少密钥时，由 NewTokenManager 拒绝启动
			Secret: getEnv("JWT_SECRET", ""),

			// TODO：有效期设为两小时
			TTL: 2 * time.Hour,
		},
	}

}

func (c DatabaseConfig) DSN() string {
	// TODO 2：
	// 根据字段拼接 MySQL DSN。

	//
	// 目标格式：
	// 用户名:密码@tcp(主机:端口)/数据库名?charset=utf8mb4&parseTime=True&loc=Local
	//
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Name)
	return dsn
	// 思考：
	// parseTime=True 为什么需要？
	// utf8mb4 和 utf8 有什么区别？

}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
