package utils

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"time"
	"github.com/dgrijalva/jwt-go"
	"gomonitor/config"
)
type Charset string
const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

type UserInfo struct {
	Id        string      `json:"id"`
	Name        string       `json:"name"`
	Avatar        string       `json:"avatar"`
	Ip        string      `json:"-"`
}

type CustomClaims struct {
	UserId string
	jwt.StandardClaims
}
/**
生成token
*/
func GetJwt(userid string, userInfo *UserInfo) string  {
	customClaims :=&CustomClaims{
		UserId: userid,//用户id
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(config.Gconfig.Expire)*time.Second).Unix(), // 过期时间，必须设置
			Issuer:userInfo.Name,   // 非必须，也可以填充用户名，
		},
	}
	//采用HMAC SHA256加密算法
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	tokenString,err:= token.SignedString([]byte(config.Gconfig.SecretKey))
	if err!=nil {
		fmt.Println("get token err",err)
		return ""
	}
	return tokenString
}

func ParseToken(tokenString string)(*CustomClaims,error)  {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Gconfig.SecretKey), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims,nil
	} else {
		return nil,err
	}
}
func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}

type ExeInfo struct {
	Exeid  	string      `json:"exeid"`
	Cmd        string      `json:"cmd"`
	Name        string      `json:"name"`
	Status        string       `json:"status"`
	Cancel        interface{}      `json:"-"`
}