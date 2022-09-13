package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/devminnu/assignment/common/constants.go"
)

const allowedChars = "ABCDEF0123456789"

func GenerateHexString(length int) (string, error) {
	max := big.NewInt(int64(len(allowedChars)))
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = allowedChars[n.Int64()]
	}
	return string(b), nil
}

func GetRegistrationURL() string {
	return fmt.Sprint(constants.API_URL)
}
