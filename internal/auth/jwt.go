package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenManager(
	secret string,
	ttl time.Duration,
) (*TokenManager, error) {
	// TODO 1：secret 少于 32 字节时，返回配置错误
	if len(secret) < 32 {
		return nil, errors.New("err:wrong config")
	}
	// TODO 2：ttl 必须大于 0，否则返回配置错误
	if ttl <= 0 {
		return nil, errors.New("err:wrong config")
	}
	// TODO 3：保存 secret 和 ttl
	// secret 转换为 []byte
	s := &TokenManager{
		secret: []byte(secret),
		ttl:    ttl,
	}
	return s, nil

}

func (m *TokenManager) Generate(userID uint64) (string, error) {
	// TODO 1：获取当前时间 now
	t := time.Now()
	// TODO 2：构造 Claims
	expire := t.Add(m.ttl)
	cl := &Claims{
		userID,
		jwt.RegisteredClaims{
			Issuer:    "mini-bilibili",
			IssuedAt:  jwt.NewNumericDate(t),
			ExpiresAt: jwt.NewNumericDate(expire),
		},
	}
	// UserID：传入的 userID
	// RegisteredClaims：
	//   Issuer：固定为 "mini-bilibili"
	//   IssuedAt：当前时间
	//   ExpiresAt：当前时间加 m.ttl

	// TODO 3：用 jwt.SigningMethodHS256 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
	// TODO 4：用 m.secret 签名
	st, err := token.SignedString(m.secret)
	if err != nil {
		return "", err
	}
	return st, nil
	// 返回签名后的字符串和错误

}
