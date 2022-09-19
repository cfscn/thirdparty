package oauth

import (
	"errors"

	"github.com/cfscn/thirdparty/result"
	"github.com/cfscn/thirdparty/utils"
)

// QQ授权登录
type AuthQq struct {
	BaseRequest
}

func NewAuthQq(conf *AuthConfig) *AuthQq {
	authRequest := &AuthQq{}
	authRequest.Set(utils.RegisterSourceQQ, conf)

	authRequest.authorizeUrl = "https://graph.qq.com/oauth2.0/authorize"
	authRequest.TokenUrl = "https://graph.qq.com/oauth2.0/token"
	authRequest.OpenIdUrl = "https://graph.qq.com/oauth2.0/me"
	authRequest.userInfoUrl = "https://graph.qq.com/user/get_user_info"

	return authRequest
}

// 获取登录地址
func (a *AuthQq) GetRedirectUrl(state string) (*result.CodeResult, error) {
	url := utils.NewUrlBuilder(a.authorizeUrl).
		AddParam("response_type", "code").
		AddParam("client_id", a.config.ClientId).
		AddParam("redirect_uri", a.config.RedirectUrl).
		AddParam("state", a.GetState(state)).
		Build()

	_, err := utils.Get(url)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// 获取token
func (a *AuthQq) GetWebAccessToken(code string) (*result.TokenResult, error) {
	url := utils.NewUrlBuilder(a.TokenUrl).
		AddParam("grant_type", "authorization_code").
		AddParam("code", code).
		AddParam("client_id", a.config.ClientId).
		AddParam("client_secret", a.config.ClientSecret).
		AddParam("redirect_uri", a.config.RedirectUrl).
		AddParam("fmt", "json").
		Build()

	body, err := utils.Get(url)
	if err != nil {
		return nil, err
	}
	m1 := utils.JsonToMSS(body)
	if _, ok := m1["error"]; ok {
		return nil, errors.New(m1["error_description"])
	}
	// 获取OpenId
	url = utils.NewUrlBuilder(a.OpenIdUrl).
		AddParam("access_token", m1["access_token"]).
		AddParam("unionid", "1").
		AddParam("fmt", "json").
		Build()
	body, err = utils.Get(url)
	if err != nil {
		return nil, err
	}
	m2 := utils.JsonToMSS(body)
	if _, ok := m2["error"]; ok {
		return nil, errors.New(m2["error_description"])
	}
	token := &result.TokenResult{
		AccessToken:  m1["access_token"],
		RefreshToken: m1["refresh_token"],
		ExpireIn:     m1["expires_in"],
		Scope:        m1["scope"],
		TokenType:    m1["token_type"],
		OpenId:       m2["openid"],
		UnionId:      m2["unionid"],
	}
	return token, nil
}

// 获取openid
func (a *AuthQq) GetAppOpenId(accessToken string) (*result.TokenResult, error) {
	url := utils.NewUrlBuilder(a.OpenIdUrl).
		AddParam("access_token", accessToken).
		AddParam("unionid", "1").
		AddParam("fmt", "json").
		Build()
	body, err := utils.Get(url)
	if err != nil {
		return nil, err
	}
	m2 := utils.JsonToMSS(body)
	if _, ok := m2["error"]; ok {
		return nil, errors.New(m2["error_description"])
	}
	token := &result.TokenResult{
		AccessToken: accessToken,
		OpenId:      m2["openid"],
		UnionId:     m2["unionid"],
	}
	return token, nil
}

// 获取第三方用户信息
func (a *AuthQq) GetUserInfo(accessToken string, openId string) (*result.UserResult, error) {
	url := utils.NewUrlBuilder(a.userInfoUrl).
		AddParam("openid", openId).
		AddParam("access_token", accessToken).
		AddParam("oauth_consumer_key", a.config.ClientId).
		AddParam("fmt", "json").
		Build()

	body, err := utils.Get(url)
	if err != nil {
		return nil, err
	}
	m := utils.JsonToMSS(body)
	if _, ok := m["error"]; ok {
		return nil, errors.New(m["error_description"])
	}
	user := &result.UserResult{
		UserName:  m["nickname"],
		NickName:  m["nickname"],
		AvatarUrl: m["figureurl_2"],
		Source:    a.registerSource,
		Gender:    utils.GetRealGender(m["gender"]).Desc,
	}
	return user, nil
}
