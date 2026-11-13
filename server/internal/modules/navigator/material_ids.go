package navigator

import "slices"

func materialIDs(citations []Citation) []string {
	var out []string
	for _, c := range citations {
		if c.MaterialID != "" && !slices.Contains(out, c.MaterialID) {
			out = append(out, c.MaterialID)
		}
	}
	return out
}
