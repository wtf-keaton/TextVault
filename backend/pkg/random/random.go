package random

import (
	"crypto/rand"
	"log"
	"math/big"
)

func String(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			log.Fatalln("Failed to generate random string: ", err)
		}
		b[i] = charset[n.Int64()]
	}

	return string(b)
}
