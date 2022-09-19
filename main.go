package main

import (
	"fmt"
	"log"

	"github.com/cfscn/thirdparty/oauth"
)

func main() {
	wxConf := &oauth.AuthConfig{ClientId: "xxx", ClientSecret: "xxx", RedirectUrl: "http://www.cfsoft.cn"}

	wxAuth := oauth.NewAuthWxWechat(wxConf)

	fmt.Print(wxAuth.GetRedirectUrl("sate")) //获取第三方登录地址

	wxRes, err := wxAuth.GetWebAccessToken("code")

	userInfo, _ := wxAuth.GetUserInfo(wxRes.AccessToken, wxRes.OpenId)

	log.Println("ssss:", err, userInfo)
}
