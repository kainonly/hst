package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm4"
	"github.com/tjfoc/gmsm/x509"
)

// SM2KeyPair SM2公私钥对
type SM2KeyPair struct {
	// 公钥 (Hex 或 Base64 格式)
	PublicKey string
	// 私钥 (Hex 或 Base64 格式)
	PrivateKey string
}

// GenerateSM2KeyPairHex 生成SM2公私钥对 (DER编码的Hex格式)
func GenerateSM2KeyPairHex() (*SM2KeyPair, error) {
	privKey, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("生成SM2密钥对失败: %w", err)
	}

	// 将私钥序列化为 PKCS8 DER 格式（与 Java 兼容）
	privKeyDER, err := x509.MarshalSm2UnecryptedPrivateKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("序列化SM2私钥失败: %w", err)
	}

	// 将公钥序列化为 X509 DER 格式（与 Java 兼容）
	pubKeyDER, err := x509.MarshalSm2PublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("序列化SM2公钥失败: %w", err)
	}

	return &SM2KeyPair{
		PublicKey:  hex.EncodeToString(pubKeyDER),
		PrivateKey: hex.EncodeToString(privKeyDER),
	}, nil
}

// GenerateSM2KeyPairBase64 生成SM2公私钥对 (DER编码的Base64格式)
func GenerateSM2KeyPairBase64() (*SM2KeyPair, error) {
	privKey, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("生成SM2密钥对失败: %w", err)
	}

	privKeyDER, err := x509.MarshalSm2UnecryptedPrivateKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("序列化SM2私钥失败: %w", err)
	}

	pubKeyDER, err := x509.MarshalSm2PublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("序列化SM2公钥失败: %w", err)
	}

	return &SM2KeyPair{
		PublicKey:  base64.StdEncoding.EncodeToString(pubKeyDER),
		PrivateKey: base64.StdEncoding.EncodeToString(privKeyDER),
	}, nil
}

// GenerateSM4KeyHex 生成SM4密钥 (Hex格式)
func GenerateSM4KeyHex() (string, error) {
	key := make([]byte, sm4.BlockSize) // SM4 密钥长度为 16 字节 (128 bit)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("生成SM4密钥失败: %w", err)
	}
	return hex.EncodeToString(key), nil
}

// GenerateSM4KeyBase64 生成SM4密钥 (Base64格式)
func GenerateSM4KeyBase64() (string, error) {
	key := make([]byte, sm4.BlockSize)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("生成SM4密钥失败: %w", err)
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// GenerateRandomIv 生成随机的SM4初始化向量（IV）
func GenerateRandomIv() ([]byte, error) {
	iv := make([]byte, sm4.BlockSize) // SM4的IV长度为16字节
	_, err := rand.Read(iv)
	if err != nil {
		return nil, fmt.Errorf("生成随机IV失败: %w", err)
	}
	return iv, nil
}

// GenerateRandomIvHex 生成随机的SM4初始化向量（IV），返回Hex格式
func GenerateRandomIvHex() string {
	iv, _ := GenerateRandomIv()
	return hex.EncodeToString(iv)
}

// GenerateRandomIvBase64 生成随机的SM4初始化向量（IV），返回Base64格式
func GenerateRandomIvBase64() string {
	iv, _ := GenerateRandomIv()
	return base64.StdEncoding.EncodeToString(iv)
}
