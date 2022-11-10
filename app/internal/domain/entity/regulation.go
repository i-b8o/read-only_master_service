package entity

import "time"

type Regulation struct {
	Id           uint64
	Pseudo       string
	Name         string
	Abbreviation string
	Title        string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	Chapters     []Chapter
}
