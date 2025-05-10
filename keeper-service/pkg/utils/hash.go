package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func CalculateFileHash(file io.Reader) (string, error) {
	seeker, ok := file.(io.Seeker)
	if ok {
		_, err := seeker.Seek(0, io.SeekStart)
		if err != nil {
			return "", err
		}
	}

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
