package sm4hex

import (
	"bytes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"

	"github.com/tjfoc/gmsm/sm4"
)

// EncryptCBCPadding 使用 SM4 CBC 模式 PKCS7 填充加密
func EncryptCBCPadding(data, keyBytes, ivBytes []byte) ([]byte, error) {
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

// DecryptCBCPadding 使用 SM4 CBC 模式 PKCS7 填充解密
func DecryptCBCPadding(encryptedData, keyBytes, ivBytes []byte) ([]byte, error) {
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

// EncryptHex 使用默认IV加密并返回Hex字符串
func EncryptHex(data string, keyHex string) (string, error) {
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", fmt.Errorf("解码密钥Hex失败: %w", err)
	}

	defaultIvHex := "cfa54e650f6ba83a83183283360ccf67"
	ivBytes, err := hex.DecodeString(defaultIvHex)
	if err != nil {
		return "", fmt.Errorf("解码IV Hex失败: %w", err)
	}

	encrypted, err := EncryptCBCPadding([]byte(data), keyBytes, ivBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encrypted), nil
}

// DecryptHex 使用默认IV解密Hex编码的密文
func DecryptHex(encryptedDataHex string, keyHex string) (string, error) {
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", fmt.Errorf("解码密钥Hex失败: %w", err)
	}

	defaultIvHex := "cfa54e650f6ba83a83183283360ccf67"
	ivBytes, err := hex.DecodeString(defaultIvHex)
	if err != nil {
		return "", fmt.Errorf("解码IV Hex失败: %w", err)
	}

	encryptedBytes, err := hex.DecodeString(encryptedDataHex)
	if err != nil {
		return "", fmt.Errorf("解码密文Hex失败: %w", err)
	}

	decrypted, err := DecryptCBCPadding(encryptedBytes, keyBytes, ivBytes)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// pkcs7Pad PKCS7 填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7Unpad PKCS7 去填充
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
