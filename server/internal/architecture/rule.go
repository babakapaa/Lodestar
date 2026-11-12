package architecture

import (
	"fmt"
	"strings"
)

const internalPrefix = "github.com/babakapaa/lodestar/server/internal/"

type Rule struct{}

func (Rule) Check(fromPackage, importPath string) error {
	if !strings.HasPrefix(importPath, internalPrefix) {
		return nil
	}
	from := strings.TrimPrefix(fromPackage, internalPrefix)
	to := strings.TrimPrefix(importPath, internalPrefix)

	switch layer(from) {
	case "platform":
		if layer(to) != "platform" {
			return fmt.Errorf("platform must not import %s", to)
		}
	case "core":
		if layer(to) == "modules" || layer(to) == "app" {
			return fmt.Errorf("core must not import %s", to)
		}
	case "modules":
		if layer(to) == "app" {
			return fmt.Errorf("module must not import %s", to)
		}
		if layer(to) == "modules" && moduleName(from) != moduleName(to) && !isPublicAPI(to) {
			return fmt.Errorf("module %s may use module %s only through its api package", moduleName(from), moduleName(to))
		}
	}
	return nil
}

func layer(path string) string {
	head, _, _ := strings.Cut(path, "/")
	return head
}

func moduleName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

func isPublicAPI(path string) bool {
	parts := strings.Split(path, "/")
	return len(parts) == 3 && parts[2] == "api"
}
