package app

import (
	"time"

	"mhelper_server/pkg/util"

	"github.com/dgrijalva/jwt-go"
	"mhelper_server/global"
)

type Claims struct {
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	jwt.StandardClaims
}

func GetJWTSecret() []byte {
	return []byte(global.JWTSetting.Secret)
}

// GenerateToken 生成Token
func GenerateToken(appKey, appSecret string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(global.JWTSetting.Expire)
	claims := Claims{
		AppKey:    util.EncodeMD5(appKey),
		AppSecret: util.EncodeMD5(appSecret),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    global.JWTSetting.Issuer,
		},
	}

	// 第一个参数是 SigningMethod，其包含 SigningMethodHS256、SigningMethodHS384、SigningMethodHS512 三种 crypto.Hash 加密算法的方案。
	// 第二个参数是 Claims，主要是用于传递用户所预定义的一些权利要求，便于后续的加密、校验等行为。
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// tokenClaims.SignedString：生成签名字符串，根据所传入 Secret 不同，进行签名并返回标准的 Token。
	token, err := tokenClaims.SignedString(GetJWTSecret())
	return token, err
}

// ParseToken 解析传入的Token
func ParseToken(token string) (*Claims, error) {

	// ParseWithClaims：用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回 *Token。
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*Claims)
		// Valid：验证基于时间的声明，例如：过期时间（ExpiresAt）、签发者（Issuer）、生效时间（Not Before）
		// 需要注意的是，如果没有任何声明在令牌中，仍然会被认为是有效的。
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
