package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

func printSecretKey() error {
	key := make([]byte, secretbox.KeySize)
	if _, err := rand.Read(key); err != nil {
		return err
	}
	fmt.Println(base64.StdEncoding.EncodeToString(key))
	return nil
}
