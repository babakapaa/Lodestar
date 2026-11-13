package team

import (
	"math/rand/v2"
	"strings"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
)

type Distribution struct {
	MinSize int
	MaxSize int
}

func (d Distribution) Plan(teams []Team, unplaced []identity.Person) []Placement {
	people := map[string]identity.Person{}
	queue := make([]string, 0, len(unplaced))
	for _, p := range unplaced {
		people[p.UserID] = p
		queue = append(queue, p.UserID)
	}
	rand.Shuffle(len(queue), func(i, j int) { queue[i], queue[j] = queue[j], queue[i] })
	var plan []Placement
	for _, t := range teams {
		need := d.MinSize - len(t.Members)
		if need <= 0 || len(queue) == 0 {
			continue
		}
		need = min(need, len(queue))
		plan = append(plan, Placement{Team: t, Add: queue[:need]})
		queue = queue[need:]
	}
	for _, size := range d.split(len(queue)) {
		plan = append(plan, Placement{Add: queue[:size]})
		queue = queue[size:]
	}
	for i := range plan {
		if plan[i].Team.ID == "" {
			plan[i].Name = d.name(plan[i].Add, people)
		}
	}
	return plan
}

func (d Distribution) split(n int) []int {
	var sizes []int
	for n > 0 {
		if n < d.MinSize {
			return d.spread(sizes, n)
		}
		for k := (n + d.MaxSize - 1) / d.MaxSize; k*d.MinSize <= n; k++ {
			if n/k >= d.MinSize && (n+k-1)/k <= d.MaxSize {
				for i := range k {
					size := n / k
					if i < n%k {
						size++
					}
					sizes = append(sizes, size)
				}
				return sizes
			}
		}
		sizes = append(sizes, d.MinSize)
		n -= d.MinSize
	}
	return sizes
}

func (d Distribution) spread(sizes []int, n int) []int {
	for i := range sizes {
		for n > 0 && sizes[i] < d.MaxSize {
			sizes[i]++
			n--
		}
	}
	if n > 0 {
		sizes = append(sizes, n)
	}
	return sizes
}

func (Distribution) name(ids []string, people map[string]identity.Person) string {
	var names []string
	for _, id := range ids[:min(3, len(ids))] {
		p := people[id]
		first, _, _ := strings.Cut(strings.TrimSpace(p.DisplayName), " ")
		if first == "" {
			first, _, _ = strings.Cut(p.Email, "@")
		}
		names = append(names, first)
	}
	name := strings.Join(names, ", ")
	if len(ids) > 3 {
		name += "…"
	}
	if len([]rune(name)) > maxName {
		name = string([]rune(name)[:maxName])
	}
	return name
}
