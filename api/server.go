package api

import (
	"fmt"
	"github.com/xinghanking/session"
	"go-dress/config"
	"go-dress/cookie"
	"go-dress/models/utils"
)

func GetUserInfo() map[string]any {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	phpSessionId := cookie.Get(cfg.Session.PhpName)
	fmt.Println(phpSessionId)
	if phpSessionId != "" {
		data := map[string]any{"session_id": phpSessionId}
		response, err := utils.PostJson("http://127.0.0.1:8006/GetSession.php", data)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		d, ok := response["USER_INFO"]
		if ok {
			return d.(map[string]any)
		}
	}
	response := session.Get("USER_INFO")
	if response != nil {
		return response.(map[string]any)
	}
	return nil
}
