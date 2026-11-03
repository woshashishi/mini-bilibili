package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

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
func (m *TokenManager) Parse(tokenString string) (*Claims, error) {
	// TODO 1：创建 Claims 指针，用来接收解析结果
	cl := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, cl, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer("mini-bilibili"),
		jwt.WithExpirationRequired())
	// TODO 2：调用 jwt.ParseWithClaims
	// 参数依次为：
	// - tokenString
	// - Claims 指针
	// - 提供验证密钥的回调函数
	// - 校验选项
	// TODO 3：解析出错，返回 nil, ErrInvalidToken
	if err != nil {
		return nil, ErrInvalidToken
	}
	// TODO 4：检查 token.Valid
	// 无效时返回 nil, ErrInvalidToken
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	// TODO 5：检查 Claims.UserID 是否大于零
	// 不符合时返回 nil, ErrInvalidToken
	if cl.UserID <= 0 {
		return nil, ErrInvalidToken
	}
	// TODO 6：返回已经验证的 Claims
	return cl, nil
}
