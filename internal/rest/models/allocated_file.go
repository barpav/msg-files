package models

import "golang.org/x/exp/slices"

type AllocatedFile struct {
	Owner  string
	Name   string
	Mime   string
	Access []string
}

func (f *AllocatedFile) HasAccess(user string) bool {
	return user == f.Owner || f.Access == nil || slices.Index(f.Access, user) != -1
}
