package user

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
)

var (
	ErrUsernameExists  = errors.New("username already exists")
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	// TODO：创建并返回 Service
	return &Service{
		repo: repo,
	}

}

func (s *Service) Register(
	ctx context.Context,
	username string,
	password string,
) (*User, error) {
	// TODO 1：去除 username 首尾空格
	username = strings.TrimSpace(username)
	// TODO 2：验证用户名
	if username == "" {
		return nil, ErrInvalidUsername
	}
	if len([]rune(username)) < 3 || len([]rune(username)) > 32 {
		return nil, ErrInvalidUsername
	}
	for _, r := range username {
		// 判断是否是小写字母、大写字母、数字或下划线
		isLower := r >= 'a' && r <= 'z'
		isUpper := r >= 'A' && r <= 'Z'
		isDigit := r >= '0' && r <= '9'
		isUnderscore := r == '_'

		if !isLower && !isUpper && !isDigit && !isUnderscore {
			return nil, ErrInvalidUsername
		}
	}

	// 要求：
	// - 长度为 3～32 个字符
	// - 只能包含字母、数字、下划线
	// 不符合要求时返回 ErrInvalidUsername

	// TODO 3：验证密码
	// 要求：
	if len(password) < 8 || len(password) > 72 {
		return nil, ErrInvalidPassword
	}
	// - 最少 8 字节
	// - 最多 72 字节
	// 不符合要求时返回 ErrInvalidPassword

	// TODO 4：调用 FindByUsername
	// 分清以下三种情况：
	existingUser, err := s.repo.FindByUsername(ctx, username)
	// - 找到用户：返回 ErrUsernameExists
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

	}
	if existingUser != nil {
		return nil, ErrUsernameExists
	}
	// - gorm.ErrRecordNotFound：继续执行
	// - 其他数据库错误：原样返回

	// TODO 5：使用 bcrypt.GenerateFromPassword 加密密码
	// cost 使用 bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// TODO 6：构造 User
	// Username 保存处理后的 username
	// PasswordHash 保存 bcrypt 生成的结果
	u := &User{
		Username:     username,
		PasswordHash: string(hash),
	}

	// TODO 7：调用 Repository.Create
	err = s.repo.Create(ctx, u)
	if err != nil {
		return nil, err
	}
	// TODO 8：返回创建成功的用户
	return u, nil
}
