package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/gin-gonic/gin"
)

func SHA256(password string) string {
	w := sha256.New()
	io.WriteString(w, password)
	bw := w.Sum(nil)
	return hex.EncodeToString(bw)
}

func MiddlewareSha256() gin.HandlerFunc {
	return func(c *gin.Context) {
		password := c.Query("password")
		if password == "" {
			password = c.PostForm("password")
		}
		c.Set("password", SHA256(password))
		c.Next()
	}
}
