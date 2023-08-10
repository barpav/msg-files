package models

type AllocatedFile struct {
	Owner  string
	Name   string
	Mime   string
	Access []string
}
