package service

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"
)

var (
	// ErrInvalidSSOToken 非法的SSO token
	ErrInvalidSSOToken = errors.New("invalid sso token")
)

// DecodeSSOToken 验证SSO token的合法性，并返回原始的token信息
func DecodeSSOToken(clientSecret string, raw string) (string, error) {
	rawByte, err := base64.URLEncoding.DecodeString(raw)
	if err != nil {
		return "", err
	}
	parts := strings.Split(string(rawByte), ":")
	if len(parts) != 2 {
		return "", ErrInvalidSSOToken
	}
	vRaw := strings.TrimSpace(parts[0])
	vHash := strings.TrimSpace(parts[1])
	if len(vRaw) == 0 || len(vHash) == 0 {
		return "", ErrInvalidSSOToken
	}
	vRaw, err = url.QueryUnescape(vRaw)
	if err != nil {
		return "", err
	}
	h := hmac.New(sha512.New, []byte(clientSecret))
	if _, err = h.Write([]byte(vRaw)); err != nil {
		return "", err
	}
	if hash := hex.EncodeToString(h.Sum(nil)); hash != vHash {
		return "", ErrInvalidSSOToken
	}
	return vRaw, nil
}
