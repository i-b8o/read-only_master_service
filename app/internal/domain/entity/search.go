package entity

import "time"

type Search struct {
	RID       *uint64
	RName     *string
	CID       *uint64
	CName     *string
	UpdatedAt *time.Time
	PID       *uint64
	Text      *string
	Count     *uint64
}
