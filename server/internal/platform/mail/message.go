package mail

import (
	"bytes"
	"encoding/base64"
	"mime"
	"net/mail"
	"time"
)

type Message struct {
	To      string
	Subject string
	Body    string
}

func (m Message) encode(from *mail.Address, id string, now time.Time) []byte {
	var b bytes.Buffer
	header := func(name, value string) {
		b.WriteString(name + ": " + value + "\r\n")
	}
	header("From", from.String())
	header("To", m.To)
	header("Subject", mime.QEncoding.Encode("utf-8", m.Subject))
	header("Date", now.Format(time.RFC1123Z))
	header("Message-ID", "<"+id+">")
	header("MIME-Version", "1.0")
	header("Content-Type", "text/plain; charset=utf-8")
	header("Content-Transfer-Encoding", "base64")
	b.WriteString("\r\n")
	encoded := base64.StdEncoding.EncodeToString([]byte(m.Body))
	for len(encoded) > 76 {
		b.WriteString(encoded[:76] + "\r\n")
		encoded = encoded[76:]
	}
	b.WriteString(encoded + "\r\n")
	return b.Bytes()
}
