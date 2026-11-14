package llm

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
)

type TranscribeClient struct {
	transport transport
}

func NewTranscribeClient(cfg Config) TranscribeClient {
	return TranscribeClient{transport: newTransport(cfg)}
}

func (c TranscribeClient) Transcribe(ctx context.Context, audio io.Reader, name string) (string, error) {
	body, writer := io.Pipe()
	form := multipart.NewWriter(writer)
	go func() {
		writer.CloseWithError(c.write(form, audio, name))
	}()
	var reply transcribeReply
	if err := c.transport.send(ctx, "/audio/transcriptions", form.FormDataContentType(), body, &reply); err != nil {
		return "", err
	}
	return reply.Text, nil
}

func (c TranscribeClient) write(form *multipart.Writer, audio io.Reader, name string) error {
	if err := form.WriteField("model", c.transport.cfg.Model); err != nil {
		return err
	}
	if err := form.WriteField("response_format", "json"); err != nil {
		return err
	}
	part, err := form.CreateFormFile("file", name)
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, audio); err != nil {
		return fmt.Errorf("read audio: %w", err)
	}
	return form.Close()
}
