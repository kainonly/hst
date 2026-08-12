package hexsign

import (
	"encoding/hex"
	"fmt"

	"github.com/kainonly/hst/model"
	sm2util "github.com/kainonly/hst/sm2hex"
	"github.com/kainonly/hst/sm4hex"
)

// EncryptAndSign 加签加密（先签名，后加密）
func EncryptAndSign(signObject *model.SignObjectReq, privateKeyHex, sm4KeyHex string) error {
	channelId := signObject.ChannelId
	ivHex := signObject.IvHex
	defaultUserIdHex := hex.EncodeToString([]byte(channelId))
	body := signObject.Body

	// 对业务数据JSON字符串进行签名
	signatureForBody, err := sm2util.Sign([]byte(body), privateKeyHex, defaultUserIdHex)
	if err != nil {
		return fmt.Errorf("SM2签名失败: %w", err)
	}
	signObject.Signature = signatureForBody

	// 对body进行加密处理
	keyBytes, err := hex.DecodeString(sm4KeyHex)
	if err != nil {
		return fmt.Errorf("解码SM4密钥失败: %w", err)
	}
	ivBytes, err := hex.DecodeString(ivHex)
	if err != nil {
		return fmt.Errorf("解码IV失败: %w", err)
	}

	encrypted, err := sm4hex.EncryptCBCPadding([]byte(body), keyBytes, ivBytes)
	if err != nil {
		return fmt.Errorf("SM4加密失败: %w", err)
	}
	signObject.Body = hex.EncodeToString(encrypted)

	return nil
}

// DecryptAndVerifyResp 解密并验签 SignObjectResp
func DecryptAndVerifyResp(resp *model.SignObjectResp, sm4KeyHex, serverPublicKeyHex string) (bool, error) {
	if err := decryptResp(resp, sm4KeyHex); err != nil {
		return false, err
	}
	return verifyResp(resp, serverPublicKeyHex)
}

// DecryptAndVerifyReq 解密并验签 SignObjectReq
func DecryptAndVerifyReq(signObject *model.SignObjectReq, sm4KeyHex, publicKeyHex string) (bool, error) {
	if err := decrypt(signObject, sm4KeyHex); err != nil {
		return false, err
	}
	return verify(signObject, publicKeyHex)
}

func decryptResp(resp *model.SignObjectResp, sm4KeyHex string) error {
	ivBytes, err := hex.DecodeString(resp.IvHex)
	if err != nil {
		return fmt.Errorf("解码IV失败: %w", err)
	}
	encryptedBytes, err := hex.DecodeString(resp.Body)
	if err != nil {
		return fmt.Errorf("解码Body Hex失败: %w", err)
	}
	keyBytes, err := hex.DecodeString(sm4KeyHex)
	if err != nil {
		return fmt.Errorf("解码SM4密钥失败: %w", err)
	}

	decrypted, err := sm4hex.DecryptCBCPadding(encryptedBytes, keyBytes, ivBytes)
	if err != nil {
		return fmt.Errorf("SM4解密失败: %w", err)
	}
	resp.Body = string(decrypted)
	return nil
}

func verifyResp(resp *model.SignObjectResp, serverPublicKeyHex string) (bool, error) {
	serverAppId := resp.ServerAppId
	defaultUserIdHex := hex.EncodeToString([]byte(serverAppId))
	return sm2util.Verify([]byte(resp.Body), resp.Signature, serverPublicKeyHex, defaultUserIdHex)
}

func decrypt(signObject *model.SignObjectReq, sm4KeyHex string) error {
	ivBytes, err := hex.DecodeString(signObject.IvHex)
	if err != nil {
		return fmt.Errorf("解码IV失败: %w", err)
	}
	encryptedBytes, err := hex.DecodeString(signObject.Body)
	if err != nil {
		return fmt.Errorf("解码Body Hex失败: %w", err)
	}
	keyBytes, err := hex.DecodeString(sm4KeyHex)
	if err != nil {
		return fmt.Errorf("解码SM4密钥失败: %w", err)
	}

	decrypted, err := sm4hex.DecryptCBCPadding(encryptedBytes, keyBytes, ivBytes)
	if err != nil {
		return fmt.Errorf("SM4解密失败: %w", err)
	}
	signObject.Body = string(decrypted)
	return nil
}

func verify(signObject *model.SignObjectReq, publicKeyHex string) (bool, error) {
	channelId := signObject.ChannelId
	defaultUserIdHex := hex.EncodeToString([]byte(channelId))
	return sm2util.Verify([]byte(signObject.Body), signObject.Signature, publicKeyHex, defaultUserIdHex)
}
