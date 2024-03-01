package jwt_plugin

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaimsPrefix struct {
	UserID     int64
	UserName   string
	GrantScope string //授权范围
}

type MyCustomClaims struct {
	MyCustomClaimsPrefix
	jwt.RegisteredClaims
}

func GenerateTokenUsingHs256(myCustomClaimsPrefix MyCustomClaimsPrefix, registeredClaims jwt.RegisteredClaims, signKey string) (string, error) {

	claim := MyCustomClaims{
		MyCustomClaimsPrefix: myCustomClaimsPrefix,
		RegisteredClaims:     registeredClaims,
		//RegisteredClaims: jwt.RegisteredClaims{
		//	Issuer:    "Auth_Server",                                   // 签发者
		//	Subject:   "Tom",                                           // 签发对象
		//	//Audience:  jwt.ClaimStrings{"Android_APP", "IOS_APP"},      //签发受众
		//	ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),   //过期时间
		//	NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //生效时间
		//	//IssuedAt:  jwt.NewNumericDate(time.Now()),                  //签发时间
		//	//ID:        encrypt_plugin.RandomString(10),                 // wt ID, 类似于盐值
		//},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(signKey))
	return token, err
}

func ParseTokenHs256(tokenStr, signKey string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil //返回签名密钥
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("claim invalid")
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		return nil, errors.New("invalid claim type")
	}

	return claims, nil
}
