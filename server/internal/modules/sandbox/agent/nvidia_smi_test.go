package agent

import "testing"

func TestParseLoadSkipsBrokenLines(t *testing.T) {
	load := parseLoad("0, 97\n1, [N/A]\nbad\n2, 4\n")
	if len(load) != 2 || load[0] != 97 || load[2] != 4 {
		t.Fatalf("load = %v", load)
	}
}
