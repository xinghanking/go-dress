package route

import (
	"github.com/gin-gonic/gin"
	"github.com/xinghanking/session"
	"go-dress/api/deepface"
	"go-dress/cookie"
	"go-dress/redis"
	"net/http"
)

func checkLogin(context *gin.Context) {
	uid := session.Get("uid")
	if uid == nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
	}
}
func Init(router *gin.Engine) {
	redis.Init()
	router.Use(session.Init(session.Options{RedisStore: redis.Client}))
	router.Use(func(context *gin.Context) {
		//context.Request.URL.Path
		cookie.Start(context)
	})
	router.Use(gin.Recovery())
	router.POST("/deepface", deepface.Exec)
}
