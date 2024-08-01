package cookie

import (
	"github.com/gin-gonic/gin"
	"go-dress/config"
)

var Context *gin.Context
var MaxAge int
var Path string
var Domain string
var Secure = false
var HttpOnly = false

func Start(context *gin.Context) {
	Context = context
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	MaxAge = cfg.Cookie.MaxAge
	Path = cfg.Cookie.Path
	Domain = cfg.Cookie.Domain
	Secure = cfg.Cookie.Secure
	HttpOnly = cfg.Cookie.HttpOnly
	if Domain == "" {
		Domain = context.GetHeader("Host")
	}
}

func Get(key string) string {
	value, err := Context.Cookie(key)
	if err != nil {
		value = ""
	}
	return value
}

func Set(key string, value string) {
	Context.SetCookie(key, value, MaxAge, Path, Domain, Secure, HttpOnly)
}

func Delete(key string) {
	Context.SetCookie(key, "", -1, Path, Domain, false, false)
}
