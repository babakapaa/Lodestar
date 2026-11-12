package delivery

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const fallbackLanguage = "en"

type Texts struct {
	byType map[string]map[string]module.NotificationText
}

func NewTexts(types ...module.NotificationType) Texts {
	t := Texts{byType: map[string]map[string]module.NotificationText{}}
	for _, nt := range types {
		t.byType[nt.Key] = nt.Texts
	}
	return t
}

func (t Texts) Render(item Item, languages ...string) (module.NotificationText, bool) {
	texts, ok := t.byType[item.Type]
	if !ok || len(texts) == 0 {
		return module.NotificationText{}, false
	}
	text := texts[pick(texts, languages)]
	var params map[string]any
	_ = json.Unmarshal(item.Params, &params)
	for key, value := range params {
		placeholder := "{" + key + "}"
		v := fmt.Sprint(value)
		text.Title = strings.ReplaceAll(text.Title, placeholder, v)
		text.Body = strings.ReplaceAll(text.Body, placeholder, v)
	}
	return text, true
}

func pick(texts map[string]module.NotificationText, languages []string) string {
	for _, lang := range languages {
		if _, ok := texts[lang]; ok && lang != "" {
			return lang
		}
	}
	if _, ok := texts[fallbackLanguage]; ok {
		return fallbackLanguage
	}
	keys := make([]string, 0, len(texts))
	for k := range texts {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys[0]
}
