package crypto

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

const defaultSaltSize = 16

// MD5Encrypt 对输入的字符串进行 MD5 加密
func MD5Encrypt(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// GenerateSalt 生成一个指定长度的随机盐 (Base64 编码)
func GenerateSalt(size int) (string, error) {
	if size <= 0 {
		return "", fmt.Errorf("salt size must be a positive number")
	}
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("生成随机盐失败: %w", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// GenerateDefaultSalt 使用默认长度 (16 字节) 生成随机盐
func GenerateDefaultSalt() (string, error) {
	return GenerateSalt(defaultSaltSize)
}

// MD5EncryptWithSalt 对输入的字符串使用指定的盐进行加盐 MD5 加密
// 加密过程：将盐和密码组合（"密码{盐}"），然后计算 MD5 哈希
func MD5EncryptWithSalt(password, salt string) string {
	saltedPassword := password + "{" + salt + "}"
	return MD5Encrypt(saltedPassword)
}

// MD5Verify 验证原始密码是否与加盐后的 MD5 哈希值匹配
func MD5Verify(password, salt, md5Password string) bool {
	newMD5 := MD5EncryptWithSalt(password, salt)
	return strings.EqualFold(md5Password, newMD5)
}
