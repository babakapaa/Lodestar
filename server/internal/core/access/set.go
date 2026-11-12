package access

import "slices"

type Set map[Permission]struct{}

func NewSet(perms ...Permission) Set {
	s := Set{}
	s.Add(perms...)
	return s
}

func (s Set) Add(perms ...Permission) {
	for _, p := range perms {
		s[p] = struct{}{}
	}
}

func (s Set) Remove(perms ...Permission) {
	for _, p := range perms {
		delete(s, p)
	}
}

func (s Set) Has(p Permission) bool {
	if _, admin := s[Administrator]; admin {
		return true
	}
	_, ok := s[p]
	return ok
}

func (s Set) Covers(perms []Permission) bool {
	for _, p := range perms {
		if !s.Has(p) {
			return false
		}
	}
	return true
}

func (s Set) List() []Permission {
	list := make([]Permission, 0, len(s))
	for p := range s {
		list = append(list, p)
	}
	slices.Sort(list)
	return list
}
