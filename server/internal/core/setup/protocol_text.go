package setup

import (
	"fmt"
	"strings"
	"time"
)

func (p Protocol) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Протокол подключения, версия %d\n\n", p.Version)
	if p.SignedAt != nil {
		fmt.Fprintf(&b, "Подписан %s.\n\n", p.SignedAt.UTC().Format(time.DateOnly))
	} else {
		b.WriteString("Не подписан.\n\n")
	}
	for _, a := range p.Answers {
		fmt.Fprintf(&b, "## %s\n\n", a.Title)
		if a.Answer == "" {
			b.WriteString("Ответа нет.\n\n")
			continue
		}
		fmt.Fprintf(&b, "%s\n\nОснование: %s\n\n", a.Answer, a.Basis)
		if a.ApprovedAt != nil {
			fmt.Fprintf(&b, "Утверждено %s.\n\n", a.ApprovedAt.UTC().Format(time.DateOnly))
		} else {
			b.WriteString("Не утверждено.\n\n")
		}
	}
	return b.String()
}
