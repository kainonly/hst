package crypto

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kainonly/hst/model"
)

// ParseCertificateToEntity 解析证书文件并创建证书信息对象
func ParseCertificateToEntity(certPath, appMgmtId, payPlatform, certType, environment, certName string) (*model.CertificateInfo, error) {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("读取证书文件失败: %w", err)
	}

	// 尝试解析 PEM 格式
	block, _ := pem.Decode(certPEM)
	var certDER []byte
	if block != nil {
		certDER = block.Bytes
	} else {
		// 可能是 DER 格式
		certDER = certPEM
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, fmt.Errorf("解析证书失败: %w", err)
	}

	certInfo := &model.CertificateInfo{
		CertName:     certName,
		PayPlatform:  payPlatform,
		CertType:     certType,
		Environment:  environment,
		AppMgmtId:    appMgmtId,
		SerialNumber: strings.ToUpper(cert.SerialNumber.Text(16)),
		Subject:      cert.Subject.String(),
		Issuer:       cert.Issuer.String(),
		NotBefore:    cert.NotBefore,
		NotAfter:     cert.NotAfter,
		CertStatus:   determineCertStatus(cert),
		CertContent:  base64.StdEncoding.EncodeToString(cert.Raw),
		PublicKey:    base64.StdEncoding.EncodeToString(cert.RawSubjectPublicKeyInfo),
		Purpose:      generatePurpose(payPlatform, certType),
	}

	return certInfo, nil
}

// determineCertStatus 判断证书状态
func determineCertStatus(cert *x509.Certificate) string {
	now := time.Now()
	if now.After(cert.NotAfter) {
		return "EXPIRED"
	}
	if now.Before(cert.NotBefore) {
		return "NOT_YET_VALID"
	}
	return "ACTIVE"
}

// generatePurpose 生成证书用途描述
func generatePurpose(payPlatform, certType string) string {
	platform := "微信支付"
	if payPlatform == "ALIPAY" {
		platform = "支付宝"
	}

	switch certType {
	case "APP_PUB_CERT":
		return platform + "应用证书，用于支付签名"
	case "PLATFORM_PUBLIC":
		return platform + "平台公钥证书，用于验签"
	case "ROOT_CERT":
		return platform + "根证书，用于证书链验证"
	default:
		return platform + "证书"
	}
}
