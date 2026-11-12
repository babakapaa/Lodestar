package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
)

const challengeLength = 43

func validChallenge(challenge string) bool {
	if len(challenge) != challengeLength {
		return false
	}
	_, err := base64.RawURLEncoding.DecodeString(challenge)
	return err == nil
}

func challengeMatches(challenge, verifier string) bool {
	sum := sha256.Sum256([]byte(verifier))
	expected := base64.RawURLEncoding.EncodeToString(sum[:])
	return subtle.ConstantTimeCompare([]byte(expected), []byte(challenge)) == 1
}
