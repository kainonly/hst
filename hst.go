package hst

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/bytedance/sonic"
	"github.com/kainonly/go/help"
	"github.com/kainonly/hst/crypto"
	"github.com/kainonly/hst/sm2base64"
	"github.com/kainonly/hst/sm4hex"
	"resty.dev/v3"
)

var (
	// JSON 编/解码使用 sonic，注册到 resty 作为 application/json 的处理函数。
	jsonEncoder resty.ContentTypeEncoder = func(w io.Writer, v any) error {
		data, err := sonic.Marshal(v)
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		return err
	}
	jsonDecoder resty.ContentTypeDecoder = func(r io.Reader, v any) error {
		data, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		return sonic.Unmarshal(data, v)
	}
)

type M map[string]any

type Hst struct {
	Option *Option
	Client *resty.Client
}

type Option struct {
	BaseURL    string `yaml:"base_url"`
	PriKey     string `yaml:"pri_key"`      // 客户端私钥
	WicoPubKey string `yaml:"wico_pub_key"` // 支付中心公钥
	EncryptKey string `yaml:"encrypt_key"`  // 加密密钥
	MerchantNo string `yaml:"merchant_no"`  // 商户号
	ChannelId  string `yaml:"channel_id"`   // 客户端身份标识
}

func NewHst(option *Option) (x *Hst, err error) {
	client := resty.New().SetBaseURL(option.BaseURL)
	client.AddContentTypeEncoder("application/json", jsonEncoder)
	client.AddContentTypeDecoder("application/json", jsonDecoder)

	x = &Hst{
		Option: option,
		Client: client,
	}
	return
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

type Body interface {
	GetTs() string
}

func (x *Hst) NewSignObjectReq(body Body) (signObjectReq *SignObjectReq, err error) {
	txn := help.SID()
	signObjectReq = &SignObjectReq{
		Timestamp: body.GetTs(),
		ReqTxn:    txn,
		IvHex:     crypto.GenerateRandomIvHex(),
		ChannelId: x.Option.ChannelId,
		Version:   "1.0",
	}
	if signObjectReq.Body, err = sonic.MarshalString(body); err != nil {
		return
	}
	if signObjectReq.Signature, err = sm2base64.Sign([]byte(signObjectReq.Body),
		x.Option.PriKey, x.Option.ChannelId); err != nil {
		return
	}
	var keyBytes []byte
	if keyBytes, err = hex.DecodeString(x.Option.EncryptKey); err != nil {
		return
	}
	var ivBytes []byte
	if ivBytes, err = hex.DecodeString(signObjectReq.IvHex); err != nil {
		return
	}
	var encryptedBytes []byte
	if encryptedBytes, err = sm4hex.EncryptCBCPadding([]byte(signObjectReq.Body), keyBytes, ivBytes); err != nil {
		return
	}
	signObjectReq.Body = base64.StdEncoding.EncodeToString(encryptedBytes)
	return
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

func (x *Hst) Request(ctx context.Context, path string, signObjectReq *SignObjectReq) (_ string, err error) {
	var resp *resty.Response
	if resp, err = x.Client.R().SetContext(ctx).
		SetBody(signObjectReq).
		Post(path); err != nil {
		return
	}

	if resp.StatusCode() != 200 {
		err = help.E(0, `第三方接口响应失败!`)
		return
	}

	var signObjectResp *SignObjectResp
	if err = sonic.Unmarshal(resp.Bytes(), &signObjectResp); err != nil {
		return
	}

	ctx = context.WithValue(ctx, "resp", signObjectResp)
	if signObjectResp.Code != "SUCCESS" {
		err = help.E(0, fmt.Sprintf(`第三方接口响应失败! code=%s msg=%s`,
			signObjectResp.Code, signObjectResp.Msg))
		return
	}
	if err = x.decryptAndVerify(signObjectResp); err != nil {
		return
	}

	return signObjectResp.Body, nil
}

func (x *Hst) decryptAndVerify(signObjectResp *SignObjectResp) (err error) {
	var ivBytes []byte
	if ivBytes, err = hex.DecodeString(signObjectResp.IvHex); err != nil {
		return
	}
	var encryptedBytes []byte
	if encryptedBytes, err = base64.StdEncoding.DecodeString(signObjectResp.Body); err != nil {
		return
	}
	var keyBytes []byte
	keyBytes, err = hex.DecodeString(x.Option.EncryptKey)
	var decrypted []byte
	if decrypted, err = sm4hex.DecryptCBCPadding(encryptedBytes, keyBytes, ivBytes); err != nil {
		err = help.E(0, `第三方接口响应解密失败!`)
		return
	}
	signObjectResp.Body = string(decrypted)

	var verified bool
	defaultUserIdHex := hex.EncodeToString([]byte(signObjectResp.ServerAppId))
	if verified, err = sm2base64.Verify(
		[]byte(signObjectResp.Body),
		signObjectResp.Signature,
		x.Option.WicoPubKey,
		defaultUserIdHex,
	); err != nil {
		return
	}
	if !verified {
		err = help.E(0, `第三方接口响应签名验证失败!`)
	}
	return
}

type SignObjectRespResult[T any] struct {
	BizSuccess bool   `json:"bizSuccess"`
	BizCode    string `json:"bizCode"`
	BizMsg     string `json:"bizMsg"`
	BizData    T      `json:"bizData"`
}
