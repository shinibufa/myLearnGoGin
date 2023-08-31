package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/xiongjsh/learn_tiktok_project/repository"
)

var secret = []byte("haohaohao")

type MyClaims struct {
	UserId int64
	jwt.StandardClaims
}

func GenerateToken(user repository.UserLogin) (string, error) {
	expiretime := time.Now().Add(7 * 24 * time.Hour)
	myClaims := &MyClaims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "tiktok_server",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expiretime.Unix(),
			Subject:   "user",
		},
	}
	tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims).SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func ParseToken(tokenStr string) (*MyClaims, bool) {
	token, _ := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method %v", t.Header["alg"])
		}
		return secret, nil
	})
	if myclaims, ok := token.Claims.(*MyClaims); ok {
		if token.Valid {
			return myclaims, true
		} else {
			return myclaims, false
		}
	}
	return nil, false
}

func MiddlewareJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			token = c.PostForm("token")
		}
		if token == "" {
			c.JSON(http.StatusOK, repository.CommonResponse{Code: 401, Msg: "user doesn't exist"})
			c.Abort()
			return
		}
		MyClaims, ok := ParseToken(token)
		if !ok {
			c.JSON(http.StatusOK, repository.CommonResponse{Code: 402, Msg: "invalid token"})
			c.Abort()
			return
		}
		if time.Now().Unix() > MyClaims.ExpiresAt {
			c.JSON(http.StatusOK, repository.CommonResponse{Code: 403, Msg: "token expires"})
			c.Abort()
			return
		}
		c.Set("user_id", MyClaims.UserId)
		c.Next()
	}
}
