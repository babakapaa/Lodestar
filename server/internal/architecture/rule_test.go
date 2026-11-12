package architecture

import "testing"

func TestRule(t *testing.T) {
	cases := []struct {
		from    string
		imports string
		allowed bool
	}{
		{"internal/core/meta", "internal/core/module", true},
		{"internal/core/meta", "internal/platform/jsonhttp", true},
		{"internal/core/meta", "internal/modules/study", false},
		{"internal/core/meta", "internal/app", false},
		{"internal/platform/jsonhttp", "internal/core/apierror", false},
		{"internal/modules/study", "internal/core/module", true},
		{"internal/modules/study/store", "internal/modules/study/domain", true},
		{"internal/modules/plan", "internal/modules/study/api", true},
		{"internal/modules/plan", "internal/modules/study/store", false},
		{"internal/modules/plan", "internal/modules/study", false},
		{"internal/modules/plan", "internal/app", false},
		{"internal/app", "internal/modules/study", true},
		{"internal/core/meta", "net/http", true},
	}
	for _, c := range cases {
		err := Rule{}.Check(prefixed(c.from), prefixed(c.imports))
		if c.allowed && err != nil {
			t.Errorf("%s -> %s rejected: %v", c.from, c.imports, err)
		}
		if !c.allowed && err == nil {
			t.Errorf("%s -> %s allowed, want rejected", c.from, c.imports)
		}
	}
}

func prefixed(path string) string {
	if len(path) > 9 && path[:9] == "internal/" {
		return internalPrefix + path[9:]
	}
	return path
}
