package attendance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"
)

const qrWindow = 15 * time.Second

type QRToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func qrFor(r sessionRow, now time.Time) QRToken {
	window := now.Unix() / int64(qrWindow/time.Second)
	return QRToken{Token: sign(r, window), ExpiresAt: time.Unix((window+1)*int64(qrWindow/time.Second), 0).UTC()}
}

func qrValid(r sessionRow, token string, now time.Time) bool {
	window := now.Unix() / int64(qrWindow/time.Second)
	for _, w := range []int64{window, window - 1} {
		if hmac.Equal([]byte(sign(r, w)), []byte(token)) {
			return true
		}
	}
	return false
}

func sign(r sessionRow, window int64) string {
	mac := hmac.New(sha256.New, r.Secret)
	mac.Write([]byte(r.ID + ":" + strconv.FormatInt(window, 10)))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)[:12])
}
