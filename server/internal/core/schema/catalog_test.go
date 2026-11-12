package schema

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestValidate(t *testing.T) {
	lo, hi := Bounds(80, 200)
	c, err := NewCatalog(
		Definition{Key: "core.theme", Type: Choice, Options: []string{"paper", "lab"}, Default: "paper"},
		Definition{Key: "core.scale", Type: Integer, Min: lo, Max: hi, Default: 100},
		Definition{Key: "core.motion", Type: Bool, Default: false},
		Definition{Key: "core.about", Type: Text, MaxLength: 3, Default: ""},
		Definition{Key: "core.zone", Type: TimeZone, Default: ""},
	)
	if err != nil {
		t.Fatalf("NewCatalog: %v", err)
	}
	cases := []struct {
		key, value string
		ok         bool
	}{
		{"core.theme", `"lab"`, true},
		{"core.theme", `"neon"`, false},
		{"core.scale", `150`, true},
		{"core.scale", `300`, false},
		{"core.scale", `"150"`, false},
		{"core.motion", `true`, true},
		{"core.about", `"абв"`, true},
		{"core.zone", `"Asia/Almaty"`, true},
		{"core.zone", `""`, true},
		{"core.zone", `"Mars/Olympus"`, false},
		{"core.zone", `"Local"`, false},
		{"core.about", `"абвг"`, false},
	}
	for _, tc := range cases {
		d, _ := c.Lookup(tc.key)
		err := d.Validate(json.RawMessage(tc.value))
		if (err == nil) != tc.ok {
			t.Errorf("%s = %s: err = %v", tc.key, tc.value, err)
		}
	}
}

func TestCatalogRejectsBadDefault(t *testing.T) {
	_, err := NewCatalog(Definition{Key: "core.theme", Type: Choice, Options: []string{"paper"}, Default: "lab"})
	if !errors.Is(err, ErrBadDefinition) {
		t.Fatalf("err = %v, want ErrBadDefinition", err)
	}
}

func TestResolveFallsBackToDefault(t *testing.T) {
	c, err := NewCatalog(Definition{Key: "core.motion", Type: Bool, Default: false})
	if err != nil {
		t.Fatalf("NewCatalog: %v", err)
	}
	values := c.Resolve(map[string]json.RawMessage{"core.motion": json.RawMessage(`"broken"`)})
	if values[0].Set || string(values[0].Value) != "false" {
		t.Fatalf("values = %+v", values)
	}
}
