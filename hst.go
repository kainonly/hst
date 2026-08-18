package hst

import (
	"bytes"
	"context"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/bytedance/sonic"
	"github.com/kainonly/go/help"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm4"
	"github.com/tjfoc/gmsm/x509"
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
	Debug      bool   `yaml:"debug"`
	BaseURL    string `yaml:"base_url"`
	PriKey     string `yaml:"pri_key"`      // 客户端私钥
	WicoPubKey string `yaml:"wico_pub_key"` // 支付中心公钥
	EncryptKey string `yaml:"encrypt_key"`  // 加密密钥
	MerchantNo string `yaml:"merchant_no"`  // 商户号
	ChannelId  string `yaml:"channel_id"`   // 客户端身份标识
}

func NewHst(option *Option) (x *Hst, err error) {
	client := resty.New().SetBaseURL(option.BaseURL).SetDebug(option.Debug)
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
		IvHex:     generateRandomIvHex(),
		ChannelId: x.Option.ChannelId,
		Version:   "1.0",
	}
	if signObjectReq.Body, err = sonic.MarshalString(body); err != nil {
		return
	}
	if signObjectReq.Signature, err = sm2Sign([]byte(signObjectReq.Body),
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
	if encryptedBytes, err = sm4EncryptCBCPadding([]byte(signObjectReq.Body), keyBytes, ivBytes); err != nil {
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

	var signObjectResp *SignObjectResp
	if err = sonic.Unmarshal(resp.Bytes(), &signObjectResp); err != nil {
		return
	}
	ctx = context.WithValue(ctx, "resp", signObjectResp)

	if resp.StatusCode() != 200 {
		err = help.E(0, fmt.Sprintf(`txn=%s code=%s，%s`,
			signObjectResp.ClientReqTxn, signObjectResp.Code, signObjectResp.Msg))
		return
	}
	if signObjectResp.Code != "SUCCESS" {
		err = help.E(0, fmt.Sprintf(`txn=%s code=%s，%s`,
			signObjectResp.ClientReqTxn, signObjectResp.Code, signObjectResp.Msg))
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
	if decrypted, err = sm4DecryptCBCPadding(encryptedBytes, keyBytes, ivBytes); err != nil {
		return
	}
	signObjectResp.Body = string(decrypted)

	var verified bool
	defaultUserIdHex := hex.EncodeToString([]byte(signObjectResp.ServerAppId))
	if verified, err = sm2Verify(
		[]byte(signObjectResp.Body),
		signObjectResp.Signature,
		x.Option.WicoPubKey,
		defaultUserIdHex,
	); err != nil {
		return
	}

	if !verified {
		err = help.E(0, fmt.Sprintf(`txn=%s，签名验证失败`, signObjectResp.ClientReqTxn))
	}
	return
}

type SignObjectRespResult[T any] struct {
	BizSuccess bool   `json:"bizSuccess"`
	BizCode    string `json:"bizCode"`
	BizMsg     string `json:"bizMsg"`
	BizData    T      `json:"bizData"`
}

// bizError 构造业务层错误，统一格式：第三方接口业务失败! bizCode=%s bizMsg=%s
func bizError(bizCode, bizMsg string) error {
	return help.E(0, fmt.Sprintf(`bizCode=%s bizMsg=%s`, bizCode, bizMsg))
}

// FileManifest 文件哈希清单。
// 每个字段对应一种资质，值为按上传顺序排列的 SM3 哈希（64 位 Hex）。
// 无需上传的字段保持 nil，序列化后等同空列表。
type FileManifest struct {
	CertPhotoAFiles             []string `json:"certPhotoAFiles,omitempty"`             // 身份证人像面
	CertPhotoBFiles             []string `json:"certPhotoBFiles,omitempty"`             // 身份证国徽面
	LicensePhotoFiles           []string `json:"licensePhotoFiles,omitempty"`           // 营业执照
	PrgPhotoFiles               []string `json:"prgPhotoFiles,omitempty"`               // 组织机构代码证
	IndustryLicensePhotoFiles   []string `json:"industryLicensePhotoFiles,omitempty"`   // 开户许可证
	ShopPhotoFiles              []string `json:"shopPhotoFiles,omitempty"`              // 门头照
	OtherPhotoFiles             []string `json:"otherPhotoFiles,omitempty"`             // 其他资料
	CertPhotoCFiles             []string `json:"certPhotoCFiles,omitempty"`             // 手持身份证
	RegisterProtocolPhotoFiles  []string `json:"registerProtocolPhotoFiles,omitempty"`  // 商户入驻协议
	ContractPhotoFiles          []string `json:"contractPhotoFiles,omitempty"`          // 租赁协议
	ShopEntrancePhotoFiles      []string `json:"shopEntrancePhotoFiles,omitempty"`      // 门店内景
	CheckstandPhotoFiles        []string `json:"checkstandPhotoFiles,omitempty"`        // 收银台
	MerchantAgreementPhotoFiles []string `json:"merchantAgreementPhotoFiles,omitempty"` // 商户协议
}

// ---- 以下函数原属子包 crypto / sm2base64 / sm4hex，现已合并入 hst 包 ----

const defaultSm2UserId = "1234567812345678"

// generateRandomIvHex 生成随机的 SM4 初始化向量（IV），返回 Hex 格式。
func generateRandomIvHex() string {
	iv := make([]byte, sm4.BlockSize)
	_, _ = rand.Read(iv)
	return hex.EncodeToString(iv)
}

// parseKeyBytes 自动识别并解码密钥字符串（支持 Hex 和 Base64 两种格式）。
func parseKeyBytes(keyStr string) ([]byte, error) {
	// 先尝试 Hex 解码
	keyBytes, err := hex.DecodeString(keyStr)
	if err == nil {
		return keyBytes, nil
	}
	// 再尝试 Base64 解码
	keyBytes, err = base64.StdEncoding.DecodeString(keyStr)
	if err == nil {
		return keyBytes, nil
	}
	return nil, fmt.Errorf("无法解码密钥，既不是有效的Hex也不是Base64格式")
}

// sm2Sign SM2 签名（私钥签名），签名结果为 Base64 编码格式。
func sm2Sign(data []byte, privateKey string, userId string) (string, error) {
	privKeyBytes, err := parseKeyBytes(privateKey)
	if err != nil {
		return "", fmt.Errorf("解码私钥失败: %w", err)
	}

	privKey, err := x509.ParsePKCS8UnecryptedPrivateKey(privKeyBytes)
	if err != nil {
		return "", fmt.Errorf("解析SM2私钥失败: %w", err)
	}

	uid := []byte(defaultSm2UserId)
	if userId != "" {
		uid = []byte(userId)
	}

	r, s, err := sm2.Sm2Sign(privKey, data, uid, rand.Reader)
	if err != nil {
		return "", fmt.Errorf("SM2签名失败: %w", err)
	}

	sigBytes, err := sm2.SignDigitToSignData(r, s)
	if err != nil {
		return "", fmt.Errorf("序列化签名失败: %w", err)
	}

	return base64.StdEncoding.EncodeToString(sigBytes), nil
}

// sm2Verify SM2 验签（公钥验证），签名为 Base64 编码格式。
func sm2Verify(data []byte, signatureBase64 string, publicKey string, userId string) (bool, error) {
	pubKeyBytes, err := parseKeyBytes(publicKey)
	if err != nil {
		return false, fmt.Errorf("解码公钥失败: %w", err)
	}

	pubKey, err := x509.ParseSm2PublicKey(pubKeyBytes)
	if err != nil {
		return false, fmt.Errorf("解析SM2公钥失败: %w", err)
	}

	uid := []byte(defaultSm2UserId)
	if userId != "" {
		uid = []byte(userId)
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false, fmt.Errorf("解码签名Base64失败: %w", err)
	}

	r, s, err := sm2.SignDataToSignDigit(sigBytes)
	if err != nil {
		return false, fmt.Errorf("反序列化签名失败: %w", err)
	}

	return sm2.Sm2Verify(pubKey, data, uid, r, s), nil
}

// sm4EncryptCBCPadding 使用 SM4 CBC 模式 PKCS7 填充加密。
func sm4EncryptCBCPadding(data, keyBytes, ivBytes []byte) ([]byte, error) {
	block, err := sm4.NewCipher(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("创建SM4 cipher失败: %w", err)
	}

	paddedData := pkcs7Pad(data, block.BlockSize())
	ciphertext := make([]byte, len(paddedData))

	mode := cipher.NewCBCEncrypter(block, ivBytes)
	mode.CryptBlocks(ciphertext, paddedData)

	return ciphertext, nil
}

// sm4DecryptCBCPadding 使用 SM4 CBC 模式 PKCS7 填充解密。
func sm4DecryptCBCPadding(encryptedData, keyBytes, ivBytes []byte) ([]byte, error) {
	block, err := sm4.NewCipher(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("创建SM4 cipher失败: %w", err)
	}

	plaintext := make([]byte, len(encryptedData))

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(plaintext, encryptedData)

	plaintext, err = pkcs7Unpad(plaintext, block.BlockSize())
	if err != nil {
		return nil, fmt.Errorf("PKCS7去填充失败: %w", err)
	}

	return plaintext, nil
}

// pkcs7Pad PKCS7 填充。
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7Unpad PKCS7 去填充。
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("无效的填充数据")
	}
	padding := int(data[len(data)-1])
	if padding > blockSize || padding == 0 {
		return nil, fmt.Errorf("无效的填充值: %d", padding)
	}
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("无效的PKCS7填充")
		}
	}
	return data[:len(data)-padding], nil
}
