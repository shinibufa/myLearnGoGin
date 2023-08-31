package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiongjsh/learn_tiktok_project/repository"
)

func MiddlewareNoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("user_id")
		if userId == "" {
			userId = c.PostForm("user_id")
		}
		if userId == "" {
			c.JSON(http.StatusOK, repository.CommonResponse{Code: 401, Msg: "user doesn't exist"})
			c.Abort()
			return
		}
		userIdInt, err := strconv.ParseInt(userId, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, repository.CommonResponse{Code: 401, Msg: "user doesn't exist"})
			c.Abort()
			return
		}
		c.Set("user_id", userIdInt)
		c.Next()
	}
}
