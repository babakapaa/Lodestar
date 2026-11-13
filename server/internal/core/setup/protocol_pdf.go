package setup

import (
	"bytes"
	_ "embed"
	"fmt"
	"time"

	"github.com/go-pdf/fpdf"
)

var (
	//go:embed fonts/Inter-Regular.ttf
	interRegular []byte
	//go:embed fonts/Inter-Bold.ttf
	interBold []byte
)

func (p Protocol) PDF() ([]byte, error) {
	doc := fpdf.New("P", "mm", "A4", "")
	doc.AddUTF8FontFromBytes("Inter", "", interRegular)
	doc.AddUTF8FontFromBytes("Inter", "B", interBold)
	doc.SetMargins(20, 20, 20)
	doc.SetAutoPageBreak(true, 20)
	doc.SetTitle(fmt.Sprintf("Протокол подключения, версия %d", p.Version), true)
	doc.SetCreator("Lodestar", true)
	doc.SetFooterFunc(func() {
		doc.SetY(-15)
		doc.SetFont("Inter", "", 8)
		doc.CellFormat(0, 10, fmt.Sprintf("%d", doc.PageNo()), "", 0, "C", false, 0, "")
	})
	doc.AddPage()
	doc.SetFont("Inter", "B", 16)
	doc.MultiCell(0, 8, fmt.Sprintf("Протокол подключения, версия %d", p.Version), "", "L", false)
	doc.Ln(2)
	doc.SetFont("Inter", "", 10)
	doc.MultiCell(0, 5, signedLine(p.SignedAt), "", "L", false)
	for _, a := range p.Answers {
		doc.Ln(5)
		doc.SetFont("Inter", "B", 12)
		doc.MultiCell(0, 6, a.Title, "", "L", false)
		doc.Ln(1)
		doc.SetFont("Inter", "", 10)
		if a.Answer == "" {
			doc.MultiCell(0, 5, "Ответа нет.", "", "L", false)
			continue
		}
		doc.MultiCell(0, 5, a.Answer, "", "L", false)
		doc.Ln(1)
		doc.MultiCell(0, 5, "Основание: "+a.Basis, "", "L", false)
		doc.Ln(1)
		doc.MultiCell(0, 5, approvedLine(a.ApprovedAt), "", "L", false)
	}
	var out bytes.Buffer
	if err := doc.Output(&out); err != nil {
		return nil, fmt.Errorf("render protocol: %w", err)
	}
	return out.Bytes(), nil
}

func signedLine(at *time.Time) string {
	if at == nil {
		return "Не подписан."
	}
	return "Подписан " + at.UTC().Format(time.DateOnly) + "."
}

func approvedLine(at *time.Time) string {
	if at == nil {
		return "Не утверждено."
	}
	return "Утверждено " + at.UTC().Format(time.DateOnly) + "."
}
