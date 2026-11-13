package team

import (
	"fmt"
	"slices"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
)

func people(n int) []identity.Person {
	list := make([]identity.Person, n)
	for i := range list {
		list[i] = identity.Person{UserID: fmt.Sprint("u", i), DisplayName: fmt.Sprint("Студент", i, " Фамилия")}
	}
	return list
}

func TestDistributionKeepsTeamsWithinSizes(t *testing.T) {
	cases := []struct {
		min, max, n int
		want        []int
	}{
		{3, 4, 10, []int{4, 3, 3}},
		{3, 3, 4, []int{3, 1}},
		{2, 4, 9, []int{3, 3, 3}},
		{3, 5, 2, []int{2}},
		{3, 4, 5, []int{4, 1}},
		{3, 3, 7, []int{3, 3, 1}},
	}
	for _, c := range cases {
		d := Distribution{MinSize: c.min, MaxSize: c.max}
		got := d.split(c.n)
		if !slices.Equal(got, c.want) {
			t.Fatalf("split(%d) in %d..%d = %v, want %v", c.n, c.min, c.max, got, c.want)
		}
	}
}

func TestDistributionFillsSmallTeamsFirstAndNamesNewOnes(t *testing.T) {
	small := Team{ID: "t1", Members: []Member{{Person: identity.Person{UserID: "x"}}}}
	plan := Distribution{MinSize: 3, MaxSize: 4}.Plan([]Team{small}, people(6))
	if len(plan) != 2 || plan[0].Team.ID != "t1" || len(plan[0].Add) != 2 || len(plan[1].Add) != 4 {
		t.Fatalf("plan = %+v", plan)
	}
	if plan[1].Name == "" || plan[1].Team.ID != "" {
		t.Fatalf("new team = %+v", plan[1])
	}
	seen := map[string]bool{}
	for _, step := range plan {
		for _, id := range step.Add {
			if seen[id] {
				t.Fatalf("%s placed twice", id)
			}
			seen[id] = true
		}
	}
}
