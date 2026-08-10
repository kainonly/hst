package hst

import (
	"context"
	"time"

	"resty.dev/v3"
)

type Hst struct {
	Option *Option
	Client *resty.Client
}

type Option struct {
	BaseURL string `yaml:"base_url" env:"BASE_URL"`
}

func NewHst(option *Option) (x *Hst, err error) {
	x = &Hst{
		Option: option,
		Client: resty.New().SetBaseURL(option.BaseURL),
	}
	return
}

type M map[string]any

func (x *Hst) SetNow(ctx context.Context, ts time.Time) context.Context {
	return context.WithValue(ctx, "now", ts)
}

func (x *Hst) GetNow(ctx context.Context) time.Time {
	return ctx.Value("now").(time.Time)
}

type SignObjectReq struct {
	Timestamp string `json:"timestamp"` // 请求时间戳（毫秒）
	ReqTxn    string `json:"reqTxn"`    // 请求流水号（建议 UUID）
	IvHex     string `json:"ivHex"`     // SM4 随机 IV（Hex）
	ChannelId string `json:"channelId"` // 客户端身份标识（平台分配）
	Version   string `json:"version"`   // 固定 1.0
	Signature string `json:"signature"` // 对明文 body 的 SM2 签名（Base64）
	Body      string `json:"body"`      // SM4-CBC 加密后的业务 JSON（Base64）
}

type SignObjectResp struct {
	Code         string `json:"code"`         // 网关响应码，SUCCESS 表示网关受理成功
	Msg          string `json:"msg"`          // 网关响应消息
	Timestamp    string `json:"timestamp"`    // 响应时间戳
	ClientReqTxn string `json:"clientReqTxn"` // 回传的客户端请求流水号
	ServerAppId  string `json:"serverAppId"`  // 服务端应用 ID（验签用）
	IvHex        string `json:"ivHex"`        // 响应 SM4 IV（Hex）
	Version      string `json:"version"`      // 固定 1.0
	Signature    string `json:"signature"`    // 对明文 body 的 SM2 签名（Base64）
	Body         string `json:"body"`         // SM4 加密后的业务 JSON（Base64）
}
