package privacy

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Archive struct {
	userID   string
	sections Sections
	data     map[string]any
	blobs    []Blob
}

func (a Archive) Write(ctx context.Context, w io.Writer) error {
	z := zip.NewWriter(w)
	m := manifest{UserID: a.userID, ExportedAt: time.Now().UTC(), Sections: []string{}, Files: []string{}}
	for _, s := range a.sections {
		m.Sections = append(m.Sections, s.Name())
	}
	for _, b := range a.blobs {
		m.Files = append(m.Files, b.Path)
	}
	if err := writeJSON(z, "export.json", m); err != nil {
		return err
	}
	for _, s := range a.sections {
		if err := writeJSON(z, s.Name()+".json", a.data[s.Name()]); err != nil {
			return err
		}
	}
	for _, b := range a.blobs {
		if err := writeBlob(ctx, z, b); err != nil {
			return err
		}
	}
	return z.Close()
}

func writeJSON(z *zip.Writer, name string, v any) error {
	f, err := z.Create(name)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func writeBlob(ctx context.Context, z *zip.Writer, b Blob) error {
	r, err := b.Open(ctx)
	if err != nil {
		return fmt.Errorf("open %s: %w", b.Path, err)
	}
	defer func() { _ = r.Close() }()
	f, err := z.Create(b.Path)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}
