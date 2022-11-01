package result

//Ticket结果
type TicketResult struct {
	AccessToken string `json:"access_token"` // TOKEN
	Ticket      string `json:"ticket"`       // Ticket
	ExpireIn    string `json:"expire_in"`    //　过期时间
	ErrMsg      string `json:"errmsg"`       // 错误消息
	ErrCode     string `json:"errcode"`      // 错误码
}
