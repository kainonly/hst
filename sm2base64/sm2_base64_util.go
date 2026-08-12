package sm2base64

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

const defaultUserId = "1234567812345678"

// parseKeyBytes 自动识别并解码密钥字符串（支持 Hex 和 Base64 两种格式）
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

// Sign SM2 签名 (私钥进行签名) - 签名结果为 Base64 编码格式
func Sign(data []byte, privateKey string, userId string) (string, error) {
	privKeyBytes, err := parseKeyBytes(privateKey)
	if err != nil {
		return "", fmt.Errorf("解码私钥失败: %w", err)
	}

	privKey, err := x509.ParsePKCS8UnecryptedPrivateKey(privKeyBytes)
	if err != nil {
		return "", fmt.Errorf("解析SM2私钥失败: %w", err)
	}

	uid := []byte(defaultUserId)
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

// Verify SM2 验签 (公钥进行验证) - 签名为 Base64 编码格式
func Verify(data []byte, signatureBase64 string, publicKey string, userId string) (bool, error) {
	pubKeyBytes, err := parseKeyBytes(publicKey)
	if err != nil {
		return false, fmt.Errorf("解码公钥失败: %w", err)
	}

	pubKey, err := x509.ParseSm2PublicKey(pubKeyBytes)
	if err != nil {
		return false, fmt.Errorf("解析SM2公钥失败: %w", err)
	}

	uid := []byte(defaultUserId)
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
