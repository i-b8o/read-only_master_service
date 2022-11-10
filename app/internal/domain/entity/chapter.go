package entity

import "time"

type Chapter struct {
	ID           uint64
	Pseudo       string
	Name         string
	Num          string
	RegulationID uint64
	OrderNum     uint64
	Paragraphs   []Paragraph
	UpdatedAt    *time.Time
}
