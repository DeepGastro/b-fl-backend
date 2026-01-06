package utils

// 해시값을 계산하는 프로그램입니다!

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func CalculateFileHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
