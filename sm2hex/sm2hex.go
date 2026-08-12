package sm2hex

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

const defaultUserId = "1234567812345678"

// Sign SM2 签名 (私钥进行签名) - Hex 编码格式
func Sign(data []byte, privateKeyHex string, userIdHex string) (string, error) {
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("解码私钥Hex失败: %w", err)
	}

	privKey, err := x509.ParsePKCS8UnecryptedPrivateKey(privKeyBytes)
	if err != nil {
		return "", fmt.Errorf("解析SM2私钥失败: %w", err)
	}

	var uid []byte
	if userIdHex != "" {
		uid, err = hex.DecodeString(userIdHex)
		if err != nil {
			return "", fmt.Errorf("解码userIdHex失败: %w", err)
		}
	} else {
		uid = []byte(defaultUserId)
	}

	r, s, err := sm2.Sm2Sign(privKey, data, uid, rand.Reader)
	if err != nil {
		return "", fmt.Errorf("SM2签名失败: %w", err)
	}

	sigBytes, err := sm2.SignDigitToSignData(r, s)
	if err != nil {
		return "", fmt.Errorf("序列化签名失败: %w", err)
	}

	return hex.EncodeToString(sigBytes), nil
}

// Verify SM2 验签 (公钥进行验证) - Hex 编码格式
func Verify(data []byte, signatureHex string, publicKeyHex string, userIdHex string) (bool, error) {
	pubKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return false, fmt.Errorf("解码公钥Hex失败: %w", err)
	}

	pubKey, err := x509.ParseSm2PublicKey(pubKeyBytes)
	if err != nil {
		return false, fmt.Errorf("解析SM2公钥失败: %w", err)
	}

	var uid []byte
	if userIdHex != "" {
		uid, err = hex.DecodeString(userIdHex)
		if err != nil {
			return false, fmt.Errorf("解码userIdHex失败: %w", err)
		}
	} else {
		uid = []byte(defaultUserId)
	}

	sigBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, fmt.Errorf("解码签名Hex失败: %w", err)
	}

	r, s, err := sm2.SignDataToSignDigit(sigBytes)
	if err != nil {
		return false, fmt.Errorf("反序列化签名失败: %w", err)
	}

	return sm2.Sm2Verify(pubKey, data, uid, r, s), nil
}
