package model

import "time"

// CertificateInfo 证书管理信息
type CertificateInfo struct {
	// 证书唯一标识
	CertId string `json:"certId"`
	// 证书名称
	CertName string `json:"certName"`
	// 支付平台：ALIPAY-支付宝，WECHAT-微信支付
	PayPlatform string `json:"payPlatform"`
	// 证书类型
	CertType string `json:"certType"`
	// 证书环境：SANDBOX-沙箱，PRODUCTION-正式环境
	Environment string `json:"environment"`
	// 证书序列号
	SerialNumber string `json:"serialNumber"`
	// 证书主题
	Subject string `json:"subject"`
	// 证书颁发者
	Issuer string `json:"issuer"`
	// 证书有效期开始时间
	NotBefore time.Time `json:"notBefore"`
	// 证书有效期结束时间
	NotAfter time.Time `json:"notAfter"`
	// 证书状态：ACTIVE-有效，EXPIRED-已过期，REVOKED-已吊销
	CertStatus string `json:"certStatus"`
	// 证书内容，Base64编码
	CertContent string `json:"certContent"`
	// 证书公钥，Base64编码
	PublicKey string `json:"publicKey"`
	// 关联的应用管理ID
	AppMgmtId string `json:"appMgmtId"`
	// 证书用途描述
	Purpose string `json:"purpose"`
	// 备注信息
	Remark string `json:"remark"`
}
