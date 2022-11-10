package entity

import "time"

type Chapter struct {
	ID           uint64      `json:"id,omitempty"`
	Pseudo       string      `json:"c_pseudo,omitempty"`
	Name         string      `json:"name"`
	Num          string      `json:"num,omitempty"`
	RegulationID uint64      `json:"regulation_id,omitempty"`
	OrderNum     uint64      `json:"order_num"`
	Paragraphs   []Paragraph `json:"paragraphs,omitempty"`
	UpdatedAt    *time.Time  `json:"updated_at,omitempty"`
}
